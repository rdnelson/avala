package main

import (
	"os"
	"path/filepath"
)

func createHomepage(site *Website, out string) error {
	last := ARTICLES_PER_PAGE - 1

	if len(site.Articles) < ARTICLES_PER_PAGE {
		last = len(site.Articles) - 1
	}

	t, err := getHomeTemplate(site)

	if err != nil {
		return err
	}

	outFile := filepath.Join(out, "index.html")

	file, err := os.Create(outFile)
	defer file.Close()

	if err != nil {
		return err
	}

	site.CurrentIndex = &IndexRange{0, last, "Homepage"}

	err = runTemplate(t, site, file)

	if err != nil {
		return err
	}

	return err
}
