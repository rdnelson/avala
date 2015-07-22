package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	ARTICLES_PER_PAGE = 4
)

type Heading struct {
	Title, Url string
}

type IndexRange struct {
	First, Last int
	Name        string
}

type WebsiteConfig struct {
	Title    string
	Subtitle string

	ArticlePaths    []string
	MediaPaths      []string
	PagePaths       []string
	ScriptPaths     []string
	StylesheetPaths []string
	TemplatePaths   []string
}

type Website struct {
	Config WebsiteConfig
	// Internal values
	Headings       []Heading
	Articles       []*Article
	Pages          []*Page
	CurrentPage    *Page
	CurrentArticle *Article
	CurrentIndex   *IndexRange
	Indices        []IndexRange
	RepoPath       string
}

func (a *Website) ActiveArticles() []*Article {
	if a.CurrentIndex == nil {
		return nil
	}

	return a.Articles[a.CurrentIndex.First : a.CurrentIndex.Last+1]
}

func loadSiteConfig(repo string) (*Website, error) {

	file, err := os.Open(filepath.Join(repo, "website.json"))

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	site := new(Website)
	site.RepoPath = repo

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &site.Config)

	if err != nil {
		return nil, err
	}

	return site, err
}
