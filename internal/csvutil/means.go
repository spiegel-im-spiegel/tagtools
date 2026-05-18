package csvutil

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/goark/csvdata"
	"github.com/goark/errs"
)

// LoadMeansMap loads existing means values keyed by tag from tagslist CSV.
func LoadMeansMap(path string) (meansMap map[string]string, err error) {
	meansMap = map[string]string{}

	cleanPath := filepath.Clean(path)
	f, err := os.Open(cleanPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return meansMap, nil
		}
		return nil, errs.Wrap(err, errs.WithContext("path", cleanPath))
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			if errors.Is(cerr, os.ErrClosed) {
				return
			}
			err = errs.Join(err, errs.Wrap(cerr, errs.WithContext("path", cleanPath)))
		}
	}()

	rows := csvdata.NewRows(csvdata.New(f), true)
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			if errors.Is(cerr, os.ErrClosed) {
				return
			}
			err = errs.Join(err, errs.Wrap(cerr, errs.WithContext("path", cleanPath)))
		}
	}()

	for {
		err := rows.Next()
		if err != nil {
			if errors.Is(err, io.EOF) || errs.Is(err, io.EOF) {
				break
			}
			return nil, errs.Wrap(err, errs.WithContext("path", cleanPath))
		}

		tag := rows.Column("tag")
		if tag == "" {
			continue
		}
		meansMap[tag] = rows.Column("means")
	}

	return meansMap, nil
}
