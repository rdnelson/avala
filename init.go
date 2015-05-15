package main

import (
    "log"
    "os"
    "path/filepath"
)

func init_repo(repo string) error {
    println("Initializing hook for repo: '" + repo + "'")

    err := init_hook(repo, "/.git/hooks/post-receive")

    if err == nil {
        return nil
    }

    return init_hook(repo, "/hooks/post-receive")
}

func init_hook(repo string, hook string) error {

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

