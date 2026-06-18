package planner

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

// Options captures user choices for v0.2 target-aware planning.
type Options struct {
	Preset    string
	Domain    string
	CI        string
	Target    string
	OutputDir string
	AppName   string
	Region    string
	EnvFile   string
}

// Plan converts a project scan detection into a structured deployment blueprint.
func Plan(det *types.Detection, preset, domain, ci string) (*types.Plan, error) {
	return PlanWithOptions(det, Options{
		Preset: preset,
		Domain: domain,
		CI:     ci,
		Target: "vps",
	})
}

// PlanWithOptions converts project detection into a provider-aware v0.2 deployment plan.
func PlanWithOptions(det *types.Detection, opts Options) (*types.Plan, error) {
	target := opts.Target
	if target == "" {
		target = "vps"
	}

	preset := opts.Preset
	if preset == "" {
		preset = "compose-traefik"
	}

	ci := opts.CI
	if ci == "" {
		ci = "github"
	}

	envFile := opts.EnvFile
	if envFile == "" {
		envFile = ".env"
	}

	appName := opts.AppName
	if appName == "" {
		appName = sanitizeName(det.Project.Name)
	}

	region := opts.Region
	if region == "" {
		region = "iad"
	}

	if target == "vps" && opts.Domain == "" {
		return nil, errors.New("domain cannot be empty")
	}

	plan, err := basePlan(det, preset, opts.Domain, ci)
	if err != nil {
		return nil, err
	}

	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = ".deploy"
	}

	plan.Version = "0.2"
	plan.Target = types.TargetPlan{
		Kind:      target,
		Preset:    preset,
		Root:      det.Project.Root,
		OutputDir: outputDir,
		AppName:   appName,
		Region:    region,
		EnvFile:   envFile,
	}
	plan.Artifacts = artifactsFor(plan)
	plan.Checks = checksFor(plan)
	plan.Credentials = credentialsFor(plan)
	plan.SecretImports = secretImportsFor(plan)
	plan.Rollback = rollbackFor(plan)
	plan.Warnings = append(plan.Warnings, targetWarnings(plan)...)
	plan.ManualSteps = manualStepsFor(plan)
	plan.Actions = actionsFor(plan)
	plan.GeneratedFiles = generatedFilesFor(plan)

	return plan, nil
}

func basePlan(det *types.Detection, preset, domain, ci string) (*types.Plan, error) {
	if domain == "" {
		domain = ""
	}

	if preset == "" {
		preset = "compose-traefik"
	}

	if ci == "" {
		ci = "github"
	}

	// 1. Establish project service names
	publicService := "web"

	// 2. Build services array
	var services []types.Service

	// App service
	appService := types.Service{
		Name:         publicService,
		Type:         "app",
		Build:        "..",
		InternalPort: det.Network.InternalPort,
		Public:       true,
	}
	services = append(services, appService)

	// Database services
	if dep, ok := det.Dependencies["postgres"]; ok && dep.Likely {
		services = append(services, types.Service{
			Name:   "postgres",
			Type:   "database",
			Image:  "postgres:16-alpine",
			Public: false,
			Volume: "postgres_data",
		})
	} else if dep, ok := det.Dependencies["mysql"]; ok && dep.Likely {
		services = append(services, types.Service{
			Name:   "mysql",
			Type:   "database",
			Image:  "mysql:8.0",
			Public: false,
			Volume: "mysql_data",
		})
	}

	if dep, ok := det.Dependencies["redis"]; ok && dep.Likely {
		services = append(services, types.Service{
			Name:   "redis",
			Type:   "database",
			Image:  "redis:7-alpine",
			Public: false,
			Volume: "redis_data",
		})
	}

	// 3. Compile Env requirements
	envRequired := make([]string, len(det.Env.Keys))
	copy(envRequired, det.Env.Keys)

	envSensitive := make([]string, len(det.Env.Sensitive))
	copy(envSensitive, det.Env.Sensitive)

	// 4. Determine generated files list
	generatedFiles := []string{
		".deploy/docker-compose.yml",
		".deploy/.env.example",
		".deploy/.gitignore",
		".deploy/scripts/deploy.sh",
		".deploy/scripts/rollback.sh",
		".deploy/scripts/status.sh",
		".deploy/scripts/backup.sh",
		".deploy/README.md",
	}

	if !det.Containerization.DockerfilePresent {
		generatedFiles = append(generatedFiles, ".deploy/Dockerfile", ".deploy/.dockerignore")
	}

	if preset == "compose-traefik" || preset == "swarm-traefik" {
		generatedFiles = append(generatedFiles, ".deploy/traefik/dynamic.yml")
	} else if preset == "compose-nginx" {
		generatedFiles = append(generatedFiles, ".deploy/nginx/default.conf")
	}

	if ci == "github" {
		generatedFiles = append(generatedFiles, ".deploy/ci/github-actions.yml")
	}

	// 5. Setup Action engine pipeline
	var actions []types.Action
	if preset == "compose-traefik" {
		actions = append(actions, types.Action{
			Type: "create_network",
			Name: "web",
		})
	} else if preset == "swarm-traefik" {
		actions = append(actions, types.Action{
			Type: "create_swarm_network",
			Name: "web",
		})
	}

	if preset == "swarm-traefik" {
		actions = append(actions, types.Action{
			Type: "swarm_deploy",
			File: filepath.Join(".deploy", "docker-compose.yml"),
		})
	} else {
		actions = append(actions, types.Action{
			Type: "compose_up",
			File: filepath.Join(".deploy", "docker-compose.yml"),
		})
	}

	// 6. Gather warnings
	warnings := make([]string, len(det.Warnings))
	copy(warnings, det.Warnings)

	// Add database warning if one was generated but no migrate command exists
	for _, svc := range services {
		if svc.Type == "database" {
			warnings = append(warnings, fmt.Sprintf("Database migration is not automated. Run migrations manually after deployment for %s.", svc.Name))
			break
		}
	}

	// 7. Manual setup steps
	manualSteps := []string{
		fmt.Sprintf("Create DNS A record for %s pointing to this VPS.", domain),
		"Fill real secret values in .env on the VPS.",
		"Review .deploy/docker-compose.yml before applying.",
	}

	plan := &types.Plan{
		Version:        "0.2",
		Project:        types.ProjectPlan{Name: det.Project.Name, Root: det.Project.Root},
		Preset:         preset,
		Domain:         domain,
		PublicService:  publicService,
		Runtime:        types.RuntimePlan{Name: det.Runtime.Name, Framework: det.Runtime.Framework, Confidence: det.Runtime.Confidence},
		Build:          det.Build,
		Network:        det.Network,
		Services:       services,
		Env:            types.EnvPlan{Required: envRequired, Sensitive: envSensitive},
		GeneratedFiles: generatedFiles,
		Actions:        actions,
		Warnings:       warnings,
		ManualSteps:    manualSteps,
		FileHashes:     make(map[string]string),
	}

	return plan, nil
}

func generatedFilesFor(plan *types.Plan) []string {
	switch plan.Target.Kind {
	case "vercel":
		files := []string{
			".deploy/plan.json",
			".deploy/.env.example",
			".deploy/.gitignore",
			".deploy/README.md",
			".deploy/reports/plan-report.md",
		}
		if plan.Target.Preset != "none" {
			files = append(files, ".deploy/ci/github-actions.yml")
		}
		return files
	case "fly":
		return []string{
			".deploy/plan.json",
			".deploy/.env.example",
			".deploy/.gitignore",
			".deploy/README.md",
			".deploy/reports/plan-report.md",
		}
	default:
		files := []string{
			".deploy/plan.json",
			".deploy/docker-compose.yml",
			".deploy/.env.example",
			".deploy/.gitignore",
			".deploy/scripts/deploy.sh",
			".deploy/scripts/rollback.sh",
			".deploy/scripts/status.sh",
			".deploy/scripts/backup.sh",
			".deploy/README.md",
		}
		if !hasArtifact(plan, "Dockerfile") {
			files = append(files, ".deploy/Dockerfile", ".deploy/.dockerignore")
		}
		if plan.Preset == "compose-traefik" || plan.Preset == "swarm-traefik" {
			files = append(files, ".deploy/traefik/dynamic.yml")
		} else if plan.Preset == "compose-nginx" {
			files = append(files, ".deploy/nginx/default.conf")
		}
		files = append(files, ".deploy/ci/github-actions.yml")
		return files
	}
}

func artifactsFor(plan *types.Plan) []types.Artifact {
	var artifacts []types.Artifact
	switch plan.Target.Kind {
	case "vercel":
		artifacts = append(artifacts, types.Artifact{
			Path:     "vercel.json",
			Kind:     "provider-config",
			Scope:    "root",
			Template: "provider/vercel.json.tmpl",
			Mode:     "0644",
		})
	case "fly":
		artifacts = append(artifacts,
			types.Artifact{Path: "Dockerfile", Kind: "dockerfile", Scope: "root", Mode: "0644"},
			types.Artifact{Path: ".dockerignore", Kind: "dockerignore", Scope: "root", Template: "dockerfile/dockerignore.tmpl", Mode: "0644"},
			types.Artifact{Path: "fly.toml", Kind: "provider-config", Scope: "root", Template: "provider/fly.toml.tmpl", Mode: "0644"},
		)
	default:
		artifacts = append(artifacts,
			types.Artifact{Path: "Dockerfile", Kind: "dockerfile", Scope: "root", Mode: "0644"},
			types.Artifact{Path: ".dockerignore", Kind: "dockerignore", Scope: "root", Template: "dockerfile/dockerignore.tmpl", Mode: "0644"},
		)
	}
	return artifacts
}

func checksFor(plan *types.Plan) []types.Check {
	switch plan.Target.Kind {
	case "vercel":
		return []types.Check{
			{Type: "command", Name: "vercel", Required: true},
			{Type: "schema", Name: "vercel.json", Required: true},
		}
	case "fly":
		return []types.Check{
			{Type: "command", Name: "fly", Required: true},
			{Type: "schema", Name: "fly.toml", Required: true},
			{Type: "schema", Name: "Dockerfile", Required: true},
		}
	default:
		return []types.Check{
			{Type: "command", Name: "docker", Required: true},
			{Type: "compose", Name: "docker-compose.yml", Required: true},
		}
	}
}

func credentialsFor(plan *types.Plan) []types.CredentialRequirement {
	switch plan.Target.Kind {
	case "vercel":
		return []types.CredentialRequirement{{
			Name:        "Vercel authentication",
			EnvVar:      "VERCEL_TOKEN",
			Command:     "vercel login",
			Description: "Use VERCEL_TOKEN for non-interactive deploys or an existing Vercel CLI login.",
			Required:    true,
		}}
	case "fly":
		return []types.CredentialRequirement{{
			Name:        "Fly.io authentication",
			EnvVar:      "FLY_ACCESS_TOKEN",
			Command:     "fly auth login",
			Description: "Use FLY_ACCESS_TOKEN for non-interactive deploys or an existing Fly CLI login.",
			Required:    true,
		}}
	default:
		return nil
	}
}

func secretImportsFor(plan *types.Plan) []types.SecretImport {
	keys := append([]string{}, plan.Env.Required...)
	sort.Strings(keys)
	if len(keys) == 0 {
		return nil
	}
	if plan.Target.Kind != "vercel" && plan.Target.Kind != "fly" {
		return nil
	}
	return []types.SecretImport{{
		SourceFile: plan.Target.EnvFile,
		Keys:       keys,
		Mode:       "import-env-file",
	}}
}

func rollbackFor(plan *types.Plan) types.RollbackPlan {
	switch plan.Target.Kind {
	case "vercel":
		return types.RollbackPlan{Strategy: "provider-rollback", StateFile: ".deploy/state/vercel.json"}
	case "fly":
		return types.RollbackPlan{Strategy: "redeploy-previous-image", StateFile: ".deploy/state/fly.json"}
	default:
		return types.RollbackPlan{Strategy: "restore-compose-backup", StateFile: ".deploy/backups"}
	}
}

func actionsFor(plan *types.Plan) []types.Action {
	switch plan.Target.Kind {
	case "vercel":
		return []types.Action{
			{Type: "vercel_link"},
			{Type: "vercel_import_env"},
			{Type: "vercel_deploy"},
			{Type: "vercel_alias"},
		}
	case "fly":
		return []types.Action{
			{Type: "fly_launch_if_missing"},
			{Type: "fly_import_secrets"},
			{Type: "fly_deploy"},
			{Type: "fly_attach_domain"},
		}
	default:
		actions := []types.Action{}
		if plan.Preset == "compose-traefik" {
			actions = append(actions, types.Action{Type: "create_network", Name: "web"})
		} else if plan.Preset == "swarm-traefik" {
			actions = append(actions, types.Action{Type: "create_swarm_network", Name: "web"})
		}

		if plan.Preset == "swarm-traefik" {
			actions = append(actions, types.Action{Type: "swarm_deploy", File: filepath.Join(".deploy", "docker-compose.yml")})
		} else {
			actions = append(actions, types.Action{Type: "compose_up", File: filepath.Join(".deploy", "docker-compose.yml")})
		}
		return actions
	}
}

func manualStepsFor(plan *types.Plan) []string {
	switch plan.Target.Kind {
	case "vercel":
		steps := []string{"Ensure Vercel CLI is installed and authenticated."}
		if plan.Domain != "" {
			steps = append(steps, fmt.Sprintf("Configure DNS for %s after depctl aliases the deployment.", plan.Domain))
		}
		if len(plan.SecretImports) > 0 {
			steps = append(steps, fmt.Sprintf("Review %s before importing secrets to Vercel.", plan.Target.EnvFile))
		}
		return steps
	case "fly":
		steps := []string{"Ensure flyctl is installed and authenticated."}
		if plan.Domain != "" {
			steps = append(steps, fmt.Sprintf("Configure DNS records shown by fly certs for %s.", plan.Domain))
		}
		if len(plan.SecretImports) > 0 {
			steps = append(steps, fmt.Sprintf("Review %s before importing secrets to Fly.", plan.Target.EnvFile))
		}
		return steps
	default:
		return []string{
			fmt.Sprintf("Create DNS A record for %s pointing to this VPS.", plan.Domain),
			fmt.Sprintf("Fill real secret values in %s on the VPS.", plan.Target.EnvFile),
			"Review .deploy/docker-compose.yml before applying.",
		}
	}
}

func targetWarnings(plan *types.Plan) []string {
	var warnings []string
	if plan.Target.Kind == "vercel" {
		if plan.Runtime.Name != "node" || (plan.Runtime.Framework != "nextjs" && plan.Runtime.Framework != "vite") {
			warnings = append(warnings, "Vercel target is best for Next.js and static frontend projects. Use fly or vps for long-running servers.")
		}
	}
	if plan.Target.Kind == "fly" && plan.Runtime.Name == "unknown" {
		warnings = append(warnings, "Fly target requires a Dockerfile-capable app. Runtime detection is unknown.")
	}
	return warnings
}

func sanitizeName(name string) string {
	if name == "" {
		return "depctl-app"
	}
	var b []rune
	lastDash := false
	for _, r := range name {
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b = append(b, r)
			lastDash = false
			continue
		}
		if !lastDash {
			b = append(b, '-')
			lastDash = true
		}
	}
	out := string(b)
	out = filepath.Base(out)
	for len(out) > 0 && out[0] == '-' {
		out = out[1:]
	}
	for len(out) > 0 && out[len(out)-1] == '-' {
		out = out[:len(out)-1]
	}
	if out == "" {
		return "depctl-app"
	}
	return out
}

func hasArtifact(plan *types.Plan, path string) bool {
	for _, artifact := range plan.Artifacts {
		if artifact.Path == path {
			return true
		}
	}
	return false
}
