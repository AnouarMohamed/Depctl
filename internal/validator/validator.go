package validator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AnouarMohamed/Depctl/internal/planfile"
	"github.com/AnouarMohamed/Depctl/internal/types"
	"gopkg.in/yaml.v3"
)

// Result contains the outcome of a validation run.
type Result struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// Validate checks the deployment kit for correctness and safety.
func Validate(outputDir string) (*Result, error) {
	res := &Result{Valid: true}

	// 1. Check if plan.json exists and parses cleanly
	plan, err := planfile.Load(outputDir)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		res.Valid = false
		return res, nil // Critical error, stop here
	}

	// 2. Domain check
	if plan.Target.Kind == "vps" && plan.Domain == "" {
		res.Errors = append(res.Errors, "domain is empty in plan.json")
		res.Valid = false
	}

	// 3. Runtime confidence check
	if plan.Runtime.Confidence < 0.5 {
		res.Warnings = append(res.Warnings, fmt.Sprintf("low detection confidence for %s (%0.2f)", plan.Runtime.Name, plan.Runtime.Confidence))
	}

	// 4. Check if required files exist
	for _, relPath := range plan.GeneratedFiles {
		cleanPath := relPath
		if len(relPath) > 8 && relPath[:8] == ".deploy/" {
			cleanPath = relPath[8:]
		}

		absPath := filepath.Join(outputDir, cleanPath)
		if _, err := os.Stat(absPath); err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("required file missing: %s", relPath))
			res.Valid = false
		}
	}

	validateRootArtifacts(plan, res)

	// 5. Compose file validation
	composePath := filepath.Join(outputDir, "docker-compose.yml")
	if plan.Target.Kind == "vps" {
		if _, err := os.Stat(composePath); err == nil {
			if err := validateCompose(composePath, plan, res); err != nil {
				return nil, err
			}
			validateComposeCLI(composePath, plan, res)
		}
	}

	// 6. .env.example validation
	envExPath := filepath.Join(outputDir, ".env.example")
	if _, err := os.Stat(envExPath); err == nil {
		validateEnvExample(envExPath, plan, res)
	}

	// 7. Check for unresolved placeholders {{ }}
	if err := checkPlaceholders(outputDir, res); err != nil {
		return nil, err
	}

	return res, nil
}

func validateRootArtifacts(plan *types.Plan, res *Result) {
	root := plan.Target.Root
	if root == "" {
		root = plan.Project.Root
	}
	if root == "" {
		root = "."
	}
	for _, artifact := range plan.Artifacts {
		if artifact.Scope != "root" {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, artifact.Path)); err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("required root artifact missing: %s", artifact.Path))
			res.Valid = false
		}
	}
}

func validateCompose(path string, plan *types.Plan, res *Result) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var composeData map[string]interface{}
	if err := yaml.Unmarshal(content, &composeData); err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("docker-compose.yml is not valid YAML: %v", err))
		res.Valid = false
		return nil
	}

	// Check for exposed database ports
	services, ok := composeData["services"].(map[string]interface{})
	if ok {
		for name, svc := range services {
			svcMap, ok := svc.(map[string]interface{})
			if !ok {
				continue
			}

			// Find if it's a database service in plan
			isDatabase := false
			for _, pSvc := range plan.Services {
				if pSvc.Name == name && pSvc.Type == "database" {
					isDatabase = true
					break
				}
			}

			if isDatabase {
				if _, ok := svcMap["ports"]; ok {
					res.Errors = append(res.Errors, fmt.Sprintf("security risk: database service '%s' has exposed ports", name))
					res.Valid = false
				}
			}
		}
	}

	return nil
}

func validateComposeCLI(path string, plan *types.Plan, res *Result) {
	if _, err := exec.LookPath("docker"); err != nil {
		res.Warnings = append(res.Warnings, "docker is not installed or not in PATH; skipped docker compose config")
		return
	}

	envPath := filepath.Join(plan.Target.Root, plan.Target.EnvFile)
	if plan.Target.EnvFile == "" {
		envPath = filepath.Join(plan.Target.Root, ".env")
	}
	if _, err := os.Stat(envPath); err != nil {
		res.Warnings = append(res.Warnings, fmt.Sprintf("%s missing; skipped docker compose config until env values exist", envPath))
		return
	}

	cmd := exec.Command("docker", "compose", "-f", path, "config")
	cmd.Dir = plan.Target.Root
	if out, err := cmd.CombinedOutput(); err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("docker compose config failed: %s", strings.TrimSpace(string(out))))
		res.Valid = false
	}
}

func validateEnvExample(path string, plan *types.Plan, res *Result) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	for _, key := range plan.Env.Required {
		if !strings.Contains(string(content), key+"=") {
			res.Errors = append(res.Errors, fmt.Sprintf(".env.example missing required key: %s", key))
			res.Valid = false
		}
	}
}

func checkPlaceholders(dir string, res *Result) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) == ".json" || strings.Contains(path, "backups") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Look for {{ but ignore if it's ${{
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(line, "{{") && !strings.Contains(line, "${{") {
				res.Errors = append(res.Errors, fmt.Sprintf("unresolved template placeholder in %s at line %d", path, i+1))
				res.Valid = false
			}
		}

		return nil
	})
}

// GenerateValidationReport produces a Markdown summary of the validation results.
func GenerateValidationReport(res *Result) string {
	var sb strings.Builder
	sb.WriteString("# Validation Report\n\n")

	if res.Valid {
		sb.WriteString("✅ **Status: Valid**\n\n")
		sb.WriteString("This deployment kit is safe to apply.\n\n")
	} else {
		sb.WriteString("❌ **Status: Invalid**\n\n")
		sb.WriteString("Please fix the following blocking errors before applying.\n\n")
	}

	if len(res.Errors) > 0 {
		sb.WriteString("### Errors\n")
		for _, err := range res.Errors {
			sb.WriteString(fmt.Sprintf("- %s\n", err))
		}
		sb.WriteString("\n")
	}

	if len(res.Warnings) > 0 {
		sb.WriteString("### Warnings\n")
		for _, warn := range res.Warnings {
			sb.WriteString(fmt.Sprintf("- %s\n", warn))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
