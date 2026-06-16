package scanner

import (
	"os"
	"path/filepath"

	"depctl/internal/types"
)

func scanCI(dir string) types.CIDetection {
	var det types.CIDetection

	if _, err := os.Stat(filepath.Join(dir, ".github", "workflows")); err == nil {
		det.GitHubActions = true
	}
	if _, err := os.Stat(filepath.Join(dir, ".gitlab-ci.yml")); err == nil {
		det.GitLab = true
	}
	if _, err := os.Stat(filepath.Join(dir, ".gitea", "workflows")); err == nil {
		det.GitEA = true
	}

	return det
}
