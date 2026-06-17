package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	planPath := filepath.Join(outputDir, "plan.json")
	planBytes, err := os.ReadFile(planPath)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("plan.json missing: %v", err))
		res.Valid = false
		return res, nil // Critical error, stop here
	}

	var plan types.Plan
	if err := json.Unmarshal(planBytes, &plan); err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("failed to parse plan.json: %v", err))
		res.Valid = false
		return res, nil // Critical error
	}

	// 2. Domain check
	if plan.Domain == "" {
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

	// 5. Compose file validation
	composePath := filepath.Join(outputDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err == nil {
		if err := validateCompose(composePath, &plan, res); err != nil {
			return nil, err
		}
	}

	// 6. .env.example validation
	envExPath := filepath.Join(outputDir, ".env.example")
	if _, err := os.Stat(envExPath); err == nil {
		validateEnvExample(envExPath, &plan, res)
	}

	// 7. Check for unresolved placeholders {{ }}
	if err := checkPlaceholders(outputDir, res); err != nil {
		return nil, err
	}

	return res, nil
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
