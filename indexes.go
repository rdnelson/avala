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

	err = t.Execute(file, site)

	if err != nil {
		return err
	}

	return err
}

/*func SameMonth(t, v time.Time) bool {
	return (t.Month() == v.Month()) && (t.Year() == t.Year())
}

func createIndices(site *Website, out string) (err error) {

	// If there are no articles, there are no indexes
	if len(site.Articles) == 0 {
		return nil
	}

	var name string
	finalUpdateRequired := false

	first, last := 0, 0

	name = fmt.Sprintf("%s %d", site.Articles[first].createdDate.Month().String(), site.Articles[first].createdDate.Year())
	progress("Creating index for %s", name)
	for i, article := range site.Articles {
		if SameMonth(article.createdDate, site.Articles[first].createdDate) {
			last = i
			finalUpdateRequired = true
		} else {
			site.Indices = append(site.Indices, IndexRange{first, last, name})
			first = i
			name = fmt.Sprintf("%s %d", site.Articles[first].createdDate.Month().String(), site.Articles[first].createdDate.Year())
			progress("Creating index for %s", name)
			finalUpdateRequired = false
		}
	}

	if finalUpdateRequired {
		site.Indices = append(site.Indices, IndexRange{first, len(site.Articles) - 1, name})
	}

	err = createIndexPages(site, out)

	return
}

func createIndexPages(site *Website, out string) error {

	for _, index := range site.Indices {
		progress("Generating index page for %s", index.Name)

		t, err := getTemplate(site, "index.tmpl")

		if err != nil {
			return err
		}

		outFile := getIndexPath(out, site.Articles[index.First].createdDate)

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
}*/
