package main

import (
    "flag"
    "os"
)

func main() {
    var init = flag.Bool("init", false, "Initialize the post-commit hook for the specified repository.")

    flag.Parse()

    if len(flag.Args()) != 1 {
        println("No repository specified.")
        os.Exit(1)
    }

    // Initialize the post-commit hook
    if (*init) {
        init_repo(flag.Arg(0))
    }
}
