package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

var sensitivePatterns = []string{
	"SECRET", "TOKEN", "PASSWORD", "PASS", "PRIVATE",
	"KEY", "API_KEY", "DB_PASS", "JWT", "SESSION",
	"CREDENTIAL", "ACCESS_KEY",
}

// isSensitiveKey checks if the key name contains any sensitive substring.
func isSensitiveKey(key string) bool {
	upperKey := strings.ToUpper(key)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(upperKey, pattern) {
			return true
		}
	}
	return false
}

// scanEnvFiles checks the project root for .env or .env.example files.
// It extracts keys only (not values) and flags sensitive keys.
func scanEnvFiles(dir string) (types.EnvDetection, []string) {
	var det types.EnvDetection
	var warnings []string

	envExamplePath := filepath.Join(dir, ".env.example")
	envPath := filepath.Join(dir, ".env")

	targetPath := ""
	if _, err := os.Stat(envExamplePath); err == nil {
		det.HasEnvExample = true
		targetPath = envExamplePath
	} else if _, err := os.Stat(envPath); err == nil {
		targetPath = envPath
		warnings = append(warnings, "No .env.example found (but .env exists).")
	}

	if targetPath == "" {
		warnings = append(warnings, "No .env.example found.")
		return det, warnings
	}

	file, err := os.Open(targetPath)
	if err != nil {
		return det, warnings
	}
	defer file.Close()

	keyMap := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	// Simple regex to match KEY=VALUE lines (ignoring comments)
	envRegex := regexp.MustCompile(`^\s*([A-Za-z0-9_]+)\s*=`)

	for scanner.Scan() {
		line := scanner.Text()
		// Skip comment lines
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		matches := envRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			key := matches[1]
			keyMap[key] = true
		}
	}

	for k := range keyMap {
		det.Keys = append(det.Keys, k)
		if isSensitiveKey(k) {
			det.Sensitive = append(det.Sensitive, k)
		}
	}

	return det, warnings
}
