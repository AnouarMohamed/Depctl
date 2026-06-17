package scanner

import (
	"os"
	"path/filepath"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

func scanGoProject(dir string) (types.RuntimeDetection, types.BuildDetection, types.NetworkDetection, bool) {
	goModPath := filepath.Join(dir, "go.mod")
	mainGoPath := filepath.Join(dir, "main.go")

	hasGoMod := false
	if _, err := os.Stat(goModPath); err == nil {
		hasGoMod = true
	}
	hasMainGo := false
	if _, err := os.Stat(mainGoPath); err == nil {
		hasMainGo = true
	}

	if !hasGoMod && !hasMainGo {
		return types.RuntimeDetection{}, types.BuildDetection{}, types.NetworkDetection{}, false
	}

	runtime := types.RuntimeDetection{
		Name:       "go",
		Confidence: 0.90,
		Evidence:   []types.Evidence{},
	}

	if hasGoMod {
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "go.mod", Reason: "go.mod file exists"})
	}
	if hasMainGo {
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "main.go", Reason: "main.go entrypoint exists"})
	}

	build := types.BuildDetection{
		PackageManager: "go",
		BuildCommand:   "go build -o app .",
		StartCommand:   "./app",
	}

	network := types.NetworkDetection{
		InternalPort: 8080, // standard Go web service port default
		Confidence:   0.60,
		Evidence: []types.Evidence{
			{File: "go.mod", Reason: "default port 8080 inferred for Go application"},
		},
	}

	return runtime, build, network, true
}
