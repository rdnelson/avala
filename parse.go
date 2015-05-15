package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
)

func parseRepo(repo, out string, bare bool) error {

    os.MkdirAll(out, 0775)

    var count int;
    err := handleMedia(repo, "/media", out, &count)

    if err != nil {
        return err
    }
    println("Finished copying media files")

    return nil
}

func handleMedia(repo, dir, out string, count *int) error {

    files, err := ioutil.ReadDir(repo + dir)

    if err != nil {
        return err
    }

    offset := 0
    *count += len(files)

    for i, file := range files {
        if file.IsDir() {
            var oldCount = *count
            handleMedia(repo, dir + "/" + file.Name(), out, count)
            offset += *count - oldCount
        }
        fmt.Printf("(%d/%d) Copying: %s\n", i + 1 + offset, *count, dir[1:] + "/" + file.Name())
    }

    return nil
}

func copyFileContents(src, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }

    defer in.Close()

    out, err := os.Create(dst)

    if err != nil {
        return
    }

    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()

    _, err = io.Copy(out, in)

    if err != nil {
        return
    }

    err = out.Sync()

    return
}
