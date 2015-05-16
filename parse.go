package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type Heading struct {
	Title, Url string
}

type Website struct {
	PageTitle, Title, Subtitle string
	Headings                   []Heading
	Articles                   []Article
}

func progress(format string, a ...interface{}) {
	fmt.Printf("  "+format+"\n", a...)
}

func parseRepo(repo, out string, bare bool) error {

	os.MkdirAll(out, 0775)

	site, err := handleSiteDescription(repo, out)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Processing website: %s\n", site.Title)

	for i, head := range site.Headings {
		progress("(%d/%d) Heading: %s -> %s", i+1, len(site.Headings), head.Title, head.Url)
	}

	fmt.Println("Processing pages:")
	err = handlePages(repo, "", out, site)

	if err != nil {
		switch err.(type) {
		case *os.PathError:
			switch err.(*os.PathError).Err.(syscall.Errno) {
			case syscall.ENOENT:
				progress(err.Error())
			case syscall.EACCES:
				progress("Permission to page directory denied")
				return err
			default:
				return err
			}
		default:
			progress(err.Error())
			return err
		}
	}

	println("Copying media files into place")

	var count int
	err = handleMedia(repo, "/media", out, &count)

	if err != nil {
		switch err.(*os.PathError).Err.(syscall.Errno) {
		case syscall.ENOENT:
			progress("No media directory found")
		case syscall.EACCES:
			progress("Permission to media directory denied")
			return err
		default:
			return err
		}
	}

	return nil
}

func handleSiteDescription(repo, out string) (*Website, error) {

	file, err := os.Open(path.Join(repo, "website.json"))

	if err != nil {
		return nil, err
	}

	site := new(Website)

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &site)

	return site, err
}

func handlePages(repo, dir, out string, site *Website) error {
	files, err := ioutil.ReadDir(path.Join(repo, "pages", dir))

	if err != nil {
		return err
	}

	os.MkdirAll(out+dir, 0775)

	for _, file := range files {
		if file.IsDir() {
			err = handlePages(repo, path.Join(dir, file.Name()), out, site)
		} else if file.Name()[len(file.Name())-3:] == ".md" {
			if dir == "" {
				progress("Parsing page: %s", file.Name())
			} else {
				progress("Parsing page: %s", path.Join(dir, file.Name()))
			}

			out, err := os.Create(path.Join(out, dir, file.Name()[:len(file.Name())-3]+".html"))
			defer out.Close()

			if err != nil {
				return nil
			}

			err = parsePage(repo, path.Join(repo, "pages", dir, file.Name()), out, site)

			if err != nil && err.(ParseError).IsFatal() {
				return err.(parseError).Err
			} else {
				continue
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func handleMedia(repo, dir, out string, count *int) error {

	files, err := ioutil.ReadDir(path.Join(repo, dir))

	if err != nil {
		return err
	}

	os.MkdirAll(path.Join(out, dir), 0775)

	offset := 0
	*count += len(files)

	for i, file := range files {

		if file.IsDir() {
			var oldCount = *count
			err = handleMedia(repo, path.Join(dir, file.Name()), out, count)
			offset += *count - oldCount

			progress("(%d/%d) Copying: %s/", i+1+offset, *count, path.Join(dir[1:], file.Name()))
		} else {
			progress("(%d/%d) Copying: %s", i+1+offset, *count, path.Join(dir[1:], file.Name()))
			err = copyFileContents(path.Join(repo, dir, file.Name()), path.Join(out, dir, file.Name()))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

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
