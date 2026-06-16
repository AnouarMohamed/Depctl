package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"depctl/internal/engine"
	"depctl/internal/output"
)

var (
	rollbackTo        string
	rollbackDryRun    bool
	rollbackYes       bool
	rollbackOutputDir string
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback deployment to a previous state",
	Long:  `Restores files and services to a previously stored backup configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := engine.RollbackOptions{
			OutputDir: rollbackOutputDir,
			To:        rollbackTo,
			DryRun:    rollbackDryRun,
			Yes:       rollbackYes,
		}

		if err := engine.Rollback(opts); err != nil {
			output.Error("Rollback failed: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rollbackCmd.Flags().StringVar(&rollbackTo, "to", "", "timestamp/identifier of the target backup state to restore")
	rollbackCmd.Flags().BoolVar(&rollbackDryRun, "dry-run", false, "display rollback steps without running them")
	rollbackCmd.Flags().BoolVar(&rollbackYes, "yes", false, "skip confirmation prompts")
	rollbackCmd.Flags().StringVar(&rollbackOutputDir, "output-dir", ".deploy", "directory containing deployment files to rollback")

	rootCmd.AddCommand(rollbackCmd)
}
