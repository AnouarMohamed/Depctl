package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/AnouarMohamed/Depctl/internal/types"
	"github.com/AnouarMohamed/Depctl/internal/writer"
)

var (
	writeOutputDir string
	writeForce     bool
)

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write deployment kit files",
	Long:  `Generates Dockerfiles, Docker Compose files, configuration files, and utility scripts in the deployment folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Step("Loading plan from %s/plan.json...", writeOutputDir)

		planPath := filepath.Join(writeOutputDir, "plan.json")
		planBytes, err := os.ReadFile(planPath)
		if err != nil {
			output.Error("Failed to read plan file: %v. Run 'depctl plan' first.", err)
			os.Exit(1)
		}

		var plan types.Plan
		if err := json.Unmarshal(planBytes, &plan); err != nil {
			output.Error("Failed to parse plan file: %v", err)
			os.Exit(1)
		}

		output.Step("Generating deployment kit in %s...", writeOutputDir)
		if err := writer.Write(&plan, writeOutputDir, writeForce); err != nil {
			output.Error("Failed to write deployment kit: %v", err)
			os.Exit(1)
		}

		output.Success("Deployment kit successfully written to %s", writeOutputDir)
		output.Info("Next steps:")
		for _, step := range plan.ManualSteps {
			output.Step("- %s", step)
		}
	},
}

func init() {
	writeCmd.Flags().StringVar(&writeOutputDir, "output-dir", ".deploy", "directory where deployment kit is written")
	writeCmd.Flags().BoolVar(&writeForce, "force", false, "force overwrite existing files in the deployment directory")

	rootCmd.AddCommand(writeCmd)
}
