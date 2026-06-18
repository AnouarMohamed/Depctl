package target

import (
	"fmt"
	"os"

	"github.com/AnouarMohamed/Depctl/internal/envfile"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/types"
	"github.com/AnouarMohamed/Depctl/internal/validator"
)

type flyProvider struct{}

type flyState struct {
	LastImage     string `json:"lastImage,omitempty"`
	PreviousImage string `json:"previousImage,omitempty"`
	AppName       string `json:"appName"`
}

func (flyProvider) Kind() string { return "fly" }

func (flyProvider) DetectSupport(det *types.Detection) Support {
	if det.Runtime.Name == "unknown" {
		return Support{Supported: false, Warnings: []string{"Fly target needs a Dockerfile-capable app; runtime detection is unknown"}}
	}
	return Support{Supported: true}
}

func (flyProvider) Plan(det *types.Detection, opts PlanOptions) (*types.Plan, error) {
	opts.Target = "fly"
	return planWithTarget(det, opts)
}

func (flyProvider) Write(plan *types.Plan, opts WriteOptions) error {
	return writeProviderFiles(plan, opts, true)
}

func (flyProvider) Validate(plan *types.Plan, opts ValidateOptions) (*validator.Result, error) {
	res, _ := validator.Validate(outputDir(plan, opts.OutputDir))
	for _, path := range []string{"Dockerfile", "fly.toml"} {
		if _, err := os.Stat(pathJoinRoot(plan, path)); err != nil {
			res.Valid = false
			res.Errors = append(res.Errors, fmt.Sprintf("root %s missing", path))
		}
	}
	if !commandExists("fly") {
		res.Warnings = append(res.Warnings, "fly CLI is not installed or not in PATH")
	}
	return res, nil
}

func (flyProvider) Apply(plan *types.Plan, opts ApplyOptions) error {
	if err := confirmApply(plan, opts.Yes, opts.DryRun); err != nil {
		return err
	}
	if !commandExists("fly") && !opts.DryRun {
		return fmt.Errorf("fly CLI is not installed or not in PATH")
	}

	token := os.Getenv("FLY_ACCESS_TOKEN")
	redact := []string{token}
	if err := ensureFlyAuth(token, opts.DryRun, redact); err != nil {
		return err
	}
	app := plan.Target.AppName
	if app == "" {
		app = plan.Project.Name
	}

	output.Step("Ensuring Fly app exists...")
	statusArgs := flyArgs(token, "status", "-a", app)
	_, statusErr := runCommand(commandSpec{Name: "fly", Args: statusArgs, Cwd: rootDir(plan), RedactValues: redact, AllowFailure: true}, opts.DryRun)
	if statusErr != nil || opts.DryRun {
		launchArgs := flyArgs(token, "launch", "--no-deploy", "--copy-config", "--yes", "--name", app, "--primary-region", plan.Target.Region, "--internal-port", fmt.Sprintf("%d", plan.Network.InternalPort))
		_, err := runCommand(commandSpec{Name: "fly", Args: launchArgs, Cwd: rootDir(plan), RedactValues: redact}, opts.DryRun)
		if err != nil {
			return err
		}
	}

	entries, values, err := providerEnv(plan, opts.EnvFile)
	if err != nil {
		return err
	}
	redact = append(redact, values...)
	if len(entries) > 0 {
		output.Step("Importing Fly secrets from env file...")
		_, err := runCommand(commandSpec{
			Name:         "fly",
			Args:         flyArgs(token, "secrets", "import", "--stage", "-a", app),
			Cwd:          rootDir(plan),
			Stdin:        envfile.AsDotenv(entries),
			RedactValues: redact,
		}, opts.DryRun)
		if err != nil {
			return err
		}
	}

	label := imageLabel(plan)
	newImage := fmt.Sprintf("registry.fly.io/%s:%s", app, label)
	output.Step("Deploying to Fly...")
	deployArgs := flyArgs(token, "deploy", "-a", app, "--image-label", label)
	if opts.SkipBuild {
		output.Warning("--skip-build is not supported for Fly source deployments; continuing with fly deploy.")
	}
	_, err = runCommand(commandSpec{Name: "fly", Args: deployArgs, Cwd: rootDir(plan), RedactValues: redact}, opts.DryRun)
	if err != nil {
		return err
	}

	if plan.Domain != "" {
		output.Step("Attaching Fly custom domain %s...", plan.Domain)
		_, err := runCommand(commandSpec{Name: "fly", Args: flyArgs(token, "certs", "add", plan.Domain, "-a", app), Cwd: rootDir(plan), RedactValues: redact}, opts.DryRun)
		if err != nil {
			return err
		}
	}

	if !opts.DryRun {
		var state flyState
		_ = readJSON(statePath(plan, opts.OutputDir, "fly"), &state)
		state.PreviousImage = state.LastImage
		state.LastImage = newImage
		state.AppName = app
		_ = writeJSON(statePath(plan, opts.OutputDir, "fly"), state)
	}

	output.Success("Fly deployment complete.")
	return nil
}

func (flyProvider) Status(plan *types.Plan, opts StatusOptions) error {
	token := os.Getenv("FLY_ACCESS_TOKEN")
	if err := ensureFlyAuth(token, false, []string{token}); err != nil {
		return err
	}
	_, err := runCommand(commandSpec{Name: "fly", Args: flyArgs(token, "status", "-a", plan.Target.AppName), Cwd: rootDir(plan), RedactValues: []string{token}}, false)
	return err
}

func (flyProvider) Logs(plan *types.Plan, opts LogsOptions) error {
	token := os.Getenv("FLY_ACCESS_TOKEN")
	if err := ensureFlyAuth(token, false, []string{token}); err != nil {
		return err
	}
	args := flyArgs(token, "logs", "-a", plan.Target.AppName)
	if opts.Tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", opts.Tail))
	}
	_, err := runCommand(commandSpec{Name: "fly", Args: args, Cwd: rootDir(plan), RedactValues: []string{token}}, false)
	return err
}

func (flyProvider) Rollback(plan *types.Plan, opts RollbackOptions) error {
	token := os.Getenv("FLY_ACCESS_TOKEN")
	if err := ensureFlyAuth(token, opts.DryRun, []string{token}); err != nil {
		return err
	}
	app := plan.Target.AppName
	image := opts.To
	if image == "" {
		var state flyState
		if err := readJSON(statePath(plan, opts.OutputDir, "fly"), &state); err != nil {
			return fmt.Errorf("no Fly rollback state found; provide --to with an image reference")
		}
		image = state.PreviousImage
		if app == "" {
			app = state.AppName
		}
	}
	if image == "" {
		return fmt.Errorf("no previous Fly image recorded; provide --to with an image reference")
	}
	_, err := runCommand(commandSpec{Name: "fly", Args: flyArgs(token, "deploy", "-a", app, "--image", image), Cwd: rootDir(plan), RedactValues: []string{token}}, opts.DryRun)
	return err
}

func flyArgs(token string, args ...string) []string {
	if token == "" {
		return args
	}
	return append([]string{"-t", token}, args...)
}

func ensureFlyAuth(token string, dryRun bool, redact []string) error {
	if token != "" {
		return nil
	}
	if dryRun {
		output.Info("[DRY RUN] Would check Fly login and run 'fly auth login' if needed")
		return nil
	}
	if !commandExists("fly") {
		return fmt.Errorf("fly CLI is not installed or not in PATH")
	}

	if _, err := runCommand(commandSpec{Name: "fly", Args: []string{"auth", "whoami"}, RedactValues: redact, AllowFailure: true}, false); err == nil {
		return nil
	}
	if !hasInteractiveTerminal() {
		return fmt.Errorf("fly is not logged in; run 'fly auth login' locally or set FLY_ACCESS_TOKEN for CI")
	}

	output.Step("Opening Fly login...")
	_, err := runCommand(commandSpec{Name: "fly", Args: []string{"auth", "login"}, Interactive: true}, false)
	return err
}

func pathJoinRoot(plan *types.Plan, path string) string {
	return rootDir(plan) + string(os.PathSeparator) + path
}
