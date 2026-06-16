package cmd

import (
	"github.com/spf13/cobra"
	"depctl/internal/output"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review deployment configuration",
	Long:  `Displays a summary of the deployment configuration, warnings, manual steps, and execution command.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Info("Running depctl review skeleton...")
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}
