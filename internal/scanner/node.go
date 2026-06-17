package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

type pkgJSON struct {
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func scanNodeProject(dir string) (types.RuntimeDetection, types.BuildDetection, types.NetworkDetection, bool) {
	pkgPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return types.RuntimeDetection{}, types.BuildDetection{}, types.NetworkDetection{}, false
	}

	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return types.RuntimeDetection{}, types.BuildDetection{}, types.NetworkDetection{}, false
	}

	var pkg pkgJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return types.RuntimeDetection{}, types.BuildDetection{}, types.NetworkDetection{}, false
	}

	runtime := types.RuntimeDetection{
		Name:       "node",
		Confidence: 0.90, // starts at 90% since package.json exists
		Evidence: []types.Evidence{
			{File: "package.json", Reason: "package.json exists"},
		},
	}

	// Detect framework
	framework := ""
	if pkg.Dependencies != nil {
		if _, ok := pkg.Dependencies["next"]; ok {
			framework = "nextjs"
			runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "package.json", Reason: "dependency next detected"})
			runtime.Confidence = 0.95
		} else if _, ok := pkg.Dependencies["express"]; ok {
			framework = "express"
			runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "package.json", Reason: "dependency express detected"})
			runtime.Confidence = 0.93
		} else if _, ok := pkg.Dependencies["fastify"]; ok {
			framework = "fastify"
			runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "package.json", Reason: "dependency fastify detected"})
			runtime.Confidence = 0.93
		} else if _, ok := pkg.Dependencies["@nestjs/core"]; ok {
			framework = "nestjs"
			runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "package.json", Reason: "dependency @nestjs/core detected"})
			runtime.Confidence = 0.93
		} else if _, ok := pkg.Dependencies["nuxt"]; ok {
			framework = "nuxt"
			runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "package.json", Reason: "dependency nuxt detected"})
			runtime.Confidence = 0.93
		} else if _, ok := pkg.Dependencies["vite"]; ok {
			framework = "vite"
			runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "package.json", Reason: "dependency vite detected"})
			runtime.Confidence = 0.90
		}
	}
	runtime.Framework = framework

	// Detect package manager
	pkgManager := "npm"
	if _, err := os.Stat(filepath.Join(dir, "pnpm-lock.yaml")); err == nil {
		pkgManager = "pnpm"
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "pnpm-lock.yaml", Reason: "pnpm lockfile detected"})
	} else if _, err := os.Stat(filepath.Join(dir, "yarn.lock")); err == nil {
		pkgManager = "yarn"
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "yarn.lock", Reason: "yarn lockfile detected"})
	} else if _, err := os.Stat(filepath.Join(dir, "bun.lockb")); err == nil {
		pkgManager = "bun"
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "bun.lockb", Reason: "bun lockfile detected"})
	} else if _, err := os.Stat(filepath.Join(dir, "package-lock.json")); err == nil {
		pkgManager = "npm"
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "package-lock.json", Reason: "npm lockfile detected"})
	}

	build := types.BuildDetection{
		PackageManager: pkgManager,
	}

	// Commands
	if pkg.Scripts != nil {
		if val, ok := pkg.Scripts["build"]; ok {
			build.BuildCommand = pkgManager + " run build"
			_ = val
		}
		if _, ok := pkg.Scripts["start"]; ok {
			build.StartCommand = pkgManager + " start"
		}
	}

	// Fallback/Defaults for commands if missing
	if build.StartCommand == "" {
		if framework == "nextjs" {
			build.StartCommand = "next start"
		} else if framework == "nuxt" {
			build.StartCommand = "nuxt start"
		} else if framework == "vite" {
			build.StartCommand = "vite preview"
		} else {
			build.StartCommand = "node index.js"
		}
	}
	if build.BuildCommand == "" && (framework == "nextjs" || framework == "nuxt" || framework == "vite") {
		build.BuildCommand = pkgManager + " run build"
	}

	// Port detection
	port := 3000
	if framework == "vite" {
		port = 4173
	}

	network := types.NetworkDetection{
		InternalPort: port,
		Confidence:   0.80,
		Evidence: []types.Evidence{
			{File: "package.json", Reason: "default framework port inferred"},
		},
	}

	return runtime, build, network, true
}
