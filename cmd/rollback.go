package cmd

import (
	"github.com/spf13/cobra"
	"depctl/internal/output"
)

var (
	rollbackTo     string
	rollbackDryRun bool
	rollbackYes    bool
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback deployment to a previous state",
	Long:  `Restores files and services to a previously stored backup configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Info("Running depctl rollback skeleton...")
		output.Step("Rollback target backup: %s", rollbackTo)
		output.Step("Dry Run: %v", rollbackDryRun)
		output.Step("Skip Confirmation (--yes): %v", rollbackYes)
	},
}

func init() {
	rollbackCmd.Flags().StringVar(&rollbackTo, "to", "", "timestamp/identifier of the target backup state to restore")
	rollbackCmd.Flags().BoolVar(&rollbackDryRun, "dry-run", false, "display rollback steps without running them")
	rollbackCmd.Flags().BoolVar(&rollbackYes, "yes", false, "skip confirmation prompts")

	rootCmd.AddCommand(rollbackCmd)
}
