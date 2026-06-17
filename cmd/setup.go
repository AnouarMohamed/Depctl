package cmd

import (
	"github.com/spf13/cobra"
	"github.com/AnouarMohamed/Depctl/internal/output"
)

var (
	setupPreset    string
	setupDomain    string
	setupOutputDir string
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run safe preparation pipeline",
	Long:  `Shortcut that sequentializes scan, plan, write, validate, and review. Does NOT execute apply.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Info("Running depctl setup skeleton...")
		output.Step("Preset: %s", setupPreset)
		output.Step("Domain: %s", setupDomain)
		output.Step("Output Directory: %s", setupOutputDir)
	},
}

func init() {
	setupCmd.Flags().StringVar(&setupPreset, "preset", "", "deployment target preset (e.g. compose-traefik)")
	setupCmd.Flags().StringVar(&setupDomain, "domain", "", "target domain name for deployment")
	setupCmd.Flags().StringVar(&setupOutputDir, "output-dir", ".deploy", "directory where deployment files are stored")

	rootCmd.AddCommand(setupCmd)
}
