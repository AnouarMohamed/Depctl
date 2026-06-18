package writer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

func TestWriteRootArtifactsBacksUpBeforeOverwrite(t *testing.T) {
	root := t.TempDir()
	out := filepath.Join(root, ".deploy")
	if err := os.MkdirAll(out, 0755); err != nil {
		t.Fatal(err)
	}
	dockerfile := filepath.Join(root, "Dockerfile")
	if err := os.WriteFile(dockerfile, []byte("FROM scratch\n"), 0644); err != nil {
		t.Fatal(err)
	}

	plan := &types.Plan{
		Project: types.ProjectPlan{Name: "app", Root: root},
		Target:  types.TargetPlan{Kind: "fly", Root: root, OutputDir: out, AppName: "app", Region: "iad"},
		Runtime: types.RuntimePlan{Name: "go", Confidence: 0.95},
		Network: types.NetworkDetection{InternalPort: 8080},
		Build:   types.BuildDetection{StartCommand: "./app"},
		Artifacts: []types.Artifact{{
			Path:  "Dockerfile",
			Kind:  "dockerfile",
			Scope: "root",
			Mode:  "0644",
		}},
	}

	if err := WriteRootArtifacts(plan, false); err == nil {
		t.Fatal("expected overwrite without force to fail")
	}
	if err := WriteRootArtifacts(plan, true); err != nil {
		t.Fatal(err)
	}

	backups, err := filepath.Glob(filepath.Join(out, "backups", "*", "root", "Dockerfile"))
	if err != nil {
		t.Fatal(err)
	}
	if len(backups) != 1 {
		t.Fatalf("expected one root artifact backup, got %d", len(backups))
	}
	data, err := os.ReadFile(backups[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "FROM scratch\n" {
		t.Fatalf("backup content changed: %q", string(data))
	}
}
