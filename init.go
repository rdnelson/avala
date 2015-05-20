package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func initRepo(repo, out, owner string, bare bool) error {
	println("Initializing hook for repo: '" + repo + "'")

	return initHook(repo, out, owner, bare)
}

func initHook(repo, out, owner string, bare bool) error {

	var hook string

	if bare {
		hook = filepath.FromSlash("/hooks/post-receive")
	} else {
		hook = filepath.FromSlash("/.git/hooks/post-receive")
	}

	progress("Adding post-receive hook")
	repo, err := filepath.Abs(repo)

	if err != nil {
		return err
	}

	out, err = filepath.Abs(out)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(repo, hook), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	defer file.Close()

	if err != nil {
		progress(err.Error())
		return err
	}

	bin_path, err := filepath.Abs(os.Args[0])

	if err != nil {
		return err
	}

	var silencer string

	switch runtime.GOOS {
	case "windows":
		silencer = ">nul 2>nul"
	default:
		silencer = ">/dev/null 2>/dev/null"
	}

	_, err = file.WriteString("\n\n# Avala render hook\n\n")

	if bare {
		// Clone a checked out copy
		_, err = file.WriteString(fmt.Sprintf("cd %s\n", os.TempDir()))

		if err != nil {
			return err
		}

		_, err = file.WriteString(fmt.Sprintf("git clone %s %s\n", repo, silencer))

		if err != nil {
			return err
		}

		repo = filepath.Join(os.TempDir(), filepath.Base(repo))
	}

	if owner == "" {
		_, err = file.WriteString(fmt.Sprintf("%s -out=\"%s\" %s\n", bin_path, out, repo))
	} else {
		_, err = file.WriteString(fmt.Sprintf("%s -out=\"%s\" -owner=\"%s\" %s\n", bin_path, out, owner, repo))
	}

	if err != nil {
		return err
	}

	if bare {
		if runtime.GOOS == "windows" {
			_, err = file.WriteString(fmt.Sprintf("rmdir /s /q %s %s\n", repo, silencer))
		} else {
			_, err = file.WriteString(fmt.Sprintf("rm -rf %s %s", repo, silencer))
		}
	}

	return err
}
