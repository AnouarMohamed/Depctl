package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AnouarMohamed/Depctl/internal/planfile"
	"github.com/AnouarMohamed/Depctl/internal/types"
)

func TestValidateRequiresRootArtifacts(t *testing.T) {
	root := t.TempDir()
	out := filepath.Join(root, ".deploy")
	plan := &types.Plan{
		Version:        "0.2",
		Project:        types.ProjectPlan{Name: "app", Root: root},
		Target:         types.TargetPlan{Kind: "fly", Root: root, OutputDir: out, AppName: "app", Region: "iad"},
		Runtime:        types.RuntimePlan{Name: "go", Confidence: 0.9},
		GeneratedFiles: []string{".deploy/README.md"},
		Artifacts: []types.Artifact{{
			Path:  "fly.toml",
			Scope: "root",
		}},
	}
	if err := planfile.Save(plan, out); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "README.md"), []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}

	res, err := Validate(out)
	if err != nil {
		t.Fatal(err)
	}
	if res.Valid {
		t.Fatal("expected validation to fail without root artifact")
	}
	if len(res.Errors) == 0 {
		t.Fatal("expected validation errors")
	}
}

func TestValidateSkipsStaleComposeForProviderTargets(t *testing.T) {
	root := t.TempDir()
	out := filepath.Join(root, ".deploy")
	plan := &types.Plan{
		Version:        "0.2",
		Project:        types.ProjectPlan{Name: "app", Root: root},
		Target:         types.TargetPlan{Kind: "vercel", Root: root, OutputDir: out, AppName: "app"},
		Runtime:        types.RuntimePlan{Name: "node", Framework: "nextjs", Confidence: 0.9},
		GeneratedFiles: []string{".deploy/README.md"},
		Artifacts: []types.Artifact{{
			Path:  "vercel.json",
			Scope: "root",
		}},
	}
	if err := planfile.Save(plan, out); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "README.md"), []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "docker-compose.yml"), []byte("services: ["), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "vercel.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	res, err := Validate(out)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Valid {
		t.Fatalf("expected stale compose to be ignored for provider target, errors: %#v", res.Errors)
	}
}
