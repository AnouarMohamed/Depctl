package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AnouarMohamed/Depctl/internal/planner"
	"github.com/AnouarMohamed/Depctl/internal/scanner"
	"github.com/AnouarMohamed/Depctl/internal/writer"
)

func TestGoldenFlows(t *testing.T) {
	fixturesDir := "../../fixtures"
	fixtures, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("failed to read fixtures dir: %v", err)
	}

	presets := []string{"compose-traefik", "compose-nginx"}

	for _, fixture := range fixtures {
		if !fixture.IsDir() {
			continue
		}

		for _, preset := range presets {
			t.Run(fixture.Name()+"/"+preset, func(t *testing.T) {
				fixturePath, _ := filepath.Abs(filepath.Join(fixturesDir, fixture.Name()))
				
				// 1. Scan
				det, err := scanner.Scan(fixturePath)
				if err != nil {
					t.Fatalf("Scan failed: %v", err)
				}

				// 2. Plan
				plan, err := planner.Plan(det, preset, "example.com", "github")
				if err != nil {
					t.Fatalf("Plan failed: %v", err)
				}
				
				// Ensure plan.json is in the list for the writer to pick it up
				plan.GeneratedFiles = append(plan.GeneratedFiles, ".deploy/plan.json")

				// 3. Write to temp dir
				tempOutputDir, err := os.MkdirTemp("", "depctl-golden-*")
				if err != nil {
					t.Fatalf("failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tempOutputDir)

				err = writer.Write(plan, tempOutputDir, true)
				if err != nil {
					t.Fatalf("Write failed: %v", err)
				}

				// 4. Verify key files exist
				expectedFiles := []string{"docker-compose.yml", "plan.json", "README.md"}
				for _, f := range expectedFiles {
					if _, err := os.Stat(filepath.Join(tempOutputDir, f)); os.IsNotExist(err) {
						t.Errorf("expected file %s missing for %s/%s", f, fixture.Name(), preset)
					}
				}
				
				if preset == "compose-traefik" {
					if _, err := os.Stat(filepath.Join(tempOutputDir, "traefik", "dynamic.yml")); os.IsNotExist(err) {
						t.Errorf("traefik/dynamic.yml missing")
					}
				} else if preset == "compose-nginx" {
					if _, err := os.Stat(filepath.Join(tempOutputDir, "nginx", "default.conf")); os.IsNotExist(err) {
						t.Errorf("nginx/default.conf missing")
					}
				}
			})
		}
	}
}
