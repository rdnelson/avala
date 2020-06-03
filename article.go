package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	ArticleDateFormat = "January 2 2006"
)

type Article struct {
	title        string
	author       string
	content      string
	permalink    string
	createdDate  time.Time
	modifiedDate time.Time
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
	bPoint := strings.Index(a.content, "\n\n\n")
	if bPoint != -1 {
		return mdownToHtml(fmt.Sprintf("%s\n\n[Full post](%s)", a.content[:bPoint], a.permalink))
	}

	return mdownToHtml(a.content)
}

func (a Article) RawTitle() string {
	return a.title
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

func handleArticle(c chan HandleResult, site *Website, articlePath string) {

	articleBytes, err := ioutil.ReadFile(articlePath)
	articleMarkdown := string(articleBytes)

	title := getTitle(articleMarkdown)
	author := getAuthor(articlePath)
	content := getContents(articleMarkdown)
	created, err := getCreatedDate(articlePath)

	if err != nil {
		c <- HandleResult{errors.New("No creation date found"), articlePath, nil}
		return
	}

	edited, _ := getModifiedDate(articlePath)
	permalink := getPermalink(articlePath)

	if err != nil {
		c <- HandleResult{err, articlePath, nil}
		return
	}

	article := &Article{
		title,
		author,
		content,
		permalink,
		created,
		edited,
	}

	c <- HandleResult{nil, articlePath, article}
}

func generateArticle(site *Website, t *template.Template, outputPath string) error {
	genPath := filepath.Join(outputPath, changeExtention(site.CurrentArticle.permalink, "html"))

	os.MkdirAll(filepath.Dir(genPath), 0775)

	out, err := os.Create(genPath)
	defer out.Close()

	if err != nil {
		return err
	}

	return t.Execute(out, site)

}
