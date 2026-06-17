package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/types"
)

var reviewOutputDir string

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review deployment configuration",
	Long:  `Displays a summary of the deployment configuration, warnings, manual steps, and execution command.`,
	Run: func(cmd *cobra.Command, args []string) {
		planPath := filepath.Join(reviewOutputDir, "plan.json")
		planBytes, err := os.ReadFile(planPath)
		if err != nil {
			output.Error("Plan file missing: %s. Run 'depctl plan' first.", planPath)
			os.Exit(1)
		}

		var plan types.Plan
		if err := json.Unmarshal(planBytes, &plan); err != nil {
			output.Error("Failed to parse plan file: %v", err)
			os.Exit(1)
		}

		output.Info("📋 Deployment Review for %s", plan.Project.Name)
		output.Step("Preset: %s", plan.Preset)
		output.Step("Domain: %s", plan.Domain)
		output.Step("Public Service: %s (Port %d)", plan.PublicService, plan.Network.InternalPort)
		
		output.Info("\n📦 Services:")
		for _, svc := range plan.Services {
			if svc.Type == "app" {
				output.Step("- %s (Build: %s)", svc.Name, svc.Build)
			} else {
				output.Step("- %s (Image: %s)", svc.Name, svc.Image)
			}
		}

		if len(plan.Warnings) > 0 {
			output.Warning("\n⚠️  Warnings:")
			for _, warn := range plan.Warnings {
				output.Step("- %s", warn)
			}
		}

		output.Info("\n📝 Manual Steps Required:")
		for _, step := range plan.ManualSteps {
			output.Step("- %s", step)
		}

		output.Success("\nReady to deploy! Run 'depctl apply' to start.")
	},
}

func init() {
	reviewCmd.Flags().StringVar(&reviewOutputDir, "output-dir", ".deploy", "directory containing the deployment plan to review")
	rootCmd.AddCommand(reviewCmd)
}
