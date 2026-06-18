package scanner

import (
	"path/filepath"
	"testing"
)

func TestScannerFixtures(t *testing.T) {
	// Root of our repository contains fixtures/
	fixturesDir := filepath.Join("..", "..", "fixtures")

	tests := []struct {
		name               string
		fixturePath        string
		expectedRuntime    string
		expectedFramework  string
		expectedPkgManager string
		expectedPort       int
		expectedDockerfile bool
		expectedCompose    bool
	}{
		{
			name:               "Node Express",
			fixturePath:        "node-express",
			expectedRuntime:    "node",
			expectedFramework:  "express",
			expectedPkgManager: "yarn",
			expectedPort:       3000,
		},
		{
			name:               "Node Next.js",
			fixturePath:        "node-next",
			expectedRuntime:    "node",
			expectedFramework:  "nextjs",
			expectedPkgManager: "pnpm",
			expectedPort:       3000,
		},
		{
			name:               "Laravel Basic",
			fixturePath:        "laravel-basic",
			expectedRuntime:    "laravel",
			expectedFramework:  "laravel",
			expectedPkgManager: "composer",
			expectedPort:       80,
		},
		{
			name:               "Python FastAPI",
			fixturePath:        "python-fastapi",
			expectedRuntime:    "python",
			expectedFramework:  "fastapi",
			expectedPkgManager: "pip",
			expectedPort:       8000,
		},
		{
			name:               "Python Django",
			fixturePath:        "python-django",
			expectedRuntime:    "python",
			expectedFramework:  "django",
			expectedPkgManager: "pip",
			expectedPort:       8000,
		},
		{
			name:               "Go Basic",
			fixturePath:        "go-basic",
			expectedRuntime:    "go",
			expectedFramework:  "",
			expectedPkgManager: "go",
			expectedPort:       8080,
		},
		{
			name:               "Existing Dockerfile",
			fixturePath:        "existing-dockerfile",
			expectedRuntime:    "node",
			expectedDockerfile: true,
		},
		{
			name:            "Existing Compose",
			fixturePath:     "existing-compose",
			expectedRuntime: "node",
			expectedCompose: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(fixturesDir, tc.fixturePath)
			det, err := Scan(path)
			if err != nil {
				t.Fatalf("failed to scan %s: %v", path, err)
			}

			if det.Runtime.Name != tc.expectedRuntime {
				t.Errorf("runtime: got %q, want %q", det.Runtime.Name, tc.expectedRuntime)
			}

			if tc.expectedFramework != "" && det.Runtime.Framework != tc.expectedFramework {
				t.Errorf("framework: got %q, want %q", det.Runtime.Framework, tc.expectedFramework)
			}

			if tc.expectedPkgManager != "" && det.Build.PackageManager != tc.expectedPkgManager {
				t.Errorf("package manager: got %q, want %q", det.Build.PackageManager, tc.expectedPkgManager)
			}

			if tc.expectedPort != 0 && det.Network.InternalPort != tc.expectedPort {
				t.Errorf("internal port: got %d, want %d", det.Network.InternalPort, tc.expectedPort)
			}

			if det.Containerization.DockerfilePresent != tc.expectedDockerfile {
				t.Errorf("dockerfile presence: got %v, want %v", det.Containerization.DockerfilePresent, tc.expectedDockerfile)
			}

			if det.Containerization.ComposePresent != tc.expectedCompose {
				t.Errorf("compose presence: got %v, want %v", det.Containerization.ComposePresent, tc.expectedCompose)
			}
		})
	}
}
