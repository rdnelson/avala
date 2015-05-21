package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func handleMedia(site *Website, out string) error {

	var count int
	offset := 0

	for _, dir := range site.Media {
		if _, err := os.Stat(filepath.Join(site.RepoPath, dir)); err != nil {
			continue
		}

		oldCount := count
		err := handleMediaPath(site, dir, out, &count, &offset)
		offset += count - oldCount

		if err != nil {
			switch {
			case os.IsNotExist(err):
				progress("%s", err.Error())
			case os.IsPermission(err):
				progress("%s", err.Error())
			}
		}
	}

	return nil
}

func handleMediaPath(site *Website, dir, out string, count *int, offset *int) error {

	files, err := ioutil.ReadDir(filepath.Join(site.RepoPath, dir))

	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Join(out, dir), 0775)

	for _, file := range files {
		if !file.IsDir() {
			*count++
		}
	}

	for i, file := range files {

		if file.IsDir() {
			var oldCount = *count
			err = handleMediaPath(site, filepath.Join(dir, file.Name()), out, count, offset)
			*offset += *count - oldCount - 1
		} else {
			progress("(%d/%d) Copying: %s", i+1+*offset, *count, filepath.Join(dir, file.Name()))
			err = copyFileContents(filepath.Join(site.RepoPath, dir, file.Name()), filepath.Join(out, dir, file.Name()))
		}

		if err != nil {
			return err
		}
	}

	return nil
}
