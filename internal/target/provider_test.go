package target

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

func TestProviderWriteCreatesExpectedRootArtifacts(t *testing.T) {
	root := t.TempDir()
	out := filepath.Join(root, ".deploy")
	det := &types.Detection{
		Project: types.ProjectDetection{Name: "next-app", Root: root},
		Runtime: types.RuntimeDetection{Name: "node", Framework: "nextjs", Confidence: 0.95},
		Build:   types.BuildDetection{PackageManager: "npm", BuildCommand: "npm run build", StartCommand: "npm start"},
		Network: types.NetworkDetection{InternalPort: 3000, Confidence: 0.8},
		Dependencies: map[string]types.Dependency{
			"postgres": {Likely: false},
			"redis":    {Likely: false},
		},
	}

	provider, err := Get("vercel")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := provider.Plan(det, PlanOptions{Target: "vercel", OutputDir: out})
	if err != nil {
		t.Fatal(err)
	}
	if err := provider.Write(plan, WriteOptions{OutputDir: out, Force: true}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(root, "vercel.json")); err != nil {
		t.Fatalf("expected vercel.json root artifact: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "plan.json")); err != nil {
		t.Fatalf("expected plan.json deploy artifact: %v", err)
	}
}
