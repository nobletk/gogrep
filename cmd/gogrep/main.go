package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/pflag"
)

type config struct {
	caseInsensitive bool
	invertMatch     bool
	recurse         bool
	printPath       bool
	matchCount      int
}

type application struct {
	config config
}

func main() {
	var cfg config

	pflag.BoolVarP(&cfg.recurse, "recursive", "r", false, "Read all files under each directory, recursively")
	pflag.BoolVarP(&cfg.invertMatch, "invert-match", "v", false, "Select non-matching lines")
	pflag.BoolVarP(&cfg.caseInsensitive, "case-insensitive", "i", false, "Case insensitive search")

	pflag.Usage = func() {
		var buf bytes.Buffer

		buf.WriteString("Usage:\n")
		buf.WriteString(" gogrep [OPTION...] PATTERNS [FILE...]\n")

		fmt.Fprintf(os.Stderr, buf.String())
		pflag.PrintDefaults()
	}

	pflag.Parse()

	/*
		if len(pflag.Args()) > 2 {
			pflag.Usage()
			os.Exit(2)
		}
	*/

	app := &application{
		config: cfg,
	}

	pattern := pflag.Arg(0)
	paths := pflag.Args()[1:]

	if app.config.caseInsensitive {
		pattern = "(?i)" + pattern
	}

	regexPattern, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error compiling regex: %v\n", err)
		os.Exit(1)
	}

	if len(paths) == 0 {
		err := app.ProcessStdin(regexPattern)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if len(paths) > 1 || app.config.recurse {
			app.config.printPath = true
		}
		err := app.ProcessPaths(paths, regexPattern)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
	}

	if app.config.matchCount == 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
