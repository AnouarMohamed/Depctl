package scanner

import (
	"os"
	"path/filepath"

	"depctl/internal/types"
)

func scanContainerization(dir string) types.ContainerizationDetection {
	var det types.ContainerizationDetection

	dockerfileNames := []string{"Dockerfile", "dockerfile", "Dockerfile.prod", "Dockerfile.production"}
	for _, name := range dockerfileNames {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			det.DockerfilePresent = true
			break
		}
	}

	composeNames := []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"}
	for _, name := range composeNames {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			det.ComposePresent = true
			break
		}
	}

	return det
}
