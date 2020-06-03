package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()
	return copyContents(in, dst)
}

func copyContents(in io.Reader, dst string) (err error) {
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
	created, err := getCreatedDate(path)

	if err != nil {
		return ""
	}

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
	dir = strings.Replace(dir, "**/", "<<wilddir>>", -1)
	dir = strings.Replace(dir, "*", "[^"+string(os.PathSeparator)+"]*", -1)
	dir = strings.Replace(dir, "<<wilddir>>", ".*/?", -1)
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

func chown(owner, dir string) error {
	if owner != "" {
		usr, err := user.Lookup(owner)

		if err != nil {
			fmt.Printf("Failed to set owner to '%s'\n", owner)
			return err
		}

		err = filepath.Walk(dir, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
			uid, _ := strconv.ParseInt(usr.Uid, 10, 32)
			gid, _ := strconv.ParseInt(usr.Gid, 10, 32)
			return os.Chown(path, int(uid), int(gid))
		}))

		if err != nil {
			return err
		}
	}

	return nil
}

func changeExtention(filename, ext string) string {
	if strings.TrimSpace(filename) == "" {
		return ""
	}

	dotIdx := strings.LastIndex(filename, ".")
	sepIdx := strings.LastIndex(filename, string(os.PathSeparator))

	if dotIdx != -1 && dotIdx > sepIdx {
		return filename[:dotIdx+1] + ext
	} else {
		return filename + "." + ext
	}
}
