package cmd

import (
	"os"
	"path/filepath"

	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/planfile"
	"github.com/AnouarMohamed/Depctl/internal/planner"
	"github.com/AnouarMohamed/Depctl/internal/scanner"
	"github.com/AnouarMohamed/Depctl/internal/target"
	"github.com/AnouarMohamed/Depctl/internal/validator"
	"github.com/spf13/cobra"
)

var (
	setupPreset    string
	setupDomain    string
	setupOutputDir string
	setupTarget    string
	setupEnvFile   string
	setupRegion    string
	setupAppName   string
	setupCI        string
	setupForce     bool
)

type pipelineOptions struct {
	Preset    string
	Domain    string
	OutputDir string
	Target    string
	EnvFile   string
	Region    string
	AppName   string
	CI        string
	Force     bool
}

var setupCmd = &cobra.Command{
	Use:   "setup [path]",
	Short: "Prepare deployment files",
	Long:  `Runs scan, plan, write, validate, and review. It never deploys.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		if err := runPreparePipeline(path, pipelineOptions{
			Preset:    setupPreset,
			Domain:    setupDomain,
			OutputDir: setupOutputDir,
			Target:    setupTarget,
			EnvFile:   setupEnvFile,
			Region:    setupRegion,
			AppName:   setupAppName,
			CI:        setupCI,
			Force:     setupForce,
		}); err != nil {
			output.Error("Setup failed: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	setupCmd.Flags().StringVar(&setupPreset, "preset", "compose-traefik", "deployment preset for VPS target")
	setupCmd.Flags().StringVar(&setupDomain, "domain", "", "target domain name")
	setupCmd.Flags().StringVar(&setupOutputDir, "output-dir", ".deploy", "directory where deployment files are stored")
	setupCmd.Flags().StringVar(&setupTarget, "target", "vps", "deployment target (vps, vercel, fly)")
	setupCmd.Flags().StringVar(&setupEnvFile, "env-file", ".env", "environment file to import during provider deploy")
	setupCmd.Flags().StringVar(&setupRegion, "region", "iad", "provider region for targets that need one")
	setupCmd.Flags().StringVar(&setupAppName, "app-name", "", "provider app/project name")
	setupCmd.Flags().StringVar(&setupCI, "ci", "github", "CI provider template (github, none)")
	setupCmd.Flags().BoolVar(&setupForce, "force", false, "backup and overwrite generated files")

	rootCmd.AddCommand(setupCmd)
}

func runPreparePipeline(path string, opts pipelineOptions) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	output.Step("Scanning project...")
	det, err := scanner.Scan(absPath)
	if err != nil {
		return err
	}
	reportsDir := filepath.Join(opts.OutputDir, "reports")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return err
	}
	if err := writeJSONFile(filepath.Join(opts.OutputDir, "detected.json"), det); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reportsDir, "scan-report.md"), []byte(scanner.GenerateScanReport(det)), 0644); err != nil {
		return err
	}

	provider, err := target.Get(opts.Target)
	if err != nil {
		return err
	}
	output.Step("Planning %s deployment...", provider.Kind())
	plan, err := provider.Plan(det, target.PlanOptions{
		Preset:    opts.Preset,
		Domain:    opts.Domain,
		CI:        opts.CI,
		Target:    opts.Target,
		OutputDir: opts.OutputDir,
		AppName:   opts.AppName,
		Region:    opts.Region,
		EnvFile:   opts.EnvFile,
	})
	if err != nil {
		return err
	}
	if err := planfile.Save(plan, opts.OutputDir); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reportsDir, "plan-report.md"), []byte(planner.GeneratePlanReport(plan)), 0644); err != nil {
		return err
	}

	output.Step("Writing deployment files...")
	if err := provider.Write(plan, target.WriteOptions{OutputDir: opts.OutputDir, Force: opts.Force}); err != nil {
		return err
	}

	output.Step("Validating deployment files...")
	res, err := provider.Validate(plan, target.ValidateOptions{OutputDir: opts.OutputDir})
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reportsDir, "validation-report.md"), []byte(validator.GenerateValidationReport(res)), 0644); err != nil {
		return err
	}
	if !res.Valid {
		for _, validationErr := range res.Errors {
			output.Error(validationErr)
		}
		return os.ErrInvalid
	}
	for _, warning := range res.Warnings {
		output.Warning(warning)
	}

	output.Success("Setup complete. Review with 'depctl review', then deploy with 'depctl apply'.")
	return nil
}

func writeJSONFile(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := jsonMarshal(value)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
