package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func fail(format string, a ...interface{}) {
	fmt.Printf("  [ERROR] "+format+"\n", a...)
}

func warning(format string, a ...interface{}) {
	fmt.Printf("  [WARNING] "+format+"\n", a...)
}

func progress(format string, a ...interface{}) {
	fmt.Printf("  "+format+"\n", a...)
}

type HandleResult struct {
	err  error
	path string
	item interface{}
}

type HandleFunc func(c chan HandleResult, site *Website, filePath string)

func parseRepo(repoPath, outPath, owner string) (err error) {
	err = os.MkdirAll(outPath, 0775)

	if err != nil {
		fail("%s", err)
		return
	}

	fmt.Printf("Loading 'website.json'\n")
	site, err := loadSiteConfig(repoPath)

	if err != nil {
		fail("%s", err)
		return
	}

	fmt.Printf("Processing website: %s\n", site.Config.Title)

	for i, head := range site.Config.Headings {
		progress("(%d/%d) Heading: %s -> %s", i+1, len(site.Config.Headings), head.Title, head.Url)
	}

	fmt.Println("Processing Pages")
	items, err := handlePaths(site.Config.PagePaths, site, handlePage)

	if err != nil {
		warning("No pages were successfully processed")
	} else {
		for _, item := range items {
			site.Pages = append(site.Pages, item.(*Page))
		}
	}

	fmt.Println("Generating HTML Pages")
	pageTemplate, err := getPageTemplate(site)

	if err != nil {
		fail("Error loading page templates: %s", err)
	} else {
		for _, page := range site.Pages {
			progress("Generating %s", page.outputPath)
			site.CurrentPage = page
			err = generatePage(site, pageTemplate, outPath)

			if err != nil {
				warning("Failed to generate %s: %s", page.outputPath, err)
			}
		}
	}

	fmt.Println("Processing Articles")
	items, err = handlePaths(site.Config.ArticlePaths, site, handleArticle)
	if err != nil {
		warning("No articles were successfully processed")
	} else {
		for _, item := range items {
			site.Articles = append(site.Articles, item.(*Article))
		}
	}

	fmt.Println("Generating HTML Articles")
	articleTemplate, err := getArticleTemplate(site)

	if err != nil {
		fail("Error loading page templates: %s", err)
	} else {
		for _, article := range site.Articles {
			progress("Generating %s", article.permalink)
			site.CurrentArticle = article
			err = generateArticle(site, articleTemplate, outPath)

			if err != nil {
				warning("Failed to generate %s: %s", article.permalink, err)
			}
		}
	}

	fmt.Println("Copying Static Files")
	items, err = handlePaths(site.Config.StaticPaths, site, handleStatic)
	if err != nil {
		warning("No static files were successfully processed")
	} else {
		for _, file := range items {
			progress("Copying %s", file.(*StaticFile).Path)
			handleStaticPath(
				filepath.Join(site.RepoPath, file.(*StaticFile).Path),
				filepath.Join(outPath, file.(*StaticFile).Path))
		}
	}

	createHomepage(site, outPath)

	return
}

func handlePaths(paths []string, site *Website, handler HandleFunc) (list []interface{}, err error) {
	anySuccess := false
	fileCount := 0

	list = []interface{}{}

	c := make(chan HandleResult)
	for _, path := range paths {
		path = filepath.Join(site.RepoPath, path)
		progress("Searching path %s", path)

		files, err := getAllFiles(path)
		if err != nil {
			fail("Failed to process path '%s': %s", path, err)
			continue
		}

		anySuccess = true

		for _, file := range files {
			fileCount++
			go handler(c, site, file)
		}
	}

	progress("Found %d items", fileCount)
	for i := 1; i <= fileCount; i++ {
		res := <-c

		progress("(%d/%d) Processing '%s'", i, fileCount, strings.Replace(res.path, site.RepoPath, "", 1))

		if res.err != nil {
			progress("  [WARNING] Error during processing: %s", res.err)
		} else {

			if res.item != nil {
				list = append(list, res.item)
			} else {
				progress("  [ERROR] Received nil item")
			}
		}
	}

	if anySuccess {
		err = nil
	}

	return
}
