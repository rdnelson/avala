package main

import (
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"text/template"
	"time"
)

type Page struct {
	Site    *Website
	Title   string
	Content string
}

type Article struct {
	Title     string
	Date      time.Time
	Content   string
	Permalink string
}

type ParseError interface {
	error
	IsFatal() bool
}

type parseError struct {
	Err   error
	Fatal bool
}

func (p parseError) Error() string {
	return p.Error()
}

func (p parseError) IsFatal() bool {
	return p.Fatal
}

func parsePage(repo, path string, out io.Writer, site *Website) ParseError {
	return parseMarkdown(repo, path, "page.tmpl", out, site)
}

func parseArticle(repo, path string, out io.Writer, site *Website) ParseError {
	return parseMarkdown(repo, path, "article.tmpl", out, site)
}

func parseMarkdown(repo, pagePath, contentTemplate string, out io.Writer, site *Website) ParseError {

	basePath := path.Join(repo, "templates")
	t, err := template.ParseFiles(path.Join(basePath, "main.tmpl"), path.Join(basePath, contentTemplate))

	if err != nil {
		return parseError{err, true}
	}

	pageContents, err := ioutil.ReadFile(pagePath)

	if err != nil {
		return parseError{err, false}
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

	page := Page{site, title, body}

	err = t.Execute(out, page)

	if err != nil {
		println(err.Error())
		return parseError{err, false}
	}

	return nil
}
