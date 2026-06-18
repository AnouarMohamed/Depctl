package target

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/types"
	"github.com/AnouarMohamed/Depctl/internal/validator"
)

type vercelProvider struct{}

type vercelState struct {
	DeploymentURL string `json:"deploymentUrl"`
	Domain        string `json:"domain,omitempty"`
}

func (vercelProvider) Kind() string { return "vercel" }

func (vercelProvider) DetectSupport(det *types.Detection) Support {
	ok := det.Runtime.Name == "node" && (det.Runtime.Framework == "nextjs" || det.Runtime.Framework == "vite")
	if ok {
		return Support{Supported: true}
	}
	return Support{Supported: false, Warnings: []string{"Vercel target is optimized here for Next.js and Vite/static apps"}}
}

func (vercelProvider) Plan(det *types.Detection, opts PlanOptions) (*types.Plan, error) {
	opts.Target = "vercel"
	return planWithTarget(det, opts)
}

func (vercelProvider) Write(plan *types.Plan, opts WriteOptions) error {
	return writeProviderFiles(plan, opts, true)
}

func (vercelProvider) Validate(plan *types.Plan, opts ValidateOptions) (*validator.Result, error) {
	res, _ := validator.Validate(outputDir(plan, opts.OutputDir))
	root := rootDir(plan)
	if _, err := os.Stat(filepath.Join(root, "vercel.json")); err != nil {
		res.Valid = false
		res.Errors = append(res.Errors, "root vercel.json missing")
	}
	if plan.Runtime.Name != "node" || (plan.Runtime.Framework != "nextjs" && plan.Runtime.Framework != "vite") {
		res.Warnings = append(res.Warnings, "vercel target is best for Next.js or Vite/static projects")
	}
	if !commandExists("vercel") {
		res.Warnings = append(res.Warnings, "vercel CLI is not installed or not in PATH")
	}
	return res, nil
}

func (vercelProvider) Apply(plan *types.Plan, opts ApplyOptions) error {
	if err := confirmApply(plan, opts.Yes, opts.DryRun); err != nil {
		return err
	}
	if !commandExists("vercel") && !opts.DryRun {
		return fmt.Errorf("vercel CLI is not installed or not in PATH")
	}

	token := os.Getenv("VERCEL_TOKEN")
	redact := []string{token}
	global := vercelGlobalArgs(token)
	if err := ensureVercelAuth(token, opts.DryRun, redact); err != nil {
		return err
	}

	output.Step("Linking Vercel project...")
	_, _ = runCommand(commandSpec{Name: "vercel", Args: append([]string{"link", "--yes"}, global...), Cwd: rootDir(plan), RedactValues: redact, AllowFailure: true}, opts.DryRun)

	entries, values, err := providerEnv(plan, opts.EnvFile)
	if err != nil {
		return err
	}
	redact = append(redact, values...)
	for _, entry := range entries {
		if entry.Value == "" {
			continue
		}
		_, _ = runCommand(commandSpec{Name: "vercel", Args: append([]string{"env", "rm", entry.Key, "production", "--yes"}, global...), Cwd: rootDir(plan), RedactValues: redact, AllowFailure: true}, opts.DryRun)
		_, err := runCommand(commandSpec{
			Name:         "vercel",
			Args:         append([]string{"env", "add", entry.Key, "production"}, global...),
			Cwd:          rootDir(plan),
			Stdin:        entry.Value + "\n",
			RedactValues: redact,
		}, opts.DryRun)
		if err != nil {
			return err
		}
	}

	output.Step("Deploying to Vercel production...")
	out, err := runCommand(commandSpec{Name: "vercel", Args: append([]string{"--prod", "--yes"}, global...), Cwd: rootDir(plan), RedactValues: redact, Capture: true}, opts.DryRun)
	if err != nil {
		return err
	}

	deploymentURL := lastURL(out)
	if plan.Domain != "" && deploymentURL != "" {
		output.Step("Assigning Vercel alias %s...", plan.Domain)
		_, err := runCommand(commandSpec{Name: "vercel", Args: append([]string{"alias", "set", deploymentURL, plan.Domain}, global...), Cwd: rootDir(plan), RedactValues: redact}, opts.DryRun)
		if err != nil {
			return err
		}
	}
	if !opts.DryRun {
		_ = writeJSON(statePath(plan, opts.OutputDir, "vercel"), vercelState{DeploymentURL: deploymentURL, Domain: plan.Domain})
	}
	output.Success("Vercel deployment complete.")
	return nil
}

func (vercelProvider) Status(plan *types.Plan, opts StatusOptions) error {
	token := os.Getenv("VERCEL_TOKEN")
	if err := ensureVercelAuth(token, false, []string{token}); err != nil {
		return err
	}
	args := append([]string{"ls"}, vercelGlobalArgs(token)...)
	_, err := runCommand(commandSpec{Name: "vercel", Args: args, Cwd: rootDir(plan), RedactValues: []string{token}}, false)
	return err
}

func (vercelProvider) Logs(plan *types.Plan, opts LogsOptions) error {
	var state vercelState
	_ = readJSON(statePath(plan, opts.OutputDir, "vercel"), &state)
	if state.DeploymentURL == "" {
		return fmt.Errorf("no Vercel deployment URL recorded; run depctl apply first")
	}
	token := os.Getenv("VERCEL_TOKEN")
	if err := ensureVercelAuth(token, false, []string{token}); err != nil {
		return err
	}
	args := append([]string{"logs", state.DeploymentURL}, vercelGlobalArgs(token)...)
	_, err := runCommand(commandSpec{Name: "vercel", Args: args, Cwd: rootDir(plan), RedactValues: []string{token}}, false)
	return err
}

func (vercelProvider) Rollback(plan *types.Plan, opts RollbackOptions) error {
	target := opts.To
	if target == "" {
		return fmt.Errorf("provide --to with a Vercel deployment id or URL")
	}
	token := os.Getenv("VERCEL_TOKEN")
	if err := ensureVercelAuth(token, opts.DryRun, []string{token}); err != nil {
		return err
	}
	args := append([]string{"rollback", target}, vercelGlobalArgs(token)...)
	_, err := runCommand(commandSpec{Name: "vercel", Args: args, Cwd: rootDir(plan), RedactValues: []string{token}}, opts.DryRun)
	return err
}

func ensureVercelAuth(token string, dryRun bool, redact []string) error {
	if token != "" {
		return nil
	}
	if dryRun {
		output.Info("[DRY RUN] Would check Vercel login and run 'vercel login' if needed")
		return nil
	}
	if !commandExists("vercel") {
		return fmt.Errorf("vercel CLI is not installed or not in PATH")
	}

	if _, err := runCommand(commandSpec{Name: "vercel", Args: []string{"whoami"}, RedactValues: redact, AllowFailure: true}, false); err == nil {
		return nil
	}
	if !hasInteractiveTerminal() {
		return fmt.Errorf("vercel is not logged in; run 'vercel login' locally or set VERCEL_TOKEN for CI")
	}

	output.Step("Opening Vercel login...")
	_, err := runCommand(commandSpec{Name: "vercel", Args: []string{"login"}, Interactive: true}, false)
	return err
}

func vercelGlobalArgs(token string) []string {
	if token == "" {
		return nil
	}
	return []string{"--token", token}
}

func lastURL(out string) string {
	fields := strings.Fields(out)
	for i := len(fields) - 1; i >= 0; i-- {
		if strings.HasPrefix(fields[i], "https://") || strings.HasPrefix(fields[i], "http://") {
			return strings.TrimSpace(fields[i])
		}
	}
	return ""
}
