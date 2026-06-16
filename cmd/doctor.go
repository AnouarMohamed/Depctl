package cmd

import (
	"github.com/spf13/cobra"
	"depctl/internal/output"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check host environment readiness",
	Long:  `Checks OS compatibility, Docker and Compose installation, permissions, and available ports.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Info("Running depctl doctor skeleton...")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
