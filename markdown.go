package main

import (
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	PageType    = iota
	ArticleType = iota
)

const (
	BlogDateFormat = "January 2 2006"
)

type Page struct {
	Site    *Website
	Title   string
	Content string
}

type Article struct {
	Site         *Website
	Title        string
	CreatedDate  time.Time
	ModifiedDate time.Time
	Content      string
	Permalink    string
}

type TemplateError interface {
	error
}

type templateError struct {
	err error
}

func (a templateError) Error() string { return a.err.Error() }

type ArticleByDate []Article

func (a ArticleByDate) Len() int           { return len(a) }
func (a ArticleByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ArticleByDate) Less(i, j int) bool { return a[i].CreatedDate.Unix() < a[j].CreatedDate.Unix() }

func (a Article) RawDate() string {
	return a.CreatedDate.Format(time.RFC3339)
}

func (a Article) FriendlyDate() string {
	return a.CreatedDate.Format(BlogDateFormat)
}

func (a Article) RawEditedDate() string {
	return a.ModifiedDate.Format(time.RFC3339)
}

func (a Article) FriendlyEditedDate() string {
	return a.ModifiedDate.Format(BlogDateFormat)
}

type ParseError interface {
	error
	IsFatal() bool
	Err() error
}

type parseError struct {
	err   error
	fatal bool
}

func (p parseError) Error() string {
	return p.err.Error()
}

func (p parseError) IsFatal() bool {
	return p.fatal
}

func (p parseError) Err() error {
	return p.err
}

func parsePage(repo, path string, out io.Writer, site *Website) (*Page, error) {
	page, err := parseMarkdown(repo, path, PageType, out, site)
	if page == nil {
		return nil, err
	}

	return page.(*Page), err
}

func parseArticle(repo, path string, out io.Writer, site *Website) (*Article, error) {
	article, err := parseMarkdown(repo, path, ArticleType, out, site)
	if article == nil {
		return nil, err
	}

	return article.(*Article), err
}

func parseMarkdown(repo, pagePath string, fileType int, out io.Writer, site *Website) (interface{}, error) {

	basePath := filepath.Join(repo, "templates")
	var contentTemplate string

	switch fileType {
	case PageType:
		contentTemplate = "page.tmpl"
	case ArticleType:
		contentTemplate = "article.tmpl"
	}

	t, err := template.ParseFiles(filepath.Join(basePath, "main.tmpl"), filepath.Join(basePath, contentTemplate))

	if err != nil {
		return nil, templateError{err}
	}

	pageContents, err := ioutil.ReadFile(pagePath)

	if err != nil {
		return nil, parseError{err, false}
	}

	lines := strings.Split(string(pageContents), "\n")

	var title string
	var content []byte

	if len(lines) > 2 && lines[1][:4] == "====" {
		title = lines[0]
		content = []byte(strings.Join(lines[2:], "\n"))
	} else {
		content = pageContents
	}

	body := string(blackfriday.MarkdownCommon(content))

	var dataObj interface{}

	switch fileType {
	case PageType:
		dataObj = &Page{site, title, body}
		err = t.Execute(out, dataObj.(*Page))

	case ArticleType:
		created := getCreatedDate(pagePath)
		modified := getModifiedDate(pagePath)
		permalink := getPermalink(pagePath)

		dataObj = &Article{site, title, created, modified, body, permalink}
		err = t.Execute(out, dataObj.(*Article))
	}

	if err != nil {
		println(err.Error())
		return nil, parseError{err, false}
	}

	return dataObj, nil
}
