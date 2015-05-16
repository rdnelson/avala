package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func isBareRepo(repo string) (bool, error) {
	out, err := exec.Command("git", "rev-parse", "--is-bare-repository").Output()

	if err != nil {
		return false, err
	}

	return string(out) == "true", nil
}

func getCreatedDate(file string) time.Time {
	return getActionTime(file, "A")
}

func getModifiedDate(file string) time.Time {
	return getActionTime(file, "M")
}

func getActionTime(file, action string) time.Time {
	if len(action) != 1 {
		return time.Unix(0, 0)
	}

	cmd := exec.Command("git", "log", "-n1", "--format=%cI", "--diff-filter="+action, "--", file)
	cmd.Dir = filepath.Dir(file)

	out, err := cmd.CombinedOutput()

	if err != nil {
		return time.Unix(0, 0)
	}

	date, err := time.Parse(time.RFC3339, strings.TrimSpace(string(out)))

	if err != nil {
		return time.Unix(0, 0)
	}

	return date
}
