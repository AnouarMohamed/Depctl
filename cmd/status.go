package cmd

import (
	"os"

	"github.com/AnouarMohamed/Depctl/internal/engine"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/spf13/cobra"
)

var statusOutputDir string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check deployment status",
	Long:  `Displays provider-specific deployment status for the saved plan.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := engine.Status(statusOutputDir); err != nil {
			output.Error("Status failed: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	statusCmd.Flags().StringVar(&statusOutputDir, "output-dir", ".deploy", "directory containing the deployment configuration")
	rootCmd.AddCommand(statusCmd)
}
