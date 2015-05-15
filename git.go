package main

import (
	"os/exec"
)

func isBareRepo(repo string) (bool, error) {
	out, err := exec.Command("git", "rev-parse", "--is-bare-repository").Output()

	if err != nil {
		return false, err
	}

	return string(out) == "true", nil
}
