package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Page struct {
	Site    *Website
	Title   string
	Content string
}

func handlePages(site *Website, out string) error {

	dirs := []string{"pages"}

	for _, dir := range dirs {
		err := handlePagePath(site, dir, "", out)

		if err != nil {

			switch {
			case os.IsNotExist(err):
				progress("%s", err.Error())
				progress("The file was likely deleted during processing")
			case os.IsPermission(err):
				progress("%s", err.Error())
				progress("Please check the file permissions on the repo")
			default:
				switch err.(type) {
				case TemplateError:
					progress("Required template missing: %s", err.Error())
					return err
				case ParseError:
					if err.(ParseError).IsFatal() {
						return err.(ParseError).Err()
					} else {
						progress(err.Error())
						return nil
					}
				}
			}
		}
	}

	return nil
}

func handlePagePath(site *Website, root, dir, out string) error {
	files, err := ioutil.ReadDir(filepath.Join(site.RepoPath, root, dir))

	if err != nil {
		return err
	}

	os.MkdirAll(out+dir, 0775)

	for _, file := range files {
		if file.IsDir() {
			err = handlePagePath(site, root, filepath.Join(dir, file.Name()), out)
		} else if file.Name()[len(file.Name())-3:] == ".md" {
			if dir == "" {
				progress("Parsing page: %s", file.Name())
			} else {
				progress("Parsing page: %s", filepath.Join(dir, file.Name()))
			}

			out, err := os.Create(filepath.Join(out, dir, file.Name()[:len(file.Name())-3]+".html"))
			defer out.Close()

			if err != nil {
				return nil
			}

			err = parsePage(site, filepath.Join(site.RepoPath, "pages", dir, file.Name()), out)

			switch err.(type) {
			case ParseError:
				if err != nil && err.(ParseError).IsFatal() {
					return err.(ParseError).Err()
				} else {
					continue
				}
			case TemplateError:
				return err
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}
