package planfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMigratesV01PlanToVPS(t *testing.T) {
	dir := t.TempDir()
	planJSON := `{
  "version": "0.1",
  "project": {"name": "app", "root": "/srv/app"},
  "preset": "compose-traefik",
  "domain": "app.example.com",
  "runtime": {"name": "node", "framework": "express", "confidence": 0.9}
}`
	if err := os.WriteFile(filepath.Join(dir, "plan.json"), []byte(planJSON), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	if plan.Version != "0.2" {
		t.Fatalf("version: got %q", plan.Version)
	}
	if plan.Target.Kind != "vps" {
		t.Fatalf("target kind: got %q", plan.Target.Kind)
	}
	if plan.Target.Preset != "compose-traefik" {
		t.Fatalf("target preset: got %q", plan.Target.Preset)
	}
	if plan.Target.OutputDir != dir {
		t.Fatalf("output dir: got %q", plan.Target.OutputDir)
	}
	if plan.Target.EnvFile != ".env" {
		t.Fatalf("env file: got %q", plan.Target.EnvFile)
	}
}
