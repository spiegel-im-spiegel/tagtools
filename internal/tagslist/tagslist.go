package tagslist

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goark/errs"

	"github.com/spiegel-im-spiegel/tagtools/internal/contents"
	"github.com/spiegel-im-spiegel/tagtools/internal/csvutil"
)

// Config represents tagslist command options.
type Config struct {
	ContentDir       string `pflag:"content-dir,c,content directory to scan" env:"CONTENT_DIR"`
	Out              string `pflag:"out,o,output CSV path" env:"OUT"`
	InheritMeansFrom string `pflag:"inherit-means-from,m,CSV path to inherit means from" env:"INHERIT_MEANS_FROM"`
}

func DefaultConfig() Config {
	return Config{
		ContentDir:       "content",
		Out:              ".github/workflows/tagslist.csv",
		InheritMeansFrom: ".github/workflows/tagslist.csv",
	}
}

type countItem struct {
	Tag   string
	Count int
}

// Run executes the tagslist workflow.
func Run(cfg Config) error {
	meansMap, err := csvutil.LoadMeansMap(cfg.InheritMeansFrom)
	if err != nil {
		return errs.Wrap(err)
	}

	counts := map[string]int{}
	err = contents.WalkMarkdownMeta(cfg.ContentDir, func(_ string, meta contents.Meta) error {
		for _, t := range meta.Tags {
			counts[t]++
		}
		return nil
	})
	if err != nil {
		return errs.Wrap(err)
	}

	items := make([]countItem, 0, len(counts))
	for tag, count := range counts {
		items = append(items, countItem{Tag: tag, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Tag < items[j].Tag
	})

	if err := os.MkdirAll(filepath.Dir(cfg.Out), 0o750); err != nil {
		return errs.Wrap(err, errs.WithContext("path", cfg.Out))
	}
	outDir := filepath.Dir(cfg.Out)
	tmp, err := os.CreateTemp(outDir, "tagslist-*.csv")
	if err != nil {
		return errs.Wrap(err)
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := fmt.Fprintln(tmp, "tag,count,means"); err != nil {
		_ = tmp.Close()
		return errs.Wrap(err)
	}
	for _, it := range items {
		means := meansMap[it.Tag]
		line := fmt.Sprintf("\"%s\",%d,\"%s\"\n", escapeCSV(it.Tag), it.Count, escapeCSV(means))
		if _, err := tmp.WriteString(line); err != nil {
			_ = tmp.Close()
			return errs.Wrap(err)
		}
	}
	if err := tmp.Close(); err != nil {
		return errs.Wrap(err)
	}

	if err := os.Rename(tmpPath, cfg.Out); err != nil {
		return errs.Wrap(err, errs.WithContext("path", cfg.Out))
	}
	return nil
}

func escapeCSV(s string) string {
	return strings.ReplaceAll(s, `"`, `""`)
}
