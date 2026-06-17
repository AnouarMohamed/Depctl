package doctor

import (
	"net"
	"os/exec"
	"strings"

	"github.com/AnouarMohamed/Depctl/internal/output"
)

// Check performs a series of system checks and reports findings.
func Check() {
	output.Step("Checking Docker installation...")
	if err := checkCommand("docker", "--version"); err != nil {
		output.Error("- Docker is not installed or not in PATH.")
	} else {
		output.Success("- Docker is installed.")
	}

	output.Step("Checking Docker Compose installation...")
	if err := checkCommand("docker", "compose", "version"); err != nil {
		output.Error("- Docker Compose is not installed.")
	} else {
		output.Success("- Docker Compose is installed.")
	}

	output.Step("Checking Docker permissions...")
	if err := checkCommand("docker", "ps"); err != nil {
		output.Error("- Current user cannot run Docker commands. Try adding user to 'docker' group.")
	} else {
		output.Success("- Docker permissions are OK.")
	}

	output.Step("Checking port 80 (HTTP)...")
	if isPortOpen("80") {
		output.Warning("- Port 80 is already in use. This may conflict with Traefik.")
	} else {
		output.Success("- Port 80 is available.")
	}

	output.Step("Checking port 443 (HTTPS)...")
	if isPortOpen("443") {
		output.Warning("- Port 443 is already in use. This may conflict with Traefik.")
	} else {
		output.Success("- Port 443 is available.")
	}

	output.Step("Checking for existing Traefik containers...")
	if hasTraefik() {
		output.Warning("- An existing Traefik container was detected. Ensure it doesn't conflict with depctl.")
	} else {
		output.Success("- No conflicting Traefik containers found.")
	}
}

func checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func isPortOpen(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	_ = ln.Close()
	return false
}

func hasTraefik() bool {
	cmd := exec.Command("docker", "ps", "--filter", "name=traefik", "--format", "{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}
