package planfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

// Load reads .deploy/plan.json and migrates older in-memory plans to v0.2.
func Load(outputDir string) (*types.Plan, error) {
	return LoadPath(filepath.Join(outputDir, "plan.json"), outputDir)
}

// LoadPath reads a plan from a specific path and migrates older in-memory plans to v0.2.
func LoadPath(path, outputDir string) (*types.Plan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan: %w", err)
	}

	var plan types.Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan: %w", err)
	}

	Migrate(&plan, outputDir)
	return &plan, nil
}

// Save writes the plan to outputDir/plan.json.
func Save(plan *types.Plan, outputDir string) error {
	Migrate(plan, outputDir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize plan: %w", err)
	}

	if err := os.WriteFile(filepath.Join(outputDir, "plan.json"), data, 0644); err != nil {
		return fmt.Errorf("failed to write plan: %w", err)
	}
	return nil
}

// Migrate normalizes v0.1 plans into the v0.2 target shape without changing behavior.
func Migrate(plan *types.Plan, outputDir string) {
	if plan.Version == "" || plan.Version == "0.1" {
		plan.Version = "0.2"
	}
	if plan.Target.Kind == "" {
		plan.Target.Kind = "vps"
	}
	if plan.Target.Preset == "" {
		plan.Target.Preset = plan.Preset
	}
	if plan.Target.Preset == "" {
		plan.Target.Preset = "compose-traefik"
	}
	if plan.Preset == "" {
		plan.Preset = plan.Target.Preset
	}
	if plan.Target.Root == "" {
		plan.Target.Root = plan.Project.Root
	}
	if plan.Target.OutputDir == "" {
		plan.Target.OutputDir = outputDir
	}
	if plan.Target.EnvFile == "" {
		plan.Target.EnvFile = ".env"
	}
	if plan.FileHashes == nil {
		plan.FileHashes = make(map[string]string)
	}
}
