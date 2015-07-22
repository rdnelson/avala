package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func getAllFiles(dir string) (files []string, err error) {

	files = []string{}

	base := dir

	if idx := strings.Index(dir, "*"); idx > -1 {
		base = dir[:idx]
	}

	dir = strings.Replace(dir, ".", "\\.", -1)
	dir = strings.Replace(dir, "**", "<<wilddir>>", -1)
	dir = strings.Replace(dir, "*", "[^"+string(os.PathSeparator)+"]*", -1)
	dir = strings.Replace(dir, "<<wilddir>>", ".*", -1)
	dir = "^" + dir + "$"

	pattern, err := regexp.Compile(dir)

	if err != nil {
		warning("Regex pattern failed: %s", err)
		return
	}

	err = filepath.Walk(base, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			warning("Walking '%s' failed: %s", path, e)
			return nil
		}

		if !info.IsDir() {
			if pattern.MatchString(path) {
				files = append(files, path)
			}
		}

		return nil
	})

	return
}
