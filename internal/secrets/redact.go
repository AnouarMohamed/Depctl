package secrets

import "strings"

var sensitiveTokens = []string{
	"SECRET",
	"TOKEN",
	"PASSWORD",
	"PASS",
	"PRIVATE",
	"KEY",
	"API_KEY",
	"DB_PASS",
	"JWT",
	"SESSION",
	"CREDENTIAL",
	"ACCESS_KEY",
}

// IsSensitiveKey reports whether a key name likely contains a secret value.
func IsSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, token := range sensitiveTokens {
		if strings.Contains(upper, token) {
			return true
		}
	}
	return false
}

// Redact replaces known secret values in command output.
func Redact(text string, values []string) string {
	redacted := text
	for _, value := range values {
		if value == "" {
			continue
		}
		redacted = strings.ReplaceAll(redacted, value, "[REDACTED]")
	}
	return redacted
}
