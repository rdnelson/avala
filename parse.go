package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"
)

func parseRepo(repo, out string, bare bool) error {

	os.MkdirAll(out, 0775)

	var count int
	err := handleMedia(repo, "/media", out, &count)

	if err != nil {
		switch err.(*os.PathError).Err.(syscall.Errno) {
		case syscall.ENOENT:
			println("No media directory found")
			return nil
		case syscall.EACCES:
			println("Permission to media directory denied")
		}

		return err
	}

	println("Finished copying media files")

	return nil
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
			fmt.Printf("(%d/%d) Copying: %s/\n", i+1+offset, *count, dir[1:]+"/"+file.Name())

			var oldCount = *count
			handleMedia(repo, dir+"/"+file.Name(), out, count)
			offset += *count - oldCount
		} else {
			fmt.Printf("(%d/%d) Copying: %s\n", i+1+offset, *count, dir[1:]+"/"+file.Name())
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
