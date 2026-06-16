package cmd

import (
	"github.com/spf13/cobra"
	"depctl/internal/output"
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
		output.Info("Running depctl write skeleton...")
		output.Step("Output Directory: %s", writeOutputDir)
		output.Step("Force Overwrite: %v", writeForce)
	},
}

func init() {
	writeCmd.Flags().StringVar(&writeOutputDir, "output-dir", ".deploy", "directory where deployment kit is written")
	writeCmd.Flags().BoolVar(&writeForce, "force", false, "force overwrite existing files in the deployment directory")

	rootCmd.AddCommand(writeCmd)
}
