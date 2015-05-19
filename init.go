package main

import (
	"os"
	"path/filepath"
	"runtime"
)

func initRepo(repo, out string, bare bool) error {
	println("Initializing hook for repo: '" + repo + "'")

	if bare {
		return initHook(repo, out, true)
	} else {
		return initHook(repo, out, false)
	}
}

func initHook(repo, out string, bare bool) error {

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
		_, err = file.WriteString("cd " + os.TempDir() + "\n")

		if err != nil {
			return err
		}

		_, err = file.WriteString("git clone " + repo + " " + silencer + "\n")

		if err != nil {
			return err
		}

		repo = filepath.Join(os.TempDir(), filepath.Base(repo))
	}

	_, err = file.WriteString(bin_path + " -out=\"" + out + "\" " + repo + "\n")

	if err != nil {
		return err
	}

	if bare {
		if runtime.GOOS == "windows" {
			_, err = file.WriteString("rmdir /s /q " + repo + " " + silencer + "\n")
		} else {
			_, err = file.WriteString("rm -rf " + repo + " " + silencer + "\n")
		}
	}

	return err
}
