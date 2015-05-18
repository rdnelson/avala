package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func createHomepage(site *Website, out string) error {
	last := ARTICLES_PER_PAGE - 1

	if len(site.Articles) < ARTICLES_PER_PAGE {
		last = len(site.Articles) - 1
	}

	t, err := getTemplate(site, "home.tmpl")

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

	err = t.Execute(file, site)

	if err != nil {
		return err
	}

	return err
}

func SameMonth(t, v time.Time) bool {
	return (t.Month() == v.Month()) && (t.Year() == t.Year())
}

func createIndices(site *Website, out string) (err error) {

	var name string
	finalUpdateRequired := false

	first, last := 0, 0

	progress("Creating index for %s %d", site.Articles[first].CreatedDate.Month().String(), site.Articles[first].CreatedDate.Year())
	for i, article := range site.Articles {
		name = fmt.Sprintf("%s %d", site.Articles[first].CreatedDate.Month().String(), site.Articles[first].CreatedDate.Year())
		if SameMonth(article.CreatedDate, site.Articles[first].CreatedDate) {
			last = i
			finalUpdateRequired = true
		} else {
			site.Indices = append(site.Indices, IndexRange{first, last, name})
			first = i
			progress("Creating index for %s", name)
			finalUpdateRequired = false
		}
	}
	progress("%d -> %d", first, last)

	if finalUpdateRequired {
		site.Indices = append(site.Indices, IndexRange{first, len(site.Articles) - 1, name})
	}

	err = createIndexPages(site, out)

	return
}

func createIndexPages(site *Website, out string) error {

	for _, index := range site.Indices {
		progress("Generating index page for %s", index.Name)
		for i := index.First; i <= index.Last; i++ {
			bPoint := strings.Index(site.Articles[i].Content, "\n\n")
			if bPoint != -1 {
				site.Articles[i].Content = site.Articles[i].Content[:bPoint]
			}
		}

		t, err := getTemplate(site, "index.tmpl")

		if err != nil {
			return err
		}

		outFile := getIndexPath(out, site.Articles[index.First].CreatedDate)

		out, err := os.Create(outFile)
		defer out.Close()

		if err != nil {
			return err
		}

		site.CurrentIndex = &index

		err = t.Execute(out, site)

		if err != nil {
			return err
		}
	}

	return nil
}
