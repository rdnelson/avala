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

	var owner string
	flag.StringVar(&owner, "owner", "", "User who will own the rendered website")
	flag.StringVar(&owner, "u", "", "User who will own the rendered website")

	flag.Parse()

	if len(flag.Args()) != 1 {
		println("No repository specified.")
		os.Exit(1)
	}

	if out == "" {
		println("No output path specified")
		os.Exit(2)
	}

	bare, err := isBareRepo(flag.Arg(0))

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(3)
	}

	// Initialize the post-commit hook
	if *init {
		initRepo(flag.Arg(0), out, owner, bare)
		return
	}

	parseRepo(flag.Arg(0), out, owner, bare)
}
