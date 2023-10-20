package main

import (
	"flag"
	"fmt"
	"github.com/mpittkin/go-todo/output"
	"github.com/mpittkin/go-todo/todo"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	RootPath        string
	OutputType      string
	SlackWebhookURL string
}

const (
	defaultRootPath        = "."
	defaultOutputType      = "console"
	defaultSlackWebhookURL = ""
)

func main() {
	cfg := config{
		RootPath:        defaultRootPath,
		OutputType:      defaultOutputType,
		SlackWebhookURL: defaultSlackWebhookURL,
	}
	if r := os.Getenv("GOTODO_ROOT_PATH"); r != "" {
		cfg.RootPath = r
	}
	if o := os.Getenv("GOTODO_OUTPUT_TYPE"); o != "" {
		cfg.OutputType = o
	}
	if u := os.Getenv("GOTODO_SLACK_URL"); u != "" {
		cfg.SlackWebhookURL = u
	}
	pathFlag := flag.String("root-path", "", "the root path from which directories will be traversed looking for files to parse")
	outputFlag := flag.String("output-type", "", "the output type (console, json, or slack-webhook")
	webhookFlag := flag.String("slack-webhook-url", "", "when output type is set to 'slack-webhook' defines the url to send the POST request")
	flag.Parse()
	if *pathFlag != "" {
		cfg.RootPath = *pathFlag
	}
	if *outputFlag != "" {
		cfg.OutputType = *outputFlag
	}
	if *webhookFlag != "" {
		cfg.SlackWebhookURL = *webhookFlag
	}

	// Output type (defaults to console)
	// Slack webhook URL (only used for that output type)

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
