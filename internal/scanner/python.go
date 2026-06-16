package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"depctl/internal/types"
)

func scanPythonProject(dir string) (types.RuntimeDetection, types.BuildDetection, types.NetworkDetection, bool) {
	reqPath := filepath.Join(dir, "requirements.txt")
	pyprojPath := filepath.Join(dir, "pyproject.toml")
	pipfilePath := filepath.Join(dir, "Pipfile")
	managePath := filepath.Join(dir, "manage.py")
	mainPath := filepath.Join(dir, "main.py")
	appPyPath := filepath.Join(dir, "app.py")

	hasRequirements := false
	if _, err := os.Stat(reqPath); err == nil {
		hasRequirements = true
	}
	hasPyProj := false
	if _, err := os.Stat(pyprojPath); err == nil {
		hasPyProj = true
	}
	hasPipfile := false
	if _, err := os.Stat(pipfilePath); err == nil {
		hasPipfile = true
	}
	hasManage := false
	if _, err := os.Stat(managePath); err == nil {
		hasManage = true
	}
	hasMain := false
	if _, err := os.Stat(mainPath); err == nil {
		hasMain = true
	}
	hasAppPy := false
	if _, err := os.Stat(appPyPath); err == nil {
		hasAppPy = true
	}

	if !hasRequirements && !hasPyProj && !hasPipfile && !hasManage && !hasMain && !hasAppPy {
		return types.RuntimeDetection{}, types.BuildDetection{}, types.NetworkDetection{}, false
	}

	runtime := types.RuntimeDetection{
		Name:       "python",
		Confidence: 0.80,
		Evidence:   []types.Evidence{},
	}

	if hasRequirements {
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "requirements.txt", Reason: "requirements.txt file exists"})
	}
	if hasPyProj {
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "pyproject.toml", Reason: "pyproject.toml file exists"})
	}
	if hasPipfile {
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "Pipfile", Reason: "Pipfile exists"})
	}
	if hasManage {
		runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "manage.py", Reason: "Django management script exists"})
	}

	// Read requirements.txt or pyproject.toml to find framework dependencies
	framework := ""
	if hasRequirements {
		file, err := os.Open(reqPath)
		if err == nil {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.ToLower(scanner.Text())
				if strings.Contains(line, "fastapi") {
					framework = "fastapi"
					runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "requirements.txt", Reason: "dependency fastapi detected"})
					runtime.Confidence = 0.93
					break
				} else if strings.Contains(line, "django") {
					framework = "django"
					runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "requirements.txt", Reason: "dependency django detected"})
					runtime.Confidence = 0.95
					break
				} else if strings.Contains(line, "flask") {
					framework = "flask"
					runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "requirements.txt", Reason: "dependency flask detected"})
					runtime.Confidence = 0.90
					break
				}
			}
			file.Close()
		}
	}

	if framework == "" && hasManage {
		framework = "django"
		runtime.Confidence = 0.95
	}

	if framework == "" && hasPyProj {
		data, err := os.ReadFile(pyprojPath)
		if err == nil {
			content := strings.ToLower(string(data))
			if strings.Contains(content, "fastapi") {
				framework = "fastapi"
				runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "pyproject.toml", Reason: "dependency fastapi detected"})
				runtime.Confidence = 0.93
			} else if strings.Contains(content, "django") {
				framework = "django"
				runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "pyproject.toml", Reason: "dependency django detected"})
				runtime.Confidence = 0.95
			} else if strings.Contains(content, "flask") {
				framework = "flask"
				runtime.Evidence = append(runtime.Evidence, types.Evidence{File: "pyproject.toml", Reason: "dependency flask detected"})
				runtime.Confidence = 0.90
			}
		}
	}

	runtime.Framework = framework

	build := types.BuildDetection{
		PackageManager: "pip",
		BuildCommand:   "pip install -r requirements.txt",
		StartCommand:   "python main.py", // fallback
	}

	if framework == "fastapi" {
		build.StartCommand = "uvicorn main:app --host 0.0.0.0 --port 8000"
	} else if framework == "django" {
		build.StartCommand = "gunicorn project.wsgi:application --bind 0.0.0.0:8000"
	} else if framework == "flask" {
		build.StartCommand = "gunicorn app:app --bind 0.0.0.0:8000"
	}

	network := types.NetworkDetection{
		InternalPort: 8000,
		Confidence:   0.80,
		Evidence: []types.Evidence{
			{File: "requirements.txt", Reason: "Python web apps commonly bind to port 8000"},
		},
	}

	return runtime, build, network, true
}
