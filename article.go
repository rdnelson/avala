package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	ArticleDateFormat = "January 2 2006"
)

type Article struct {
	Title        string
	CreatedDate  time.Time
	ModifiedDate time.Time
	Content      string
	Permalink    string
}

type ArticleByDate []*Article

func (a ArticleByDate) Len() int           { return len(a) }
func (a ArticleByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ArticleByDate) Less(i, j int) bool { return a[i].CreatedDate.Unix() < a[j].CreatedDate.Unix() }

func (a Article) RawDate() string {
	return a.CreatedDate.Format(time.RFC3339)
}

func (a Article) FriendlyDate() string {
	return a.CreatedDate.Format(ArticleDateFormat)
}

func (a Article) RawEditedDate() string {
	return a.ModifiedDate.Format(time.RFC3339)
}

func (a Article) FriendlyEditedDate() string {
	return a.ModifiedDate.Format(ArticleDateFormat)
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
