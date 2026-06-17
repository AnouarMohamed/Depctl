package planner

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

// Plan converts a project scan detection into a structured deployment blueprint.
func Plan(det *types.Detection, preset, domain, ci string) (*types.Plan, error) {
	if domain == "" {
		return nil, errors.New("domain cannot be empty")
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
		Build:        ".",
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

	if preset == "compose-traefik" {
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
	}

	actions = append(actions, types.Action{
		Type: "compose_up",
		File: filepath.Join(".deploy", "docker-compose.yml"),
	})

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
		Version:        "0.1",
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
