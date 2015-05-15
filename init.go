package main

import (
    "os"
    "path/filepath"
)

func initRepo(repo string, bare bool) error {
    println("Initializing hook for repo: '" + repo + "'")

    if bare {
        return initHook(repo, "/hooks/post-receive")
    } else {
        return initHook(repo, "/.git/hooks/post-receive")
    }
}

func initHook(repo string, hook string) error {

    repo, err := filepath.Abs(repo)

    if err != nil {
        return err
    }

    file, err := os.OpenFile(repo + hook, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)
    defer file.Close()

    if err != nil {
        return err
    }

    bin_path, err := filepath.Abs(os.Args[0])

    if err != nil {
        return err
    }

    _, err = file.WriteString("\n\n# Avala render hook\n\n" + bin_path + " " + repo)

    return err
}

