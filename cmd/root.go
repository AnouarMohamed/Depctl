package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/spf13/cobra"
)

var (
	quiet         bool
	noBanner      bool
	bannerText    string
	bannerPrinted bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "depctl",
	Short: "depctl prepares and deploys production app configs",
	Long: `depctl scans a single-app repository, generates reviewable deployment
configuration, validates it, and applies it through supported targets such as
VPS, Fly.io, and Vercel.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if quiet {
			output.SetQuiet(true)
		}
		maybePrintBanner()
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// SetBanner configures the startup banner embedded by the main package.
func SetBanner(text string) {
	bannerText = strings.TrimRight(text, "\n")
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
	rootCmd.PersistentFlags().BoolVar(&noBanner, "no-banner", false, "suppress the interactive startup banner")
}

func maybePrintBanner() {
	if bannerPrinted || quiet || noBanner || bannerText == "" {
		return
	}
	if os.Getenv("DEPCTL_NO_BANNER") != "" || os.Getenv("CI") != "" {
		return
	}
	if !isInteractiveStdout() {
		return
	}
	fmt.Fprintln(os.Stdout, bannerText)
	fmt.Fprintln(os.Stdout)
	bannerPrinted = true
}

func isInteractiveStdout() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
