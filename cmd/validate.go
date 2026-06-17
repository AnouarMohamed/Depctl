package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/validator"
)

var (
	validateOutputDir string
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate deployment files and environment",
	Long:  `Validates generated deployment configurations for syntax errors, missing configurations, and potential security issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Step("Validating deployment kit in %s...", validateOutputDir)

		res, err := validator.Validate(validateOutputDir)
		if err != nil {
			output.Error("Validation failed: %v", err)
			os.Exit(1)
		}

		// Ensure reports directory exists
		reportsDir := filepath.Join(validateOutputDir, "reports")
		if err := os.MkdirAll(reportsDir, 0755); err != nil {
			output.Error("Failed to create reports directory: %v", err)
			os.Exit(1)
		}

		// Write report
		reportContent := validator.GenerateValidationReport(res)
		reportPath := filepath.Join(reportsDir, "validation-report.md")
		if err := os.WriteFile(reportPath, []byte(reportContent), 0644); err != nil {
			output.Error("Failed to write validation report: %v", err)
			os.Exit(1)
		}
		output.Success("Wrote validation report: %s", reportPath)

		if !res.Valid {
			output.Error("Deployment kit is INVALID.")
			for _, err := range res.Errors {
				output.Step("- [ERROR] %s", err)
			}
			os.Exit(1)
		}

		if len(res.Warnings) > 0 {
			output.Warning("Deployment kit has warnings.")
			for _, warn := range res.Warnings {
				output.Step("- [WARN] %s", warn)
			}
		}

		output.Success("Deployment kit is VALID and safe to apply.")
	},
}

func init() {
	validateCmd.Flags().StringVar(&validateOutputDir, "output-dir", ".deploy", "directory containing deployment files to validate")

	rootCmd.AddCommand(validateCmd)
}
