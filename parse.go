package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

type Website struct {
	PageTitle, Title, Subtitle string
	Headings                   []Heading
	Articles                   []*Article
	Pages                      []*Page
	CurrentPage                *Page
	CurrentArticle             *Article
	CurrentIndex               *IndexRange
	Indices                    []IndexRange
	RepoPath                   string
	Media                      []string
}

func (a *Website) ActiveArticles() []*Article {
	if a.CurrentIndex == nil {
		return nil
	}

	return a.Articles[a.CurrentIndex.First : a.CurrentIndex.Last+1]
}

func progress(format string, a ...interface{}) {
	fmt.Printf("  "+format+"\n", a...)
}

func parseRepo(repo, out string, bare bool) error {

	os.MkdirAll(out, 0775)

	site, err := handleSiteDescription(repo, out)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if site == nil {
		println("No website.json found, cannot generate website.")
		return nil
	}

	fmt.Printf("Processing website: %s\n", site.Title)

	for i, head := range site.Headings {
		progress("(%d/%d) Heading: %s -> %s", i+1, len(site.Headings), head.Title, head.Url)
	}

	fmt.Println("Processing pages")
	err = handlePages(site, out)

	if err != nil {
		return err
	}

	fmt.Println("Processing articles")
	err = handleArticles(site, out)

	if err != nil {
		return err
	}

	// Sort the articles to newest first
	sort.Sort(sort.Reverse(ArticleByDate(site.Articles)))

	println("Creating homepage")
	err = createHomepage(site, out)

	if err != nil {
		progress(err.Error())
		return err
	}
	progress("Successfully generated homepage")

	println("Creating archive index pages")
	err = createIndices(site, out)

	if err != nil {
		progress(err.Error())
		return err
	}

	println("Copying media files into place")

	if len(site.Media) == 0 {
		progress("No media directories")
	} else {
		err = handleMedia(site, out)

		if err != nil {
			return err
		}
	}

	return nil
}

func handleSiteDescription(repo, out string) (*Website, error) {

	file, err := os.Open(filepath.Join(repo, "website.json"))

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	site := new(Website)

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &site)

	site.RepoPath = repo

	return site, err
}
