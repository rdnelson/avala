package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

type Heading struct {
	Title, Url string
}

type Website struct {
	PageTitle, Title, Subtitle string
	Headings                   []Heading
	Articles                   []*Article
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
	err = handlePages(repo, out, site)

	if err != nil {
		return err
	}

	fmt.Println("Processing articles:")
	err = handleArticles(repo, out, site)

	if err != nil {
		return err
	}

	println("Copying media files into place")

	err = handleMedia(repo, out)

	if err != nil {
		return err
	}

	fmt.Printf("Processed %d articles successfully.\n", len(site.Articles))

	return nil
}

func handleSiteDescription(repo, out string) (*Website, error) {

	file, err := os.Open(filepath.Join(repo, "website.json"))

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

func handlePages(repo, out string, site *Website) error {

	err := handlePagePath(repo, "", out, site)

	if err != nil {
		switch err.(type) {
		case *os.PathError:
			switch err.(*os.PathError).Err.(syscall.Errno) {
			case syscall.ENOENT:
				progress("Page directory not found")
			case syscall.EACCES:
				progress("Permission to pages directory denied")
				return err
			default:
				return err
			}
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
		default:
			progress(err.Error())
			return err
		}
	}

	return nil
}

func handlePagePath(repo, dir, out string, site *Website) error {
	files, err := ioutil.ReadDir(filepath.Join(repo, "pages", dir))

	if err != nil {
		return err
	}

	os.MkdirAll(out+dir, 0775)

	for _, file := range files {
		if file.IsDir() {
			err = handlePagePath(repo, filepath.Join(dir, file.Name()), out, site)
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

			_, err = parsePage(repo, filepath.Join(repo, "pages", dir, file.Name()), out, site)

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

func handleArticles(repo, out string, site *Website) error {

	err := handleArticlePath(repo, "", out, site)

	if err != nil {
		switch err.(type) {
		case *os.PathError:
			switch err.(*os.PathError).Err.(syscall.Errno) {
			case syscall.ENOENT:
				progress(err.Error())
				progress("No article directory found")
			case syscall.EACCES:
				progress("Permission to article directory denied")
				return err
			default:
				return err
			}
		case TemplateError:
			progress("Template missing: %s", err.Error())
			return err
		case ParseError:
			if err.(ParseError).IsFatal() {
				return err.(ParseError).Err()
			} else {
				progress(err.Error())
				return nil
			}
		default:
			progress(err.Error())
			return err
		}
	}

	return nil
}

func handleArticlePath(repo, dir, out string, site *Website) error {
	files, err := ioutil.ReadDir(filepath.Join(repo, "articles", dir))

	if err != nil {
		return err
	}

	os.MkdirAll(out+dir, 0775)

	for _, file := range files {
		if file.IsDir() {
			err = handleArticlePath(repo, filepath.Join(dir, file.Name()), out, site)
		} else if file.Name()[len(file.Name())-3:] == ".md" {

			filePath := filepath.Join(repo, "articles", dir, file.Name())

			if getCreatedDate(filePath).Unix() == 0 {
				continue
			}

			outPath := filepath.Join(out, filepath.FromSlash(getPermalink(filePath)+".html"))

			if dir == "" {
				progress("Parsing article: %s -> %s", file.Name(), outPath)
			} else {
				progress("Parsing article: %s -> %s", filepath.Join(dir, file.Name()), outPath)
			}

			os.MkdirAll(filepath.Dir(outPath), 0775)

			out, err := os.Create(outPath)
			defer out.Close()

			if err != nil {
				return nil
			}

			article, err := parseArticle(repo, filePath, out, site)

			if err != nil {
				if err.(ParseError).IsFatal() {
					return err.(ParseError).Err()
				} else {
					progress(err.Error())
					continue
				}
			}

			site.Articles = append(site.Articles, article)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func handleMedia(repo, out string) error {

	var count int
	err := handleMediaPath(repo, "/media", out, &count)

	if err != nil {
		switch err.(type) {
		case *os.PathError:
			switch err.(*os.PathError).Err.(syscall.Errno) {
			case syscall.ENOENT:
				progress("No media directory found")
			case syscall.EACCES:
				progress("Permission to media directory denied")
				return err
			default:
				progress(err.Error())
				return err
			}
		default:
			progress(err.Error())
			return err
		}
	}

	return nil
}

func handleMediaPath(repo, dir, out string, count *int) error {

	files, err := ioutil.ReadDir(filepath.Join(repo, dir))

	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Join(out, dir), 0775)

	offset := 0
	*count += len(files)

	for i, file := range files {

		if file.IsDir() {
			var oldCount = *count
			err = handleMediaPath(repo, filepath.Join(dir, file.Name()), out, count)
			offset += *count - oldCount

			progress("(%d/%d) Copying: %s/", i+1+offset, *count, filepath.Join(dir[1:], file.Name()))
		} else {
			progress("(%d/%d) Copying: %s", i+1+offset, *count, filepath.Join(dir[1:], file.Name()))
			err = copyFileContents(filepath.Join(repo, dir, file.Name()), filepath.Join(out, dir, file.Name()))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func createIndices(site *Website) {
	// TODO: Populate
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

func getPermalink(path string) string {
	created := getCreatedDate(path)

	file := filepath.Base(path)
	file = file[:len(file)-len(filepath.Ext(file))]

	return fmt.Sprintf("/%d/%d/%d/%s", created.Year(), created.Month(), created.Day(), file)
}
