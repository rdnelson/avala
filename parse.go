package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"
	"time"
)

type Heading struct {
	Title, Url string
}

type Article struct {
	Title   string
	Date    time.Time
	Content string
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

	println("Copying media files into place")

	var count int
	err = handleMedia(repo, "/media", out, &count)

	if err != nil {
		switch err.(*os.PathError).Err.(syscall.Errno) {
		case syscall.ENOENT:
			progress("No media directory found")
			return nil
		case syscall.EACCES:
			progress("Permission to media directory denied")
		}

		return err
	}

	return nil
}

func handleSiteDescription(repo, out string) (*Website, error) {

	file, err := os.Open(repo + "/website.json")

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

func handleMedia(repo, dir, out string, count *int) error {

	files, err := ioutil.ReadDir(repo + dir)

	if err != nil {
		return err
	}

	os.MkdirAll(out+dir, 0775)

	offset := 0
	*count += len(files)

	for i, file := range files {

		if file.IsDir() {
			var oldCount = *count
			handleMedia(repo, dir+"/"+file.Name(), out, count)
			offset += *count - oldCount

			progress("(%d/%d) Copying: %s/", i+1+offset, *count, dir[1:]+"/"+file.Name())
		} else {
			progress("(%d/%d) Copying: %s", i+1+offset, *count, dir[1:]+"/"+file.Name())
			err = copyFileContents(repo+dir+"/"+file.Name(), out+dir+"/"+file.Name())
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
