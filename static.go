package main

import (
	"bytes"
	"mime"
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

func handleStaticPath(site *Website, file, outPath string) error {
	outDir := filepath.Dir(outPath)
	os.MkdirAll(outDir, 0775)

	if !site.Config.Settings.Minify {
		return copyFileContents(file, outPath)
	}

	buf := new(bytes.Buffer)

	in, err := os.Open(file)

	if err != nil {
		return err
	}

	defer in.Close()

	err = minifier.Minify(mime.TypeByExtension(filepath.Ext(file)), buf, in)

	if err != nil {
		return copyFileContents(file, outPath)
	}

	return copyContents(buf, outPath)
}
