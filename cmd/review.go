package cmd

import (
	"os"

	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/planfile"
	"github.com/spf13/cobra"
)

var reviewOutputDir string

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review deployment configuration",
	Long:  `Displays a summary of the deployment configuration, warnings, manual steps, and execution command.`,
	Run: func(cmd *cobra.Command, args []string) {
		plan, err := planfile.Load(reviewOutputDir)
		if err != nil {
			output.Error("Plan file missing or invalid. Run 'depctl plan' first: %v", err)
			os.Exit(1)
		}

		output.Info("📋 Deployment Review for %s", plan.Project.Name)
		output.Step("Target: %s", plan.Target.Kind)
		output.Step("Preset: %s", plan.Preset)
		if plan.Domain != "" {
			output.Step("Domain: %s", plan.Domain)
		}
		if plan.Target.AppName != "" {
			output.Step("Provider App: %s", plan.Target.AppName)
		}
		output.Step("Public Service: %s (Port %d)", plan.PublicService, plan.Network.InternalPort)

		output.Info("\n📦 Services:")
		for _, svc := range plan.Services {
			if svc.Type == "app" {
				output.Step("- %s (Build: %s)", svc.Name, svc.Build)
			} else {
				output.Step("- %s (Image: %s)", svc.Name, svc.Image)
			}
		}

		if len(plan.Artifacts) > 0 {
			output.Info("\n🧾 Root Artifacts:")
			for _, artifact := range plan.Artifacts {
				if artifact.Scope == "root" {
					output.Step("- %s", artifact.Path)
				}
			}
		}

		if len(plan.Credentials) > 0 {
			output.Info("\n🔐 Credentials:")
			for _, cred := range plan.Credentials {
				output.Step("- %s (%s or %s)", cred.Name, cred.EnvVar, cred.Command)
			}
		}

		if len(plan.SecretImports) > 0 {
			output.Info("\n🔒 Secret Imports:")
			for _, imp := range plan.SecretImports {
				output.Step("- %s keys from %s", len(imp.Keys), imp.SourceFile)
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

		output.Success("\nReady to deploy! Run 'depctl apply' to start, or 'depctl apply --dry-run' to preview.")
	},
}

func init() {
	reviewCmd.Flags().StringVar(&reviewOutputDir, "output-dir", ".deploy", "directory containing the deployment plan to review")
	rootCmd.AddCommand(reviewCmd)
}
