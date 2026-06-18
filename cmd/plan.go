package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/planfile"
	"github.com/AnouarMohamed/Depctl/internal/planner"
	"github.com/AnouarMohamed/Depctl/internal/target"
	"github.com/AnouarMohamed/Depctl/internal/types"
	"github.com/spf13/cobra"
)

var (
	planPreset    string
	planDomain    string
	planCI        string
	planOutputDir string
	planTarget    string
	planEnvFile   string
	planRegion    string
	planAppName   string
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Create deployment plan",
	Long:  `Transforms detected signals into a deployable plan according to specified preset, domain, and CI provider.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Step("Compiling deployment plan...")

		// Load detected.json
		detPath := filepath.Join(planOutputDir, "detected.json")
		if _, err := os.Stat(detPath); os.IsNotExist(err) {
			output.Error("detection file missing: %s does not exist. Please run 'depctl scan' first.", detPath)
			os.Exit(1)
		}

		detBytes, err := os.ReadFile(detPath)
		if err != nil {
			output.Error("failed to read %s: %v", detPath, err)
			os.Exit(1)
		}

		var det types.Detection
		if err := json.Unmarshal(detBytes, &det); err != nil {
			output.Error("failed to parse detection data: %v", err)
			os.Exit(1)
		}

		provider, err := target.Get(planTarget)
		if err != nil {
			output.Error("planning failed: %v", err)
			os.Exit(1)
		}

		// Compile plan
		plan, err := provider.Plan(&det, target.PlanOptions{
			Preset:    planPreset,
			Domain:    planDomain,
			CI:        planCI,
			Target:    planTarget,
			OutputDir: planOutputDir,
			AppName:   planAppName,
			Region:    planRegion,
			EnvFile:   planEnvFile,
		})
		if err != nil {
			output.Error("planning failed: %v", err)
			os.Exit(1)
		}

		// Ensure output directories exist
		reportsDir := filepath.Join(planOutputDir, "reports")
		if err := os.MkdirAll(reportsDir, 0755); err != nil {
			output.Error("failed to create output directory %s: %v", reportsDir, err)
			os.Exit(1)
		}

		planPath := filepath.Join(planOutputDir, "plan.json")
		if err := planfile.Save(plan, planOutputDir); err != nil {
			output.Error("failed to write %s: %v", planPath, err)
			os.Exit(1)
		}
		output.Success("Wrote plan blueprint: %s", planPath)

		// Write plan-report.md
		reportContent := planner.GeneratePlanReport(plan)
		reportPath := filepath.Join(reportsDir, "plan-report.md")
		if err := os.WriteFile(reportPath, []byte(reportContent), 0644); err != nil {
			output.Error("failed to write %s: %v", reportPath, err)
			os.Exit(1)
		}
		output.Success("Wrote plan report: %s", reportPath)

		output.Success("Plan generated successfully for target '%s'.", plan.Target.Kind)
	},
}

func init() {
	planCmd.Flags().StringVar(&planPreset, "preset", "", "deployment target preset (e.g. compose-traefik)")
	planCmd.Flags().StringVar(&planDomain, "domain", "", "target domain name for deployment")
	planCmd.Flags().StringVar(&planCI, "ci", "", "CI provider template (e.g. github, none)")
	planCmd.Flags().StringVar(&planOutputDir, "output-dir", ".deploy", "directory where deployment plan is saved")
	planCmd.Flags().StringVar(&planTarget, "target", "vps", "deployment target (vps, vercel, fly)")
	planCmd.Flags().StringVar(&planEnvFile, "env-file", ".env", "environment file to import during provider deploy")
	planCmd.Flags().StringVar(&planRegion, "region", "iad", "provider region for targets that need one")
	planCmd.Flags().StringVar(&planAppName, "app-name", "", "provider app/project name")

	rootCmd.AddCommand(planCmd)
}
