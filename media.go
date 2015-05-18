package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func handleMedia(site *Website, out string) error {

	dirs := []string{"media", "scripts"}

	var count int

	for _, dir := range dirs {
		if _, err := os.Stat(filepath.Join(site.RepoPath, dir)); err != nil {
			continue
		}

		err := handleMediaPath(site, dir, out, &count)

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

func handleMediaPath(site *Website, dir, out string, count *int) error {

	files, err := ioutil.ReadDir(filepath.Join(site.RepoPath, dir))

	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Join(out, dir), 0775)

	offset := 0
	*count += len(files)

	for i, file := range files {

		if file.IsDir() {
			var oldCount = *count
			err = handleMediaPath(site, filepath.Join(dir, file.Name()), out, count)
			offset += *count - oldCount

			progress("(%d/%d) Copying: %s/", i+1+offset, *count, filepath.Join(dir[1:], file.Name()))
		} else {
			progress("(%d/%d) Copying: %s", i+1+offset, *count, filepath.Join(dir[1:], file.Name()))
			err = copyFileContents(filepath.Join(site.RepoPath, dir, file.Name()), filepath.Join(out, dir, file.Name()))
		}

		if err != nil {
			return err
		}
	}

	return nil
}
