package frontmatter

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/goark/errs"
)

// Meta holds front matter fields used by tagtools.
type Meta struct {
	Date string
	Tags []string
}

type parsedMeta struct {
	Date any `yaml:"date" toml:"date" json:"date"`
	Tags any `yaml:"tags" toml:"tags" json:"tags"`
}

// ParseFile parses TOML front matter and extracts only date and tags fields.
func ParseFile(path string) (meta Meta, err error) {
	cleanPath := filepath.Clean(path)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return Meta{}, errs.Wrap(err, errs.WithContext("path", cleanPath))
	}

	var pm parsedMeta
	if _, err := frontmatter.Parse(bytes.NewReader(data), &pm); err != nil {
		return Meta{}, errs.Wrap(err, errs.WithContext("path", cleanPath))
	}

	meta = Meta{
		Date: normalizeDate(pm.Date),
		Tags: normalizeTags(pm.Tags),
	}

	return meta, nil
}

func normalizeDate(v any) string {
	switch d := v.(type) {
	case nil:
		return ""
	case string:
		for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01-02T15:04:05"} {
			if t, err := time.Parse(layout, d); err == nil {
				return t.Format("2006-01-02")
			}
		}
		if len(d) >= 10 {
			if _, err := time.Parse("2006-01-02", d[:10]); err == nil {
				return d[:10]
			}
		}
		return ""
	case time.Time:
		return d.Format("2006-01-02")
	default:
		return ""
	}
}

func normalizeTags(v any) []string {
	switch tags := v.(type) {
	case nil:
		return nil
	case []string:
		res := make([]string, 0, len(tags))
		for _, t := range tags {
			if t != "" {
				res = append(res, t)
			}
		}
		return res
	case []any:
		res := make([]string, 0, len(tags))
		for _, t := range tags {
			if s, ok := t.(string); ok && s != "" {
				res = append(res, s)
			}
		}
		return res
	case string:
		if tags == "" {
			return nil
		}
		return []string{tags}
	default:
		return nil
	}
}
