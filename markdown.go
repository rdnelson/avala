package main

import (
	"errors"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	PageType        = iota
	ArticleType     = iota
	TemplateDir     = "templates"
	DefaultTemplate = "main.tmpl"
)

type TemplateError interface {
	error
}

type templateError error

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

func mdownToHtml(mdown string) string {
	return string(blackfriday.MarkdownCommon([]byte(mdown)))
}

func parsePage(site *Website, path string, out io.Writer) error {
	return parseMarkdown(site, path, PageType, out)
}

func parseArticle(site *Website, path string, out io.Writer) error {
	return parseMarkdown(site, path, ArticleType, out)
}

func getTemplate(site *Website, names ...string) (t *template.Template, err error) {

	if len(names) == 0 {
		return nil, errors.New("At least one template must be specified")
	}

	t, err = template.ParseFiles(filepath.Join(site.RepoPath, TemplateDir, DefaultTemplate))

	if err != nil {
		return
	}

	for _, name := range names {
		t, err = t.ParseFiles(filepath.Join(site.RepoPath, TemplateDir, name))

		if err != nil {
			err = templateError(err)
			return
		}
	}

	return
}

func parseMarkdown(site *Website, pagePath string, fileType int, out io.Writer) error {

	var contentTemplate string

	switch fileType {
	case PageType:
		contentTemplate = "page.tmpl"
	case ArticleType:
		contentTemplate = "article.tmpl"
	}

	t, err := getTemplate(site, contentTemplate)

	if err != nil {
		return templateError(err)
	}

	pageContents, err := ioutil.ReadFile(pagePath)

	if err != nil {
		return parseError{err, false}
	}

	lines := strings.Split(string(pageContents), "\n")

	var title string
	var content []byte

	if len(lines) > 2 && lines[1] == strings.Repeat("=", len(lines[0])) {
		title = lines[0]
		content = []byte(strings.Join(lines[2:], "\n"))
	} else {
		content = pageContents
	}

	body := string(content)
	author := getAuthor(pagePath)

	switch fileType {
	case PageType:
		dataObj := &Page{title, author, body}
		site.Pages = append(site.Pages, dataObj)
		site.CurrentPage = dataObj
	case ArticleType:
		created := getCreatedDate(pagePath)
		modified := getModifiedDate(pagePath)
		permalink := getPermalink(pagePath)

		dataObj := &Article{title, created, modified, body, permalink, author}
		site.Articles = append(site.Articles, dataObj)
		site.CurrentArticle = dataObj
	}

	err = t.Execute(out, site)

	site.CurrentPage = nil
	site.CurrentArticle = nil

	if err != nil {
		println(err.Error())
		return parseError{err, false}
	}

	return nil
}
