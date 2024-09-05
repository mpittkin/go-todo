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
	RepoTitle       string
}

const (
	outTypeConsole      = "console"
	outTypeSlackWebhook = "slack-webhook"
)

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
	cfg.RepoTitle = os.Getenv("GOTODO_REPO_TITLE")

	pathFlag := flag.String("root-path", "", "the root path from which directories will be traversed looking for files to parse")
	outputFlag := flag.String("output-type", "", "the output type ('console' or 'slack-webhook'")
	webhookFlag := flag.String("slack-webhook-url", "", "when output type is set to 'slack-webhook' defines the url to send the POST request")
	titleFlag := flag.String("repo-title", "", "when output type is set to slack-webhook, this is included in the report to indicate to the reader the source of the todos")

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
	if *titleFlag != "" {
		cfg.SlackWebhookURL = *titleFlag
	}

	absPath, err := filepath.Abs(cfg.RootPath)
	if err != nil {
		log.Fatalf("Unable to convert path '%s' to absolute path: %s", cfg.RootPath, err)
	}
	fmt.Println("Starting at " + absPath)

	var result []todo.Todo

	if err := os.Chdir(absPath); err != nil {
		panic(err)
	}

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

	switch cfg.OutputType {
	case outTypeConsole:
		output.ToConsole(result)
	case outTypeSlackWebhook:
		if err := output.ToSlackWebhook(result, cfg.SlackWebhookURL, cfg.RepoTitle); err != nil {
			log.Fatalf("error posting result to slack webhook: %s", err)
		}
		fmt.Println("Todo report sent successfully to Slack")
	default:
		log.Printf("invalid output type %s\n", cfg.OutputType)
	}
}
