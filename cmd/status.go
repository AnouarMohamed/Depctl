package cmd

import (
	"github.com/spf13/cobra"
	"depctl/internal/output"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check deployment status",
	Long:  `Displays running containers, internal port mapping, reverse proxy routes, and container health.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Info("Running depctl status skeleton...")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
