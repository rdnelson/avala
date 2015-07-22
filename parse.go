package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
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

func parseRepo(repoPath, outPath, owner string) (err error) {
	os.MkdirAll(outPath, 0775)

	fmt.Printf("Loading 'website.json'\n")
	site, err := loadSiteConfig(repoPath)

	if err != nil {
		fail("%s", err)
		return
	}

	fmt.Printf("Processing website: %s\n", site.Config.Title)

	for i, head := range site.Headings {
		progress("(%d/%d) Heading: %s -> %s", i+1, len(site.Headings), head.Title, head.Url)
	}

	fmt.Println("Processing Pages")
	for _, path := range site.Config.PagePaths {
		path = filepath.Join(repoPath, path)
		files, err := getAllFiles(path)
		if err != nil {
			fail("Failed to process page path '%s': %s", path, err)
			continue
		}

		for _, file := range files {
			progress("Found page '%s'", file)
		}
	}

	fmt.Println("Processing Articles")
	for _, path := range site.Config.ArticlePaths {
		path = filepath.Join(repoPath, path)
		files, err := getAllFiles(path)
		if err != nil {
			fail("Failed to process article path '%s': %s", path, err)
			continue
		}

		for _, file := range files {
			progress("Found article '%s'", file)
		}
	}

	return
}

func parseRepoOld(repo, out, owner string) error {

	os.MkdirAll(out, 0775)

	site, err := loadSiteConfig(repo)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if site == nil {
		println("No website.json found, cannot generate website.")
		return nil
	}

	fmt.Printf("Processing website: %s\n", site.Config.Title)

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

	if len(site.Config.MediaPaths) == 0 {
		progress("No media directories")
	} else {
		err = handleMedia(site, out)

		if err != nil {
			return err
		}
	}

	if owner != "" {
		usr, err := user.Lookup(owner)

		if err != nil {
			fmt.Printf("Failed to set owner to '%s'\n", owner)
			return err
		}

		err = filepath.Walk(out, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
			uid, _ := strconv.ParseInt(usr.Uid, 10, 32)
			gid, _ := strconv.ParseInt(usr.Gid, 10, 32)
			return os.Chown(path, int(uid), int(gid))
		}))

		if err != nil {
			return err
		}
	}

	return nil
}
