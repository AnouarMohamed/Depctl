package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDotenv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	err := os.WriteFile(path, []byte(`
# comment
DATABASE_URL="postgres://example"
TOKEN='abc123'
PUBLIC_URL=https://example.com
DATABASE_URL=postgres://override
`), 0600)
	if err != nil {
		t.Fatal(err)
	}

	entries, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	rendered := AsDotenv(entries)
	if rendered == "" || rendered[0] == '#' {
		t.Fatalf("unexpected dotenv render: %q", rendered)
	}
	keys := Keys(entries)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
}
