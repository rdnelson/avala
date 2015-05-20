package main

import (
	"os"
	"os/exec"
	"path/filepath"
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
	cmd := exec.Command("git", "-C", repo, "rev-parse", "--is-bare-repository")
	cmd.Env = getEnvironment()

	out, err := cmd.CombinedOutput()

	if err != nil {
		return false, err
	}

	return string(out) == "true\n", nil
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
	cmd.Env = getEnvironment()

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
