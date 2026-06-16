package cmd

import (
	"github.com/spf13/cobra"
	"depctl/internal/output"
)

var (
	applyYes             bool
	applyDryRun          bool
	applySkipBuild       bool
	applySkipHealthcheck bool
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply deployment plan",
	Long:  `Executes the deployment actions outlined in the saved plan.json, backing up the current state first.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Info("Running depctl apply skeleton...")
		output.Step("Skip Confirmation (--yes): %v", applyYes)
		output.Step("Dry Run: %v", applyDryRun)
		output.Step("Skip Build: %v", applySkipBuild)
		output.Step("Skip Healthcheck: %v", applySkipHealthcheck)
	},
}

func init() {
	applyCmd.Flags().BoolVar(&applyYes, "yes", false, "skip confirmation prompts")
	applyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false, "display execution steps without running them")
	applyCmd.Flags().BoolVar(&applySkipBuild, "skip-build", false, "skip building new Docker images")
	applyCmd.Flags().BoolVar(&applySkipHealthcheck, "skip-healthcheck", false, "skip post-deployment health check checks")

	rootCmd.AddCommand(applyCmd)
}
