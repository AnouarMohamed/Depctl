package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

func TestWrite(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "depctl-writer-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	plan := &types.Plan{
		Project: types.ProjectPlan{Name: "test-project"},
		Runtime: types.RuntimePlan{Name: "node", Framework: "nextjs"},
		Build: types.BuildDetection{
			PackageManager: "npm",
		},
		Network: types.NetworkDetection{
			InternalPort: 3000,
		},
		Domain: "example.com",
		Services: []types.Service{
			{Name: "web", Type: "app", InternalPort: 3000},
		},
		GeneratedFiles: []string{
			".deploy/docker-compose.yml",
			".deploy/Dockerfile",
		},
	}

	err = Write(plan, tempDir, true)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Check if files exist
	expectedFiles := []string{
		"docker-compose.yml",
		"Dockerfile",
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(filepath.Join(tempDir, f)); os.IsNotExist(err) {
			t.Errorf("expected file %s does not exist", f)
		}
	}
}
