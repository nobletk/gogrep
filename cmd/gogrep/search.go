package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (app *application) ProcessStdin(pattern string) error {
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
					app.invertMatchStdin(line, []byte(pattern))
					continue
				}
				if pattern == "" || bytes.Contains(line, []byte(pattern)) {
					fmt.Println(string(line))
				}
			}
		}
		if err == io.EOF {
			if len(leftover) > 0 && app.config.invertMatch {
				app.invertMatchStdin(leftover, []byte(pattern))
				continue
			}
			if len(leftover) > 0 && (pattern == "" || bytes.Contains(leftover, []byte(pattern))) {
				fmt.Println(string(leftover))
			}
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading stdin: %W", err)
		}
	}
	return nil
}

func (app *application) ProcessPaths(paths []string, pattern string) error {
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			continue
		}

		if info.IsDir() && !app.config.recurse {
			fmt.Printf("gogrep: %s: Is a directory\n", path)
			continue
		}
		if info.IsDir() {
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

func (app *application) processFile(path, pattern string) error {
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
					app.invertMatch(path, line, []byte(pattern))
					continue
				}
				if pattern == "" || bytes.Contains(line, []byte(pattern)) {
					if app.config.printPath {
						fmt.Printf("%s:", path)
					}
					fmt.Println(string(line))
				}
			}
		}
		if err == io.EOF {
			if len(leftover) > 0 && app.config.invertMatch {
				app.invertMatch(path, leftover, []byte(pattern))
				continue
			}
			if len(leftover) > 0 && (pattern == "" || bytes.Contains(leftover, []byte(pattern))) {
				fmt.Println(string(leftover))
			}
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading file: %W", err)
		}
	}
	return nil
}

func (app *application) invertMatch(path string, line, pattern []byte) {
	if !bytes.Contains(line, pattern) {
		if app.config.printPath {
			fmt.Printf("%s:", path)
		}
		fmt.Println(string(line))
	}
}

func (app *application) invertMatchStdin(line, pattern []byte) {
	if !bytes.Contains(line, pattern) {
		fmt.Println(string(line))
	}
}
