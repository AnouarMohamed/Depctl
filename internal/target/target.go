package target

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AnouarMohamed/Depctl/internal/envfile"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/planner"
	"github.com/AnouarMohamed/Depctl/internal/secrets"
	"github.com/AnouarMohamed/Depctl/internal/types"
	"github.com/AnouarMohamed/Depctl/internal/validator"
	"github.com/AnouarMohamed/Depctl/internal/writer"
)

// Support describes whether a target is a sensible fit for a detected project.
type Support struct {
	Supported bool
	Warnings  []string
}

// PlanOptions contains provider-aware planning inputs.
type PlanOptions struct {
	Preset    string
	Domain    string
	CI        string
	Target    string
	OutputDir string
	AppName   string
	Region    string
	EnvFile   string
}

// WriteOptions controls provider file rendering.
type WriteOptions struct {
	OutputDir string
	Force     bool
}

// ValidateOptions controls provider validation.
type ValidateOptions struct {
	OutputDir string
}

// ApplyOptions controls provider apply/deploy actions.
type ApplyOptions struct {
	OutputDir       string
	Yes             bool
	DryRun          bool
	SkipBuild       bool
	SkipHealthcheck bool
	EnvFile         string
}

// StatusOptions controls provider status output.
type StatusOptions struct {
	OutputDir string
}

// LogsOptions controls provider logs output.
type LogsOptions struct {
	OutputDir string
	Service   string
	Tail      int
}

// RollbackOptions controls provider rollback actions.
type RollbackOptions struct {
	OutputDir string
	To        string
	DryRun    bool
	Yes       bool
}

// Provider implements a deployment target.
type Provider interface {
	Kind() string
	DetectSupport(*types.Detection) Support
	Plan(*types.Detection, PlanOptions) (*types.Plan, error)
	Write(*types.Plan, WriteOptions) error
	Validate(*types.Plan, ValidateOptions) (*validator.Result, error)
	Apply(*types.Plan, ApplyOptions) error
	Status(*types.Plan, StatusOptions) error
	Logs(*types.Plan, LogsOptions) error
	Rollback(*types.Plan, RollbackOptions) error
}

var providers = map[string]Provider{
	"vps":    vpsProvider{},
	"vercel": vercelProvider{},
	"fly":    flyProvider{},
}

// Get returns a registered provider.
func Get(kind string) (Provider, error) {
	if kind == "" {
		kind = "vps"
	}
	provider, ok := providers[kind]
	if !ok {
		return nil, fmt.Errorf("unknown target %q; supported targets: vps, vercel, fly", kind)
	}
	return provider, nil
}

func planWithTarget(det *types.Detection, opts PlanOptions) (*types.Plan, error) {
	return planner.PlanWithOptions(det, planner.Options{
		Preset:    opts.Preset,
		Domain:    opts.Domain,
		CI:        opts.CI,
		Target:    opts.Target,
		OutputDir: opts.OutputDir,
		AppName:   opts.AppName,
		Region:    opts.Region,
		EnvFile:   opts.EnvFile,
	})
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

type commandSpec struct {
	Name         string
	Args         []string
	Cwd          string
	Stdin        string
	RedactValues []string
	AllowFailure bool
	Capture      bool
}

func runCommand(spec commandSpec, dryRun bool) (string, error) {
	cmdLine := strings.Join(append([]string{spec.Name}, spec.Args...), " ")
	cmdLine = secrets.Redact(cmdLine, spec.RedactValues)
	if dryRun {
		output.Info("[DRY RUN] Would run: %s", cmdLine)
		return "", nil
	}

	cmd := exec.Command(spec.Name, spec.Args...)
	cmd.Dir = spec.Cwd
	if spec.Stdin != "" {
		cmd.Stdin = strings.NewReader(spec.Stdin)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	combined := stdout.String() + stderr.String()
	combined = secrets.Redact(combined, spec.RedactValues)
	if combined != "" {
		fmt.Print(combined)
	}

	if err != nil && !spec.AllowFailure {
		return combined, fmt.Errorf("%s failed: %w", cmdLine, err)
	}
	if err != nil {
		output.Warning("Ignoring non-blocking command failure: %s", cmdLine)
		return combined, err
	}
	return combined, nil
}

func rootDir(plan *types.Plan) string {
	if plan.Target.Root != "" {
		return plan.Target.Root
	}
	if plan.Project.Root != "" {
		return plan.Project.Root
	}
	return "."
}

func outputDir(plan *types.Plan, fallback string) string {
	if fallback != "" {
		return fallback
	}
	if plan.Target.OutputDir != "" {
		return plan.Target.OutputDir
	}
	return ".deploy"
}

func envFilePath(plan *types.Plan, override string) string {
	envPath := override
	if envPath == "" {
		envPath = plan.Target.EnvFile
	}
	if envPath == "" {
		envPath = ".env"
	}
	if filepath.IsAbs(envPath) {
		return envPath
	}
	return filepath.Join(rootDir(plan), envPath)
}

func providerEnv(plan *types.Plan, override string) ([]envfile.Entry, []string, error) {
	path := envFilePath(plan, override)
	entries, err := envfile.Parse(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			output.Warning("Env file %s does not exist; skipping provider secret import.", path)
			return nil, nil, nil
		}
		return nil, nil, err
	}

	values := make([]string, 0, len(entries))
	for _, entry := range entries {
		values = append(values, entry.Value)
	}
	return entries, values, nil
}

func statePath(plan *types.Plan, outDir, name string) string {
	return filepath.Join(outputDir(plan, outDir), "state", name+".json")
}

func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func confirmApply(plan *types.Plan, yes bool, dryRun bool) error {
	if yes || dryRun {
		return nil
	}
	target := plan.Target.Kind
	if target == "" {
		target = "vps"
	}
	output.Warning("This will deploy %s using target %s.", plan.Project.Name, target)
	if plan.Domain != "" {
		output.Warning("Target domain: %s", plan.Domain)
	}
	fmt.Print("Are you sure? [y/N]: ")
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" {
		return fmt.Errorf("apply cancelled by user")
	}
	return nil
}

func imageLabel(plan *types.Plan) string {
	root := rootDir(plan)
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = root
	out, err := cmd.Output()
	if err == nil {
		sha := strings.TrimSpace(string(out))
		if sha != "" {
			return "depctl-" + sha
		}
	}
	return "depctl-" + time.Now().Format("20060102150405")
}

func writeProviderFiles(plan *types.Plan, opts WriteOptions, includeDeployKit bool) error {
	out := outputDir(plan, opts.OutputDir)
	plan.Target.OutputDir = out
	if includeDeployKit {
		if err := writer.Write(plan, out, opts.Force); err != nil {
			return err
		}
	}
	return writer.WriteRootArtifacts(plan, opts.Force)
}
