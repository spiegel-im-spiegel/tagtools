package frontmatter

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseFile(t *testing.T) {
	t.Parallel()

	t.Run("toml quoted date and tags", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "post.md")
		content := "+++\n" +
			"date = \"2026-05-10\"\n" +
			"tags = [\"go\", \"cli\"]\n" +
			"+++\n\nbody\n"
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write fixture: %v", err)
		}

		meta, err := ParseFile(path)
		if err != nil {
			t.Fatalf("ParseFile() error = %v", err)
		}
		if meta.Date != "2026-05-10" {
			t.Fatalf("Date = %q, want %q", meta.Date, "2026-05-10")
		}
		if want := []string{"go", "cli"}; !reflect.DeepEqual(meta.Tags, want) {
			t.Fatalf("Tags = %#v, want %#v", meta.Tags, want)
		}
	})

	t.Run("toml unquoted date", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "post.md")
		content := "+++\n" +
			"date = 2026-05-11\n" +
			"tags = [\"go\"]\n" +
			"+++\n\nbody\n"
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write fixture: %v", err)
		}

		meta, err := ParseFile(path)
		if err != nil {
			t.Fatalf("ParseFile() error = %v", err)
		}
		if meta.Date != "2026-05-11" {
			t.Fatalf("Date = %q, want %q", meta.Date, "2026-05-11")
		}
	})

	t.Run("no front matter", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "post.md")
		if err := os.WriteFile(path, []byte("body only\n"), 0o600); err != nil {
			t.Fatalf("write fixture: %v", err)
		}

		meta, err := ParseFile(path)
		if err != nil {
			t.Fatalf("ParseFile() error = %v", err)
		}
		if meta.Date != "" {
			t.Fatalf("Date = %q, want empty", meta.Date)
		}
		if len(meta.Tags) != 0 {
			t.Fatalf("Tags = %#v, want empty", meta.Tags)
		}
	})
}
