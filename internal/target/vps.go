package target

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/types"
	"github.com/AnouarMohamed/Depctl/internal/validator"
	"github.com/AnouarMohamed/Depctl/internal/writer"
)

type vpsProvider struct{}

func (vpsProvider) Kind() string { return "vps" }

func (vpsProvider) DetectSupport(det *types.Detection) Support {
	if det.Runtime.Name == "unknown" {
		return Support{Supported: false, Warnings: []string{"unknown runtime cannot be deployed safely to VPS without manual Dockerfile review"}}
	}
	return Support{Supported: true}
}

func (vpsProvider) Plan(det *types.Detection, opts PlanOptions) (*types.Plan, error) {
	opts.Target = "vps"
	return planWithTarget(det, opts)
}

func (vpsProvider) Write(plan *types.Plan, opts WriteOptions) error {
	out := outputDir(plan, opts.OutputDir)
	plan.Target.OutputDir = out
	if err := writer.Write(plan, out, opts.Force); err != nil {
		return err
	}
	return writer.WriteRootArtifacts(plan, opts.Force)
}

func (vpsProvider) Validate(plan *types.Plan, opts ValidateOptions) (*validator.Result, error) {
	return validator.Validate(outputDir(plan, opts.OutputDir))
}

func (vpsProvider) Apply(plan *types.Plan, opts ApplyOptions) error {
	if err := confirmApply(plan, opts.Yes, opts.DryRun); err != nil {
		return err
	}

	res, err := validator.Validate(outputDir(plan, opts.OutputDir))
	if err != nil {
		return err
	}
	if !res.Valid {
		return fmt.Errorf("deployment kit is invalid, run 'depctl validate' for details")
	}

	if !opts.DryRun {
		if err := backupDeployState(plan, opts.OutputDir); err != nil {
			return err
		}
	}

	for _, action := range plan.Actions {
		switch action.Type {
		case "create_network":
			if err := createNetwork(action.Name, opts.DryRun); err != nil {
				return err
			}
		case "create_swarm_network":
			if err := createSwarmNetwork(action.Name, opts.DryRun); err != nil {
				return err
			}
		case "compose_up":
			if err := composeUp(plan, action.File, opts); err != nil {
				return err
			}
		case "swarm_deploy":
			if err := swarmDeploy(plan, action.File, opts); err != nil {
				return err
			}
		default:
			output.Warning("Unknown VPS action type: %s", action.Type)
		}
	}

	output.Success("VPS apply complete.")
	return nil
}

func (vpsProvider) Status(plan *types.Plan, opts StatusOptions) error {
	if plan.Preset == "swarm-traefik" {
		_, err := runCommand(commandSpec{Name: "docker", Args: []string{"stack", "services", plan.Project.Name}, Cwd: rootDir(plan)}, false)
		return err
	}
	composePath := filepath.Join(outputDir(plan, opts.OutputDir), "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		return fmt.Errorf("compose file missing: %s", composePath)
	}
	_, err := runCommand(commandSpec{Name: "docker", Args: []string{"compose", "-f", composePath, "ps"}, Cwd: rootDir(plan)}, false)
	return err
}

func (vpsProvider) Logs(plan *types.Plan, opts LogsOptions) error {
	if plan.Preset == "swarm-traefik" {
		serviceName := opts.Service
		if serviceName == "" {
			serviceName = plan.PublicService
		}
		fullServiceName := fmt.Sprintf("%s_%s", plan.Project.Name, serviceName)
		args := []string{"service", "logs"}
		if opts.Tail > 0 {
			args = append(args, "--tail", fmt.Sprintf("%d", opts.Tail))
		}
		args = append(args, fullServiceName)
		_, err := runCommand(commandSpec{Name: "docker", Args: args, Cwd: rootDir(plan)}, false)
		return err
	}
	composePath := filepath.Join(outputDir(plan, opts.OutputDir), "docker-compose.yml")
	args := []string{"compose", "-f", composePath, "logs"}
	if opts.Tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", opts.Tail))
	}
	if opts.Service != "" {
		args = append(args, opts.Service)
	}
	_, err := runCommand(commandSpec{Name: "docker", Args: args, Cwd: rootDir(plan)}, false)
	return err
}

func (vpsProvider) Rollback(plan *types.Plan, opts RollbackOptions) error {
	backupBaseDir := filepath.Join(outputDir(plan, opts.OutputDir), "backups")
	if opts.To == "" {
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
	if _, err := os.Stat(targetBackupDir); err != nil {
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
		output.Info("[DRY RUN] Would run docker compose up")
		return nil
	}

	for _, name := range []string{"docker-compose.yml", "plan.json"} {
		src := filepath.Join(targetBackupDir, name)
		dst := filepath.Join(outputDir(plan, opts.OutputDir), name)
		_ = copyFile(src, dst)
	}
	return composeUp(plan, filepath.Join(outputDir(plan, opts.OutputDir), "docker-compose.yml"), ApplyOptions{SkipBuild: true, Yes: true})
}

func backupDeployState(plan *types.Plan, outDir string) error {
	backupDir := filepath.Join(outputDir(plan, outDir), "backups", time.Now().Format("20060102150405"))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}
	for _, name := range []string{"docker-compose.yml", "plan.json"} {
		src := filepath.Join(outputDir(plan, outDir), name)
		if _, err := os.Stat(src); err == nil {
			if err := copyFile(src, filepath.Join(backupDir, name)); err != nil {
				output.Warning("Failed to backup %s: %v", name, err)
			}
		}
	}
	return nil
}

func createNetwork(name string, dryRun bool) error {
	if dryRun {
		output.Info("[DRY RUN] Would create Docker network: %s", name)
		return nil
	}
	if name == "" {
		return nil
	}

	output.Step("Ensuring Docker network exists: %s", name)
	if err := exec.Command("docker", "network", "inspect", name).Run(); err == nil {
		output.Info("Network %s already exists.", name)
		return nil
	}
	_, err := runCommand(commandSpec{Name: "docker", Args: []string{"network", "create", name}}, false)
	return err
}

func composeUp(plan *types.Plan, file string, opts ApplyOptions) error {
	if file == "" {
		file = filepath.Join(outputDir(plan, opts.OutputDir), "docker-compose.yml")
	}
	if file == filepath.Join(".deploy", "docker-compose.yml") {
		file = filepath.Join(outputDir(plan, opts.OutputDir), "docker-compose.yml")
	}
	args := []string{"compose", "-f", file, "up", "-d"}
	if !opts.SkipBuild {
		args = append(args, "--build")
	}
	_, err := runCommand(commandSpec{Name: "docker", Args: args, Cwd: rootDir(plan)}, opts.DryRun)
	if err != nil && !opts.DryRun {
		_, _ = runCommand(commandSpec{Name: "docker", Args: []string{"compose", "-f", file, "logs", "--tail", "80"}, Cwd: rootDir(plan), AllowFailure: true}, false)
	}
	return err
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func createSwarmNetwork(name string, dryRun bool) error {
	if dryRun {
		output.Info("[DRY RUN] Would create Swarm overlay network: %s", name)
		return nil
	}
	if name == "" {
		return nil
	}

	output.Step("Ensuring Swarm overlay network exists: %s", name)
	if err := exec.Command("docker", "network", "inspect", name).Run(); err == nil {
		output.Info("Network %s already exists.", name)
		return nil
	}
	_, err := runCommand(commandSpec{Name: "docker", Args: []string{"network", "create", "--driver", "overlay", "--attachable", name}}, false)
	return err
}

func swarmDeploy(plan *types.Plan, file string, opts ApplyOptions) error {
	if file == "" {
		file = filepath.Join(outputDir(plan, opts.OutputDir), "docker-compose.yml")
	}
	if file == filepath.Join(".deploy", "docker-compose.yml") {
		file = filepath.Join(outputDir(plan, opts.OutputDir), "docker-compose.yml")
	}

	if !opts.DryRun {
		out, err := exec.Command("docker", "info", "--format", "{{.Swarm.LocalNodeState}}").Output()
		if err != nil || strings.TrimSpace(string(out)) != "active" {
			return fmt.Errorf("Docker Swarm is not initialized on this machine. Run 'docker swarm init' first")
		}
	}

	// 1. Build image locally since Swarm doesn't build images
	if !opts.SkipBuild {
		output.Step("Building image locally for Swarm...")
		_, err := runCommand(commandSpec{Name: "docker", Args: []string{"compose", "-f", file, "build"}, Cwd: rootDir(plan)}, opts.DryRun)
		if err != nil {
			return err
		}
	}

	// 2. Deploy stack
	output.Step("Deploying Docker Swarm stack: %s...", plan.Project.Name)
	_, err := runCommand(commandSpec{Name: "docker", Args: []string{"stack", "deploy", "-c", file, plan.Project.Name}, Cwd: rootDir(plan)}, opts.DryRun)
	return err
}
