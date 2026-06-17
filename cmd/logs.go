package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/AnouarMohamed/Depctl/internal/output"
)

var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "Inspect deployment container logs",
	Long:  `Tails stdout and stderr from running containers. Optionally filter by service type (app, proxy, db).`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Info("Running depctl logs skeleton...")
		if len(args) > 0 {
			output.Step("Filtering logs for service: %s", strings.Join(args, ", "))
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
