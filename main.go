package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var init = flag.Bool("init", false, "Initialize the post-commit hook for the specified repository.")

	var out string
	flag.StringVar(&out, "out", "", "Path to render website into")
	flag.StringVar(&out, "o", "", "Path to render website into")

	flag.Parse()

	if len(flag.Args()) != 1 {
		println("No repository specified.")
		os.Exit(1)
	}

	bare, err := isBareRepo(flag.Arg(0))

	if err != nil {
		fmt.Println(err)
	}

	// Initialize the post-commit hook
	if *init {
		initRepo(flag.Arg(0), bare)
		return
	}

	if out != "" {
		parseRepo(flag.Arg(0), out, bare)
	} else {
		println("No output path specified")
		os.Exit(2)
	}
}
