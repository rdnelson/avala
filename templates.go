package main

import (
	"errors"
	"path/filepath"
	"text/template"
)

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

	t, err = template.ParseFiles(paths[0])

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
