package engine

import (
	"fmt"

	"github.com/AnouarMohamed/Depctl/internal/planfile"
	"github.com/AnouarMohamed/Depctl/internal/target"
)

// ApplyOptions holds flags and settings for the apply operation.
type ApplyOptions struct {
	OutputDir       string
	Yes             bool
	DryRun          bool
	SkipBuild       bool
	SkipHealthcheck bool
	EnvFile         string
}

// Apply executes the deployment plan through its target provider.
func Apply(opts ApplyOptions) error {
	plan, err := planfile.Load(opts.OutputDir)
	if err != nil {
		return err
	}
	provider, err := target.Get(plan.Target.Kind)
	if err != nil {
		return err
	}
	return provider.Apply(plan, target.ApplyOptions{
		OutputDir:       opts.OutputDir,
		Yes:             opts.Yes,
		DryRun:          opts.DryRun,
		SkipBuild:       opts.SkipBuild,
		SkipHealthcheck: opts.SkipHealthcheck,
		EnvFile:         opts.EnvFile,
	})
}

// Status shows provider-specific deployment status.
func Status(outputDir string) error {
	plan, err := planfile.Load(outputDir)
	if err != nil {
		return err
	}
	provider, err := target.Get(plan.Target.Kind)
	if err != nil {
		return err
	}
	return provider.Status(plan, target.StatusOptions{OutputDir: outputDir})
}

// Logs shows provider-specific deployment logs.
func Logs(outputDir, service string, tail int) error {
	plan, err := planfile.Load(outputDir)
	if err != nil {
		return err
	}
	provider, err := target.Get(plan.Target.Kind)
	if err != nil {
		return err
	}
	return provider.Logs(plan, target.LogsOptions{OutputDir: outputDir, Service: service, Tail: tail})
}

// RollbackOptions holds flags and settings for the rollback operation.
type RollbackOptions struct {
	OutputDir string
	To        string
	DryRun    bool
	Yes       bool
}

// Rollback restores a previous deployment state through its target provider.
func Rollback(opts RollbackOptions) error {
	plan, err := planfile.Load(opts.OutputDir)
	if err != nil {
		return err
	}
	provider, err := target.Get(plan.Target.Kind)
	if err != nil {
		return err
	}
	if plan.Target.Kind != "vps" && opts.To == "" && plan.Rollback.Strategy == "" {
		return fmt.Errorf("rollback target requires --to")
	}
	return provider.Rollback(plan, target.RollbackOptions{
		OutputDir: opts.OutputDir,
		To:        opts.To,
		DryRun:    opts.DryRun,
		Yes:       opts.Yes,
	})
}
