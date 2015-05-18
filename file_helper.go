package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}

	defer in.Close()

	out, err := os.Create(dst)

	if err != nil {
		return
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(out, in)

	if err != nil {
		return
	}

	err = out.Sync()

	return
}

func getPermalink(path string) string {
	created := getCreatedDate(path)

	file := filepath.Base(path)
	file = file[:len(file)-len(filepath.Ext(file))]

	return fmt.Sprintf("/%d/%d/%d/%s", created.Year(), created.Month(), created.Day(), file)
}

func getIndexPath(out string, date time.Time) string {
	return filepath.FromSlash(fmt.Sprintf("%s/%d/%d/index.html", out, date.Year(), date.Month()))
}
