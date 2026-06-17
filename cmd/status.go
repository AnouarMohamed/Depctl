package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/AnouarMohamed/Depctl/internal/output"
)

var statusOutputDir string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check deployment status",
	Long:  `Displays running containers, internal port mapping, reverse proxy routes, and container health.`,
	Run: func(cmd *cobra.Command, args []string) {
		composePath := filepath.Join(statusOutputDir, "docker-compose.yml")
		if _, err := os.Stat(composePath); os.IsNotExist(err) {
			output.Error("Compose file missing: %s. Run 'depctl write' and 'depctl apply' first.", composePath)
			os.Exit(1)
		}

		output.Info("📊 Current Deployment Status:")
		
		psCmd := exec.Command("docker", "compose", "-f", composePath, "ps")
		psCmd.Stdout = os.Stdout
		psCmd.Stderr = os.Stderr
		if err := psCmd.Run(); err != nil {
			output.Error("Failed to get container status: %v", err)
		}

		output.Info("\n📝 Recent Logs (last 10 lines):")
		logsCmd := exec.Command("docker", "compose", "-f", composePath, "logs", "--tail=10")
		logsCmd.Stdout = os.Stdout
		logsCmd.Stderr = os.Stderr
		if err := logsCmd.Run(); err != nil {
			output.Error("Failed to get container logs: %v", err)
		}
	},
}

func init() {
	statusCmd.Flags().StringVar(&statusOutputDir, "output-dir", ".deploy", "directory containing the deployment configuration")
	rootCmd.AddCommand(statusCmd)
}
