package scanner

import (
	"os"
	"path/filepath"

	"depctl/internal/types"
)

// Scan analyzes the given directory and returns a structural representation of all detected signals.
func Scan(dir string) (*types.Detection, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	projectName := filepath.Base(absPath)

	// Check if gitRoot exists
	gitRoot := false
	if _, err := os.Stat(filepath.Join(absPath, ".git")); err == nil {
		gitRoot = true
	}

	projectDet := types.ProjectDetection{
		Name:    projectName,
		Root:    absPath,
		GitRoot: gitRoot,
	}

	var runtime types.RuntimeDetection
	var build types.BuildDetection
	var network types.NetworkDetection
	matched := false

	// Try Node first
	if r, b, n, ok := scanNodeProject(absPath); ok {
		runtime, build, network = r, b, n
		matched = true
	}

	// Try PHP/Laravel next
	if !matched {
		if r, b, n, ok := scanPHPProject(absPath); ok {
			runtime, build, network = r, b, n
			matched = true
		}
	}

	// Try Python next
	if !matched {
		if r, b, n, ok := scanPythonProject(absPath); ok {
			runtime, build, network = r, b, n
			matched = true
		}
	}

	// Try Go next
	if !matched {
		if r, b, n, ok := scanGoProject(absPath); ok {
			runtime, build, network = r, b, n
			matched = true
		}
	}

	// Fallback if no runtime matches
	if !matched {
		runtime = types.RuntimeDetection{
			Name:       "unknown",
			Framework:  "",
			Confidence: 0.0,
			Evidence:   []types.Evidence{},
		}
		build = types.BuildDetection{}
		network = types.NetworkDetection{
			InternalPort: 80,
			Confidence:   0.0,
		}
	}

	// Scan environment variables
	env, envWarnings := scanEnvFiles(absPath)

	// Detect DB dependencies
	dependencies := detectDependencies(absPath, env.Keys, runtime)

	// Scan existing container configuration
	container := scanContainerization(absPath)

	// Scan CI pipelines
	ci := scanCI(absPath)

	// Compile Warnings
	var warnings []string
	warnings = append(warnings, envWarnings...)

	if !container.DockerfilePresent {
		warnings = append(warnings, "No Dockerfile found.")
	}

	// If a database connection is required (has DATABASE_URL) but no database dependency is marked
	hasDbUrl := false
	for _, key := range env.Keys {
		if key == "DATABASE_URL" || key == "DB_CONNECTION" {
			hasDbUrl = true
			break
		}
	}

	isDbLikely := false
	for _, dep := range dependencies {
		if dep.Likely {
			isDbLikely = true
			break
		}
	}

	if hasDbUrl && !isDbLikely {
		warnings = append(warnings, "DATABASE_URL referenced but no database service exists.")
	}

	return &types.Detection{
		Version:          "0.1",
		Project:          projectDet,
		Runtime:          runtime,
		Build:            build,
		Network:          network,
		Env:              env,
		Dependencies:     dependencies,
		Containerization: container,
		CI:               ci,
		Warnings:         warnings,
	}, nil
}
