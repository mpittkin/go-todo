package main

import (
	"fmt"
	"github.com/mpittkin/go-todo/output"
	"github.com/mpittkin/go-todo/todo"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {

	// Set root path to start parsing
	rootPath := "."
	if len(os.Args) > 1 {
		rootPath = os.Args[1]
		if err := os.Chdir(rootPath); err != nil {
			log.Fatalf("Unable to change directory to %s\n", rootPath)
		}
	}

	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		log.Fatalf("Unable to convert path '%s' to absolute path: %s", rootPath, err)
	}
	fmt.Println("Starting at " + absPath)

	var result []todo.Todo

	// Walk through the files and process all .go files
	err = fs.WalkDir(os.DirFS(absPath), ".", func(path string, d fs.DirEntry, err error) error {
		// Skip walking .git because it contains so many small files it slows down the program substantially
		if d.IsDir() {
			if d.Name() == ".git" {
				return fs.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".go" {
			return nil
		}

		fmt.Println("Found go file: " + path)

		todos, err := todo.ParseGo(path)
		if err != nil {
			log.Fatalf("parse %s: %s", path, err)
		}

		result = append(result, todos...)

		return nil
	})
	if err != nil {
		log.Fatalf("walk path %s: %s", absPath, err)
	}
	output.ToConsole(result)
}
