package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"depctl/internal/output"
)

var (
	quiet bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "depctl",
	Short: "depctl is a repo-aware VPS deployment-kit compiler",
	Long: `depctl is a command-line DevOps tool that scans your code repository,
generates clean deployment files (.deploy/), validates them, and applies them
safely to your target VPS.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if quiet {
			output.SetQuiet(true)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress standard output logging")
}
