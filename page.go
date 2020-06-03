package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Page struct {
	title   string
	author  string
	content string

	outputPath string
}

func (a Page) Title() string {
	return mdownToHtml(a.title)
}

func (a Page) Author() string {
	return a.author
}

func (a Page) Content() string {
	return mdownToHtml(a.content)
}

func handlePage(c chan HandleResult, site *Website, pagePath string) {

	pageBytes, err := ioutil.ReadFile(pagePath)
	pageMarkdown := string(pageBytes)

	title := getTitle(pageMarkdown)
	content := getContents(pageMarkdown)
	author := getAuthor(pagePath)
	outpath := getPageOutputPath(pagePath, pageMarkdown)

	if err != nil {
		c <- HandleResult{err, pagePath, nil}
		return
	}

	page := &Page{
		title,
		author,
		content,
		outpath,
	}

	c <- HandleResult{nil, pagePath, page}
}

func getPageOutputPath(pagePath, contents string) string {

	if path := getOutputPath(contents); path != "" {
		return path
	}

	idx := strings.LastIndex(pagePath, string(os.PathSeparator))

	if idx == -1 {
		return changeExtention(pagePath, "html")
	} else {
		return changeExtention(pagePath[idx:], "html")
	}
}

func generatePage(site *Website, t *template.Template, outputPath string) error {

	genPath := filepath.Join(outputPath, site.CurrentPage.outputPath)

	os.MkdirAll(filepath.Dir(genPath), 0775)

	out, err := os.Create(genPath)
	defer out.Close()

	if err != nil {
		return err
	}

	return t.Execute(out, site)
}
