package cmd

import (
	"github.com/spf13/cobra"
	"depctl/internal/output"
)

var (
	validateOutputDir string
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate deployment files and environment",
	Long:  `Validates generated deployment configurations for syntax errors, missing configurations, and potential security issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Info("Running depctl validate skeleton...")
		output.Step("Output Directory: %s", validateOutputDir)
	},
}

func init() {
	validateCmd.Flags().StringVar(&validateOutputDir, "output-dir", ".deploy", "directory containing deployment files to validate")

	rootCmd.AddCommand(validateCmd)
}
