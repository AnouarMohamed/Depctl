package cmd

import (
	"os"

	"github.com/AnouarMohamed/Depctl/internal/engine"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/spf13/cobra"
)

var (
	applyYes             bool
	applyDryRun          bool
	applySkipBuild       bool
	applySkipHealthcheck bool
	applyOutputDir       string
	applyEnvFile         string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply deployment plan",
	Long:  `Executes the deployment actions outlined in the saved plan.json, backing up the current state first.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := engine.ApplyOptions{
			OutputDir:       applyOutputDir,
			Yes:             applyYes,
			DryRun:          applyDryRun,
			SkipBuild:       applySkipBuild,
			SkipHealthcheck: applySkipHealthcheck,
			EnvFile:         applyEnvFile,
		}

		if err := engine.Apply(opts); err != nil {
			output.Error("Apply failed: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	applyCmd.Flags().BoolVar(&applyYes, "yes", false, "skip confirmation prompts")
	applyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false, "display execution steps without running them")
	applyCmd.Flags().BoolVar(&applySkipBuild, "skip-build", false, "skip building new Docker images")
	applyCmd.Flags().BoolVar(&applySkipHealthcheck, "skip-healthcheck", false, "skip post-deployment health check checks")
	applyCmd.Flags().StringVar(&applyOutputDir, "output-dir", ".deploy", "directory containing deployment files to apply")
	applyCmd.Flags().StringVar(&applyEnvFile, "env-file", "", "override environment file used for provider secret import")

	rootCmd.AddCommand(applyCmd)
}
