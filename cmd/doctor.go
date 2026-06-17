package cmd

import (
	"github.com/spf13/cobra"
	"github.com/AnouarMohamed/Depctl/internal/doctor"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check host environment readiness",
	Long:  `Checks OS compatibility, Docker and Compose installation, permissions, and available ports.`,
	Run: func(cmd *cobra.Command, args []string) {
		doctor.Check()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
