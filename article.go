package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ArticleDateFormat = "January 2 2006"
)

type Article struct {
	title        string
	createdDate  time.Time
	modifiedDate time.Time
	content      string
	permalink    string
	author       string
}

type ArticleByDate []*Article

func (a ArticleByDate) Len() int           { return len(a) }
func (a ArticleByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ArticleByDate) Less(i, j int) bool { return a[i].createdDate.Unix() < a[j].createdDate.Unix() }

func (a Article) IsEdited() bool {
	edited := true

	edited = edited && a.modifiedDate.Unix() != 0
	edited = edited && (a.modifiedDate.Day() != a.createdDate.Day() ||
		a.modifiedDate.Year() != a.createdDate.Year() ||
		a.modifiedDate.Month() != a.createdDate.Month())

	return edited
}

func (a Article) RawDate() string {
	return a.createdDate.Format(time.RFC3339)
}

func (a Article) FriendlyDate() string {
	return a.createdDate.Format(ArticleDateFormat)
}

func (a Article) RawEditedDate() string {
	return a.modifiedDate.Format(time.RFC3339)
}

func (a Article) FriendlyEditedDate() string {
	return a.modifiedDate.Format(ArticleDateFormat)
}

func (a Article) Summary() string {
	bPoint := strings.Index(a.content, "\n\n")
	if bPoint != -1 {
		return mdownToHtml(fmt.Sprintf("%s\n\n[Full post](%s)", a.content[:bPoint], a.permalink))
	}

	return mdownToHtml(a.content)
}

func (a Article) Title() string {
	return mdownToHtml(a.title)
}

func (a Article) Author() string {
	return a.author
}

func (a Article) Permalink() string {
	return a.permalink
}

func (a Article) Content() string {
	return mdownToHtml(a.content)
}

func handleArticles(site *Website, out string) error {

	dirs := []string{"articles"}

	for _, dir := range dirs {
		err := handleArticlePath(site, dir, "", out)

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

func handleArticlePath(site *Website, root, dir, out string) error {
	files, err := ioutil.ReadDir(filepath.Join(site.RepoPath, root, dir))

	if err != nil {
		return err
	}

	os.MkdirAll(out+dir, 0775)

	for _, file := range files {
		if file.IsDir() {
			err = handleArticlePath(site, root, filepath.Join(dir, file.Name()), out)
		} else if file.Name()[len(file.Name())-3:] == ".md" {

			filePath := filepath.Join(site.RepoPath, root, dir, file.Name())

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

			err = parseArticle(site, filePath, out)

			if err != nil {
				if err.(ParseError).IsFatal() {
					return err.(ParseError).Err()
				} else {
					progress(err.Error())
					continue
				}
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}
