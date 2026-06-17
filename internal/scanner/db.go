package scanner

import (
	"strings"

	"github.com/AnouarMohamed/Depctl/internal/types"
)

func detectDependencies(dir string, envKeys []string, runtime types.RuntimeDetection) map[string]types.Dependency {
	deps := map[string]types.Dependency{
		"postgres": {Likely: false, Confidence: 0.1},
		"redis":    {Likely: false, Confidence: 0.1},
		"mysql":    {Likely: false, Confidence: 0.1},
	}

	// 1. Env key hints
	hasPgUrl := false
	hasRedisUrl := false
	hasMysqlUrl := false

	for _, k := range envKeys {
		upperK := strings.ToUpper(k)
		if strings.Contains(upperK, "DATABASE_URL") || strings.Contains(upperK, "POSTGRES") || strings.Contains(upperK, "PG_") {
			hasPgUrl = true
		}
		if strings.Contains(upperK, "REDIS") {
			hasRedisUrl = true
		}
		if strings.Contains(upperK, "MYSQL") {
			hasMysqlUrl = true
		}
	}

	if hasPgUrl {
		deps["postgres"] = types.Dependency{Likely: true, Confidence: 0.70}
	}
	if hasRedisUrl {
		deps["redis"] = types.Dependency{Likely: true, Confidence: 0.70}
	}
	if hasMysqlUrl {
		deps["mysql"] = types.Dependency{Likely: true, Confidence: 0.70}
	}

	// 2. Runtime specific hints
	if runtime.Name == "node" {
		// Node dependencies can be checked if we scan package.json
		// We'll pass package.json check here or keep it simple.
	} else if runtime.Name == "laravel" {
		// Laravel default is MySQL or pgsql
		// Check standard laravel env keys
	}

	return deps
}
