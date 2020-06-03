package main

import (
	"os"
	"path/filepath"
	"strings"
)

type StaticFile struct {
	Path string
}

func handleStatic(c chan HandleResult, site *Website, staticFile string) {
	c <- HandleResult{
		nil,
		staticFile,
		&StaticFile{
			strings.Replace(staticFile, site.RepoPath, "", 1),
		},
	}
}

func handleStaticPath(file, outPath string) error {
	outDir := filepath.Dir(outPath)

	os.MkdirAll(outDir, 0775)
	return copyFileContents(file, outPath)
}
