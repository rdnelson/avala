package main

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

var minifier = buildMinifier()

func buildMinifier() *minify.M {
	m := minify.New()

	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)

	return m
}

func runTemplate(t *template.Template, site *Website, out io.Writer) error {
	if !site.Config.Settings.Minify {
		return t.Execute(out, site)
	}

	buf := new(bytes.Buffer)
	err := t.Execute(buf, site)
	if err != nil {
		return err
	}

	mb, err := minifier.Bytes("text/html", buf.Bytes())

	if err != nil {
		return err
	}

	_, err = out.Write(mb)
	return err
}

func getArticleTemplate(site *Website) (t *template.Template, err error) {

	paths := make([]string, len(site.Config.GlobalTemplates))
	copy(paths, site.Config.GlobalTemplates)

	paths = append(paths, site.Config.ArticleTemplates...)

	for i := 0; i < len(paths); i++ {
		paths[i] = filepath.Join(site.RepoPath, paths[i])
	}

	return loadTemplates(paths)
}

func getPageTemplate(site *Website) (t *template.Template, err error) {

	paths := make([]string, len(site.Config.GlobalTemplates))
	copy(paths, site.Config.GlobalTemplates)

	paths = append(paths, site.Config.PageTemplates...)

	for i := 0; i < len(paths); i++ {
		paths[i] = filepath.Join(site.RepoPath, paths[i])
	}

	return loadTemplates(paths)
}

func getHomeTemplate(site *Website) (t *template.Template, err error) {

	paths := make([]string, len(site.Config.GlobalTemplates))
	copy(paths, site.Config.GlobalTemplates)

	paths = append(paths, site.Config.HomeTemplates...)

	for i := 0; i < len(paths); i++ {
		paths[i] = filepath.Join(site.RepoPath, paths[i])
	}

	return loadTemplates(paths)
}

func getIndexTemplate(site *Website) (t *template.Template, err error) {

	paths := make([]string, len(site.Config.GlobalTemplates))
	copy(paths, site.Config.GlobalTemplates)

	paths = append(paths, site.Config.IndexTemplates...)

	for i := 0; i < len(paths); i++ {
		paths[i] = filepath.Join(site.RepoPath, paths[i])
	}

	return loadTemplates(paths)
}

func loadTemplates(paths []string) (t *template.Template, err error) {

	if len(paths) == 0 {
		return nil, errors.New("At least one template must be specified")
	}

	t = template.New(filepath.Base(paths[0])).Funcs(sprig.TxtFuncMap())

	t, err = t.ParseFiles(paths[0])

	if err != nil {
		return
	}

	for i := 1; i < len(paths); i++ {
		t, err = t.ParseFiles(paths[i])

		if err != nil {
			return
		}
	}

	return
}
