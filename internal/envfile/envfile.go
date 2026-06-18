package envfile

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Entry is a parsed KEY=VALUE pair from a dotenv file.
type Entry struct {
	Key   string
	Value string
}

var lineRE = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)\s*$`)

// Parse reads a dotenv file and returns key/value entries. It does not expand variables.
func Parse(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open env file %s: %w", path, err)
	}
	defer file.Close()

	var entries []Entry
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		match := lineRE.FindStringSubmatch(line)
		if len(match) != 3 {
			continue
		}

		key := match[1]
		value := strings.TrimSpace(match[2])
		value = strings.Trim(value, `"`)
		value = strings.Trim(value, `'`)
		if seen[key] {
			for i := range entries {
				if entries[i].Key == key {
					entries[i].Value = value
					break
				}
			}
			continue
		}

		seen[key] = true
		entries = append(entries, Entry{Key: key, Value: value})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read env file %s: %w", path, err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries, nil
}

// Keys returns only dotenv keys, safe for plan files and reports.
func Keys(entries []Entry) []string {
	keys := make([]string, 0, len(entries))
	for _, entry := range entries {
		keys = append(keys, entry.Key)
	}
	return keys
}

// AsDotenv renders entries back to KEY=VALUE lines for provider stdin.
func AsDotenv(entries []Entry) string {
	var b strings.Builder
	for _, entry := range entries {
		b.WriteString(entry.Key)
		b.WriteString("=")
		b.WriteString(entry.Value)
		b.WriteString("\n")
	}
	return b.String()
}
