package toptags

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goark/errs"

	"github.com/spiegel-im-spiegel/tagtools/internal/contents"
)

// Config represents toptags command options.
type Config struct {
	ContentDir string `pflag:"content-dir,c,content directory to scan" env:"CONTENT_DIR"`
	Out        string `pflag:"out,o,output JSON path" env:"OUT"`
	TopN       int    `pflag:"top-n,n,number of top tags" env:"TOP_N"`
	Today      string `pflag:"today,t,override today date (YYYY-MM-DD)" env:"TODAY"`
	Window     string `pflag:"window,w,window duration (e.g. 1y, 6m, 90d, 1y2m10d)" env:"WINDOW"`
}

func DefaultConfig() Config {
	return Config{
		ContentDir: "content",
		Out:        "data/toptags.json",
		TopN:       15,
		Today:      "",
		Window:     "1y",
	}
}

type countItem struct {
	Tag   string
	Count int
}

// Run executes the toptags workflow.
func Run(cfg Config) (err error) {
	today, err := resolveToday(cfg.Today)
	if err != nil {
		return errs.Wrap(err)
	}
	cutoff, err := resolveCutoff(today, cfg.Window)
	if err != nil {
		return errs.Wrap(err)
	}

	counts := map[string]int{}
	err = contents.WalkMarkdownMeta(cfg.ContentDir, func(_ string, meta contents.Meta) error {
		if meta.Date == "" {
			return nil
		}
		postDate, err := time.Parse("2006-01-02", meta.Date)
		if err != nil {
			return nil
		}
		if postDate.Before(cutoff) || postDate.After(today) {
			return nil
		}
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

	n := cfg.TopN
	if n < 0 {
		n = 0
	}
	if n > len(items) {
		n = len(items)
	}

	tags := make([]string, 0, n)
	for i := 0; i < n; i++ {
		tags = append(tags, items[i].Tag)
	}
	sort.Strings(tags)

	if err := os.MkdirAll(filepath.Dir(cfg.Out), 0o750); err != nil {
		return errs.Wrap(err, errs.WithContext("path", cfg.Out))
	}
	outDir := filepath.Dir(cfg.Out)
	tmp, err := os.CreateTemp(outDir, "toptags-*.json")
	if err != nil {
		return errs.Wrap(err, errs.WithContext("path", outDir))
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = os.Remove(tmpPath)
		if tmp == nil {
			return
		}
		if cerr := tmp.Close(); cerr != nil {
			err = errs.Join(err, errs.Wrap(cerr, errs.WithContext("path", tmpPath)))
		}
	}()

	if _, err := tmp.WriteString("["); err != nil {
		return errs.Wrap(err)
	}
	for i, t := range tags {
		if i > 0 {
			if _, err := tmp.WriteString(", "); err != nil {
				return errs.Wrap(err)
			}
		}
		if _, err := tmp.WriteString(strconv.Quote(t)); err != nil {
			return errs.Wrap(err)
		}
	}
	if _, err := tmp.WriteString("]\n"); err != nil {
		return errs.Wrap(err)
	}
	if err := tmp.Close(); err != nil {
		return errs.Wrap(err, errs.WithContext("path", tmpPath))
	}
	tmp = nil
	if err := os.Rename(tmpPath, cfg.Out); err != nil {
		return errs.Wrap(err, errs.WithContext("path", cfg.Out))
	}

	return nil
}

func resolveToday(s string) (time.Time, error) {
	if s == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errs.Wrap(err, errs.WithContext("today", s))
	}
	return d, nil
}

func resolveCutoff(today time.Time, window string) (time.Time, error) {
	years, months, days, err := parseWindow(window)
	if err != nil {
		return time.Time{}, errs.Wrap(err, errs.WithContext("window", window))
	}
	return today.AddDate(-years, -months, -days), nil
}

func parseWindow(s string) (int, int, int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, 0, errs.New("window is empty")
	}

	years, months, days := 0, 0, 0
	for i := 0; i < len(s); {
		start := i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		if start == i || i >= len(s) {
			return 0, 0, 0, errs.New("invalid window format")
		}
		n, err := strconv.Atoi(s[start:i])
		if err != nil {
			return 0, 0, 0, errs.Wrap(err)
		}
		switch s[i] {
		case 'y':
			years += n
		case 'm':
			months += n
		case 'd':
			days += n
		default:
			return 0, 0, 0, errs.New("invalid window unit", errs.WithContext("unit", string(s[i])))
		}
		i++
	}
	return years, months, days, nil
}
