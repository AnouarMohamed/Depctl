package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"depctl/internal/types"
	"depctl/templates"
)

// Write generates the deployment kit files based on the provided plan.
func Write(plan *types.Plan, outputDir string, force bool) error {
	// 1. Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 2. Prepare backup directory if forcing overwrite
	var backupDir string
	if force {
		timestamp := time.Now().Format("20060102150405")
		backupDir = filepath.Join(outputDir, "backups", timestamp)
	}

	// 2. Map files to templates
	fileTemplateMap := map[string]string{
		".env.example":           "env/env.example.tmpl",
		"scripts/deploy.sh":      "scripts/deploy.sh.tmpl",
		"scripts/rollback.sh":    "scripts/rollback.sh.tmpl",
		"scripts/status.sh":      "scripts/status.sh.tmpl",
		"scripts/backup.sh":      "scripts/backup.sh.tmpl",
		"README.md":              "README.md.tmpl",
		"traefik/dynamic.yml":    "proxy/traefik-dynamic.yml.tmpl",
		"nginx/default.conf":     "proxy/nginx-default.conf.tmpl",
		"ci/github-actions.yml":  "ci/github-actions.yml.tmpl",
		".dockerignore":          "dockerfile/dockerignore.tmpl",
		".gitignore":             "deploy_gitignore.tmpl",
	}

	// Preset-based docker-compose template selection
	switch plan.Preset {
	case "compose-traefik":
		fileTemplateMap["docker-compose.yml"] = "compose/docker-compose.yml.tmpl"
	case "compose-nginx":
		fileTemplateMap["docker-compose.yml"] = "compose/nginx/docker-compose.yml.tmpl"
	default:
		fileTemplateMap["docker-compose.yml"] = "compose/docker-compose.yml.tmpl"
	}


	// Dockerfile selection
	if plan.Runtime.Name != "" {
		dockerfileTemplate := ""
		switch plan.Runtime.Name {
		case "node":
			if plan.Runtime.Framework == "nextjs" {
				dockerfileTemplate = "dockerfile/node-next.Dockerfile.tmpl"
			} else {
				dockerfileTemplate = "dockerfile/node-server.Dockerfile.tmpl"
			}
		case "laravel":
			dockerfileTemplate = "dockerfile/laravel.Dockerfile.tmpl"
		case "python":
			if plan.Runtime.Framework == "fastapi" {
				dockerfileTemplate = "dockerfile/python-fastapi.Dockerfile.tmpl"
			} else {
				dockerfileTemplate = "dockerfile/python-django.Dockerfile.tmpl"
			}
		case "go":
			dockerfileTemplate = "dockerfile/go.Dockerfile.tmpl"
		}

		if dockerfileTemplate != "" {
			fileTemplateMap["Dockerfile"] = dockerfileTemplate
		}
	}

	// 4. Render and write each file
	for _, relPath := range plan.GeneratedFiles {
		// Strip .deploy/ prefix if present
		cleanPath := relPath
		if len(relPath) > 8 && relPath[:8] == ".deploy/" {
			cleanPath = relPath[8:]
		}

		if cleanPath == "plan.json" {
			planBytes, err := json.MarshalIndent(plan, "", "  ")
			if err != nil {
				return err
			}
			err = os.WriteFile(filepath.Join(outputDir, cleanPath), planBytes, 0644)
			if err != nil {
				return err
			}
			continue
		}

		tmplPath, ok := fileTemplateMap[cleanPath]
		if !ok {
			// Some files might not have templates or are handled differently
			continue
		}

		if err := renderAndWrite(tmplPath, filepath.Join(outputDir, cleanPath), plan, force, backupDir); err != nil {
			return err
		}
	}

	return nil
}

func renderAndWrite(tmplPath, targetPath string, plan *types.Plan, force bool, backupDir string) error {
	// Check if file exists
	if _, err := os.Stat(targetPath); err == nil {
		if !force {
			return fmt.Errorf("file %s already exists, use --force to overwrite", targetPath)
		}

		// Backup existing file if backupDir is provided
		if backupDir != "" {
			if err := os.MkdirAll(backupDir, 0755); err == nil {
				content, err := os.ReadFile(targetPath)
				if err == nil {
					safeName := filepath.Base(targetPath)
					backupFile := filepath.Join(backupDir, safeName)
					_ = os.WriteFile(backupFile, content, 0644)
				}
			}
		}
	}

	// Read template
	tmplBytes, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmplPath, err)
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, plan); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", tmplPath, err)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
	}

	// Write file
	mode := os.FileMode(0644)
	if filepath.Ext(targetPath) == ".sh" {
		mode = 0755
	}

	if err := os.WriteFile(targetPath, buf.Bytes(), mode); err != nil {
		return fmt.Errorf("failed to write file %s: %w", targetPath, err)
	}

	return nil
}
