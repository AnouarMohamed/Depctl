package secrets

import "testing"

func TestRedact(t *testing.T) {
	got := Redact("token=abc123 password=secret", []string{"abc123", "secret"})
	if got != "token=[REDACTED] password=[REDACTED]" {
		t.Fatalf("unexpected redaction: %s", got)
	}
}

func TestIsSensitiveKey(t *testing.T) {
	for _, key := range []string{"API_KEY", "SESSION_SECRET", "DATABASE_PASSWORD"} {
		if !IsSensitiveKey(key) {
			t.Fatalf("expected %s to be sensitive", key)
		}
	}
	if IsSensitiveKey("PUBLIC_URL") {
		t.Fatal("PUBLIC_URL should not be sensitive")
	}
}
