package contents

import (
	"os"
	"path/filepath"

	"github.com/goark/errs"

	"github.com/spiegel-im-spiegel/tagtools/internal/frontmatter"
)

// Meta holds front matter fields used by content walkers.
type Meta = frontmatter.Meta

// WalkMarkdownMeta walks markdown files under contentDir and calls fn with parsed front matter.
func WalkMarkdownMeta(contentDir string, fn func(path string, meta Meta) error) error {
	err := filepath.WalkDir(contentDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return errs.Wrap(err, errs.WithContext("path", path))
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		meta, err := frontmatter.ParseFile(path)
		if err != nil {
			return errs.Wrap(err)
		}
		if err := fn(path, meta); err != nil {
			return errs.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}
