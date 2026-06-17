package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

type compJSON struct {
	Require map[string]string `json:"require"`
}

func scanPHPProject(dir string) (types.RuntimeDetection, types.BuildDetection, types.NetworkDetection, bool) {
	composerPath := filepath.Join(dir, "composer.json")
	artisanPath := filepath.Join(dir, "artisan")

	hasComposer := false
	if _, err := os.Stat(composerPath); err == nil {
		hasComposer = true
	}
	hasArtisan := false
	if _, err := os.Stat(artisanPath); err == nil {
		hasArtisan = true
	}

	if !hasComposer && !hasArtisan {
		return types.RuntimeDetection{}, types.BuildDetection{}, types.NetworkDetection{}, false
	}

	runtime := types.RuntimeDetection{
		Name:       "php",
		Confidence: 0.85,
		Evidence:   []types.Evidence{},
	}

	if hasComposer {
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "composer.json", Reason: "composer.json exists"})
	}
	if hasArtisan {
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "artisan", Reason: "artisan entrypoint exists"})
	}

	isLaravel := false
	if hasComposer {
		data, err := os.ReadFile(composerPath)
		if err == nil {
			var comp compJSON
			if err := json.Unmarshal(data, &comp); err == nil {
				if comp.Require != nil {
					if _, ok := comp.Require["laravel/framework"]; ok {
						isLaravel = true
						runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "composer.json", Reason: "dependency laravel/framework detected"})
					}
				}
			}
		}
	}

	if isLaravel || hasArtisan {
		runtime.Name = "laravel"
		runtime.Framework = "laravel"
		runtime.Confidence = 0.95
	}

	build := types.BuildDetection{
		PackageManager: "composer",
		BuildCommand:   "composer install --no-dev --optimize-autoloader",
		StartCommand:   "php artisan serve --host=0.0.0.0 --port=8000", // local fallback, but production Nginx is 80
	}

	if runtime.Name == "laravel" {
		build.StartCommand = "php artisan serve --host=0.0.0.0 --port=8000"
	}

	network := types.NetworkDetection{
		InternalPort: 80, // Default Nginx port in front of php-fpm
		Confidence:   0.75,
		Evidence: []types.Evidence{
			{File: "composer.json", Reason: "PHP/Laravel deployment uses standard web port 80"},
		},
	}

	return runtime, build, network, true
}
