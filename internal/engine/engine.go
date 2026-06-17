package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/types"
	"github.com/AnouarMohamed/Depctl/internal/validator"
)

// ApplyOptions holds flags and settings for the apply operation.
type ApplyOptions struct {
	OutputDir       string
	Yes             bool
	DryRun          bool
	SkipBuild       bool
	SkipHealthcheck bool
}

// Apply executes the deployment plan.
func Apply(opts ApplyOptions) error {
	// 1. Load plan
	planPath := filepath.Join(opts.OutputDir, "plan.json")
	planBytes, err := os.ReadFile(planPath)
	if err != nil {
		return fmt.Errorf("failed to read plan: %w", err)
	}

	var plan types.Plan
	if err := json.Unmarshal(planBytes, &plan); err != nil {
		return fmt.Errorf("failed to parse plan: %w", err)
	}

	// 2. Validate
	output.Step("Validating deployment kit...")
	res, err := validator.Validate(opts.OutputDir)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if !res.Valid {
		return fmt.Errorf("deployment kit is invalid, run 'depctl validate' for details")
	}

	// 3. Confirmation
	if !opts.Yes && !opts.DryRun {
		output.Warning("This will deploy %s to %s.", plan.Project.Name, plan.Domain)
		fmt.Print("Are you sure? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("apply cancelled by user")
		}
	}

	// 4. Create backup of existing state
	if !opts.DryRun {
		output.Step("Creating backup...")
		backupDir := filepath.Join(opts.OutputDir, "backups", time.Now().Format("20060102150405"))
		if err := os.MkdirAll(backupDir, 0755); err == nil {
			// Copy important files to backup
			_ = copyFile(filepath.Join(opts.OutputDir, "docker-compose.yml"), filepath.Join(backupDir, "docker-compose.yml"))
			_ = copyFile(filepath.Join(opts.OutputDir, "plan.json"), filepath.Join(backupDir, "plan.json"))
		}
	}

	// 5. Execute actions
	for _, action := range plan.Actions {
		if err := executeAction(action, opts, plan); err != nil {
			return err
		}
	}

	output.Success("Apply complete!")
	return nil
}

func executeAction(action types.Action, opts ApplyOptions, plan types.Plan) error {
	switch action.Type {
	case "create_network":
		return createNetwork(action.Name, opts)
	case "compose_up":
		return composeUp(action.File, opts)
	default:
		output.Warning("Unknown action type: %s", action.Type)
	}
	return nil
}

func createNetwork(name string, opts ApplyOptions) error {
	if opts.DryRun {
		output.Info("[DRY RUN] Would create Docker network: %s", name)
		return nil
	}

	output.Step("Creating Docker network: %s", name)
	// Check if network exists
	checkCmd := exec.Command("docker", "network", "inspect", name)
	if err := checkCmd.Run(); err == nil {
		output.Info("Network %s already exists, skipping.", name)
		return nil
	}

	cmd := exec.Command("docker", "network", "create", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create network %s: %s", name, string(out))
	}
	return nil
}

func composeUp(file string, opts ApplyOptions) error {
	if opts.DryRun {
		args := []string{"compose", "-f", file, "up", "-d"}
		if !opts.SkipBuild {
			args = append(args, "--build")
		}
		output.Info("[DRY RUN] Would run: docker %s", strings.Join(args, " "))
		return nil
	}

	output.Step("Running docker compose up...")
	args := []string{"compose", "-f", file, "up", "-d"}
	if !opts.SkipBuild {
		args = append(args, "--build")
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose up failed: %w", err)
	}
	return nil
}

// RollbackOptions holds flags and settings for the rollback operation.
type RollbackOptions struct {
	OutputDir string
	To        string
	DryRun    bool
	Yes       bool
}

// Rollback restores a previous deployment state from backup.
func Rollback(opts RollbackOptions) error {
	backupBaseDir := filepath.Join(opts.OutputDir, "backups")

	if opts.To == "" {
		// List backups
		entries, err := os.ReadDir(backupBaseDir)
		if err != nil {
			return fmt.Errorf("failed to list backups: %w", err)
		}

		output.Info("Available backups:")
		for _, entry := range entries {
			if entry.IsDir() {
				output.Step("- %s", entry.Name())
			}
		}
		return nil
	}

	targetBackupDir := filepath.Join(backupBaseDir, opts.To)
	if _, err := os.Stat(targetBackupDir); os.IsNotExist(err) {
		return fmt.Errorf("backup %s does not exist", opts.To)
	}

	if !opts.Yes && !opts.DryRun {
		output.Warning("This will restore deployment from backup %s.", opts.To)
		fmt.Print("Are you sure? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("rollback cancelled by user")
		}
	}

	if opts.DryRun {
		output.Info("[DRY RUN] Would restore files from %s", targetBackupDir)
		output.Info("[DRY RUN] Would run: docker compose -f .deploy/docker-compose.yml up -d")
		return nil
	}

	// Restore files
	output.Step("Restoring files from %s...", opts.To)
	filesToRestore := []string{"docker-compose.yml", "plan.json"}
	for _, f := range filesToRestore {
		src := filepath.Join(targetBackupDir, f)
		dst := filepath.Join(opts.OutputDir, f)
		if _, err := os.Stat(src); err == nil {
			if err := copyFile(src, dst); err != nil {
				output.Warning("Failed to restore %s: %v", f, err)
			}
		}
	}

	// Re-apply (compose up)
	output.Step("Restarting services...")
	cmd := exec.Command("docker", "compose", "-f", filepath.Join(opts.OutputDir, "docker-compose.yml"), "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart services during rollback: %w", err)
	}

	output.Success("Rollback complete!")
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
