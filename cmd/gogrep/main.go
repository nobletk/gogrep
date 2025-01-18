package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type config struct {
	invertMatch bool
	recurse     bool
	printPath   bool
}

type application struct {
	config config
}

func main() {
	var cfg config

	pflag.BoolVarP(&cfg.recurse, "recursive", "r", false, "Read all files under each directory, recursively")
	pflag.BoolVarP(&cfg.invertMatch, "invert-match", "v", false, "Select non-matching lines")

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

	if len(paths) == 0 {
		err := app.ProcessStdin(pattern)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(2)
		}
	} else {
		if len(paths) > 1 || app.config.recurse {
			app.config.printPath = true
		}
		err := app.ProcessPaths(paths, pattern)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(2)
		}
	}

	os.Exit(0)
}
