package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func (app *application) ProcessStdin(pattern *regexp.Regexp) error {
	const bufferSize = 4096
	buffer := make([]byte, bufferSize)

	var leftover []byte

	for {
		n, err := os.Stdin.Read(buffer)
		if n > 0 {
			data := append(leftover, buffer[:n]...)

			lines := bytes.Split(data, []byte("\n"))

			leftover = lines[len(lines)-1]

			for _, line := range lines[:len(lines)-1] {
				if app.config.invertMatch {
					app.invertMatchStdin(line, pattern)
					continue
				}
				app.patternMatchStdin(line, pattern)
			}
		}
		if err == io.EOF {
			if len(leftover) > 0 {
				if app.config.invertMatch {
					app.invertMatchStdin(leftover, pattern)
					continue
				}
				app.patternMatchStdin(leftover, pattern)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading stdin: %W", err)
		}
	}

	return nil
}

func (app *application) ProcessPaths(paths []string, pattern *regexp.Regexp) error {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			continue
		}

		if info.IsDir() {
			if !app.config.recurse {
				fmt.Printf("gogrep: %s: Is a directory\n", path)
				continue
			}

			err := filepath.WalkDir(path, func(subPath string, d os.DirEntry, err error) error {
				if err != nil {
					fmt.Printf("Error accessing path %q: %v\n", subPath, err)
					return nil
				}
				if d.IsDir() {
					return nil
				}
				return app.processFile(subPath, pattern)
			})

			if err != nil {
				fmt.Printf("Error processing directory %q: %v\n", path, err)
			}
		}
		if !info.IsDir() {
			err := app.processFile(path, pattern)
			if err != nil {
				fmt.Printf("Error processing file %q: %v\n", path, err)
			}
		}
	}

	return nil
}

func (app *application) processFile(path string, pattern *regexp.Regexp) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	const bufferSize = 4096
	buffer := make([]byte, bufferSize)

	var leftover []byte

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			data := append(leftover, buffer[:n]...)

			lines := bytes.Split(data, []byte("\n"))

			leftover = lines[len(lines)-1]

			for _, line := range lines[:len(lines)-1] {
				if app.config.invertMatch {
					app.invertMatch(path, line, pattern)
					continue
				}
				app.patternMatch(path, line, pattern)
			}
		}
		if err == io.EOF {
			if len(leftover) > 0 {
				if app.config.invertMatch {
					app.invertMatch(path, leftover, pattern)
					continue
				}
				app.patternMatch(path, leftover, pattern)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading file: %W", err)
		}
	}

	return nil
}

func (app *application) invertMatch(path string, line []byte, pattern *regexp.Regexp) {
	if !pattern.Match(line) {
		app.config.matchCount++
		if app.config.printPath {
			fmt.Printf("%s:", path)
		}
		fmt.Println(string(line))
	}
}

func (app *application) invertMatchStdin(line []byte, pattern *regexp.Regexp) {
	if !pattern.Match(line) {
		app.config.matchCount++
		fmt.Println(string(line))
	}
}

func (app *application) patternMatch(path string, line []byte, pattern *regexp.Regexp) {
	if pattern.Match(line) {
		app.config.matchCount++
		if app.config.printPath {
			fmt.Printf("%s:", path)
		}
		fmt.Println(string(line))
	}
}

func (app *application) patternMatchStdin(line []byte, pattern *regexp.Regexp) {
	if pattern.Match(line) {
		app.config.matchCount++
		fmt.Println(string(line))
	}
}
