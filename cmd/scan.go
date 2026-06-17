package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/scanner"
)

var (
	scanOutputDir string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project directory for deployment signals",
	Long:  `Analyzes the codebase for runtimes, frameworks, environment variables, databases, and container configurations.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output.Step("Starting project scan...")

		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}

		absTargetDir, err := filepath.Abs(targetDir)
		if err != nil {
			output.Error("failed to get absolute path for %s: %v", targetDir, err)
			os.Exit(1)
		}

		// Execute Scan
		det, err := scanner.Scan(absTargetDir)
		if err != nil {
			output.Error("scan failed: %v", err)
			os.Exit(1)
		}

		// Ensure output directories exist
		reportsDir := filepath.Join(scanOutputDir, "reports")
		if err := os.MkdirAll(reportsDir, 0755); err != nil {
			output.Error("failed to create output directory %s: %v", reportsDir, err)
			os.Exit(1)
		}

		// Write detected.json
		detBytes, err := json.MarshalIndent(det, "", "  ")
		if err != nil {
			output.Error("failed to serialize scan data: %v", err)
			os.Exit(1)
		}

		detPath := filepath.Join(scanOutputDir, "detected.json")
		if err := os.WriteFile(detPath, detBytes, 0644); err != nil {
			output.Error("failed to write %s: %v", detPath, err)
			os.Exit(1)
		}
		output.Success("Wrote detection artifact: %s", detPath)

		// Write scan-report.md
		reportContent := scanner.GenerateScanReport(det)
		reportPath := filepath.Join(reportsDir, "scan-report.md")
		if err := os.WriteFile(reportPath, []byte(reportContent), 0644); err != nil {
			output.Error("failed to write %s: %v", reportPath, err)
			os.Exit(1)
		}
		output.Success("Wrote scan report: %s", reportPath)

		band := scanner.ConfidenceBand(det.Runtime.Confidence)
		output.Success("Scan complete: %s (%s) detected with %s confidence.",
			det.Runtime.Name, det.Runtime.Framework, band)
	},
}

func init() {
	scanCmd.Flags().StringVar(&scanOutputDir, "output-dir", ".deploy", "directory where scan results are saved")
	rootCmd.AddCommand(scanCmd)
}

