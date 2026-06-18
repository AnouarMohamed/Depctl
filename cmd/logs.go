package cmd

import (
	"os"

	"github.com/AnouarMohamed/Depctl/internal/engine"
	"github.com/AnouarMohamed/Depctl/internal/output"
	"github.com/spf13/cobra"
)

var (
	logsOutputDir string
	logsTail      int
)

var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "Inspect deployment logs",
	Long:  `Shows provider-specific logs. For VPS, an optional service name filters docker compose logs.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := ""
		if len(args) > 0 {
			service = args[0]
		}
		if err := engine.Logs(logsOutputDir, service, logsTail); err != nil {
			output.Error("Logs failed: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	logsCmd.Flags().StringVar(&logsOutputDir, "output-dir", ".deploy", "directory containing the deployment configuration")
	logsCmd.Flags().IntVar(&logsTail, "tail", 100, "number of log lines to show where supported")
	rootCmd.AddCommand(logsCmd)
}
