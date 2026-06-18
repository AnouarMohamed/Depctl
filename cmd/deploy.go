package cmd

import (
	"os"

	"github.com/AnouarMohamed/Depctl/internal/engine"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/spf13/cobra"
)

var (
	deployPreset          string
	deployDomain          string
	deployOutputDir       string
	deployTarget          string
	deployEnvFile         string
	deployRegion          string
	deployAppName         string
	deployCI              string
	deployForce           bool
	deployYes             bool
	deployDryRun          bool
	deploySkipBuild       bool
	deploySkipHealthcheck bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy [path]",
	Short: "Prepare and deploy an app",
	Long:  `Runs scan, plan, write, validate, review, and apply with provider-specific automation.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		err := runPreparePipeline(path, pipelineOptions{
			Preset:    deployPreset,
			Domain:    deployDomain,
			OutputDir: deployOutputDir,
			Target:    deployTarget,
			EnvFile:   deployEnvFile,
			Region:    deployRegion,
			AppName:   deployAppName,
			CI:        deployCI,
			Force:     deployForce,
		})
		if err != nil {
			output.Error("Deploy preparation failed: %v", err)
			os.Exit(1)
		}
		if err := engine.Apply(engine.ApplyOptions{
			OutputDir:       deployOutputDir,
			Yes:             deployYes,
			DryRun:          deployDryRun,
			SkipBuild:       deploySkipBuild,
			SkipHealthcheck: deploySkipHealthcheck,
			EnvFile:         deployEnvFile,
		}); err != nil {
			output.Error("Deploy failed: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	deployCmd.Flags().StringVar(&deployPreset, "preset", "compose-traefik", "deployment preset for VPS target")
	deployCmd.Flags().StringVar(&deployDomain, "domain", "", "target domain name")
	deployCmd.Flags().StringVar(&deployOutputDir, "output-dir", ".deploy", "directory where deployment files are stored")
	deployCmd.Flags().StringVar(&deployTarget, "target", "vps", "deployment target (vps, vercel, fly)")
	deployCmd.Flags().StringVar(&deployEnvFile, "env-file", ".env", "environment file to import during provider deploy")
	deployCmd.Flags().StringVar(&deployRegion, "region", "iad", "provider region for targets that need one")
	deployCmd.Flags().StringVar(&deployAppName, "app-name", "", "provider app/project name")
	deployCmd.Flags().StringVar(&deployCI, "ci", "github", "CI provider template (github, none)")
	deployCmd.Flags().BoolVar(&deployForce, "force", false, "backup and overwrite generated files")
	deployCmd.Flags().BoolVar(&deployYes, "yes", false, "skip confirmation prompts")
	deployCmd.Flags().BoolVar(&deployDryRun, "dry-run", false, "display apply steps without running provider commands")
	deployCmd.Flags().BoolVar(&deploySkipBuild, "skip-build", false, "skip Docker build where supported")
	deployCmd.Flags().BoolVar(&deploySkipHealthcheck, "skip-healthcheck", false, "skip post-deployment health checks where supported")

	rootCmd.AddCommand(deployCmd)
}
