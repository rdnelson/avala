package main

import (
	"regexp"
	"strings"

	"github.com/russross/blackfriday"
)

var paramRegex = regexp.MustCompile("\\[([^:]*):\\s*(.*)\\]")

func mdownToHtml(mdown string) string {
	return string(blackfriday.MarkdownCommon([]byte(mdown)))
}

func getTitle(contents string) string {
	lines := strings.Split(contents, "\n")

	// Look for header information
	for _, line := range lines {
		if match := paramRegex.FindStringSubmatch(line); len(match) > 0 {
			if strings.ToLower(match[1]) == "title" {
				return match[2]
			}
		} else {
			break
		}
	}

	if len(lines) > 2 && lines[1] == strings.Repeat("=", len(lines[0])) {
		return lines[0]
	}

	return ""
}

func getContents(mdown string) string {
	lines := strings.Split(mdown, "\n")

	var content string

	// If [param: value] style metadata is present
	if len(lines) > 0 && paramRegex.MatchString(lines[0]) {
		for i, line := range lines[1:] {
			if !paramRegex.MatchString(line) {
				return strings.Join(lines[i+1:], "\n")
			}
		}
	}

	if len(lines) > 2 && lines[1] == strings.Repeat("=", len(lines[0])) {
		content = strings.Join(lines[2:], "\n")
	} else {
		content = mdown
	}

	return content
}

func getOutputPath(mdown string) string {
	lines := strings.Split(mdown, "\n")

	for _, line := range lines {
		if match := paramRegex.FindStringSubmatch(line); len(match) > 0 {
			if strings.ToLower(match[1]) == "path" {
				return strings.Replace(match[2], "..", "", -1)
			}
		} else {
			break
		}
	}

	return ""
}
