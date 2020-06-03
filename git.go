package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getEnvironment() []string {
	env := os.Environ()

	for i := len(env) - 1; i >= 0; i-- {
		if strings.HasPrefix(env[i], "GIT_") {
			env = append(env[:i], env[i+1:]...)
		}
	}

	return env
}

func isBareRepo(repo string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-bare-repository")
	cmd.Dir = repo
	cmd.Env = getEnvironment()

	out, err := cmd.CombinedOutput()

	if err != nil {
		println(string(out))
		return false, err
	}

	return string(out) == "true\n", nil
}

func getCreatedDate(file string) (time.Time, error) {
	return getActionTime(file, "A")
}

func getModifiedDate(file string) (time.Time, error) {
	return getActionTime(file, "M")
}

func getActionTime(file, action string) (time.Time, error) {
	if len(action) != 1 {
		return time.Unix(0, 0), errors.New("Invalid git diff action")
	}

	cmd := exec.Command("git", "log",
		"-n1",                   // One result
		"--format=%ct",          // Unix timestamp
		"--diff-filter="+action, // Filter by action
		"-i", "-E", "--invert-grep", "--grep='\\[(typo|draft)\\]'", // Search for non-typo commit
		"--", file)

	cmd.Dir = filepath.Dir(file)
	cmd.Env = getEnvironment()

	out, err := cmd.CombinedOutput()

	if err != nil {
		return time.Unix(0, 0), err
	}

	unix, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)

	if err != nil {
		return time.Unix(0, 0), err
	}

	return time.Unix(unix, 0), nil
}

func getAuthor(file string) string {
	cmd := exec.Command("git", "log", "-n1", "--format=%an", file)
	cmd.Dir = filepath.Dir(file)
	cmd.Env = getEnvironment()

	out, err := cmd.CombinedOutput()

	if err != nil {
		return "Unknown Author"
	}

	return strings.TrimSpace(string(out))
}
