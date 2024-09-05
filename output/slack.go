package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mpittkin/go-todo/todo"
	"io"
	"log"
	"net/http"
	"sort"
)

type AuthorTodos struct {
	AuthorMail string
	Todos      []todo.Todo
}

func ToSlackWebhook(todos []todo.Todo, webhookUrl string, repo string) error {
	byAuthor := make(map[string][]todo.Todo)
	for _, td := range todos {
		authorTodos := byAuthor[td.Mail]
		authorTodos = append(authorTodos, td)
		byAuthor[td.Mail] = authorTodos
	}

	var byAuthorSl []AuthorTodos
	for authorMail, authorTodos := range byAuthor {
		byAuthorSl = append(byAuthorSl, AuthorTodos{
			AuthorMail: authorMail,
			Todos:      authorTodos,
		})
	}

	sort.Slice(byAuthorSl, func(i, j int) bool {
		return byAuthorSl[i].AuthorMail < byAuthorSl[j].AuthorMail
	})

	message := fmt.Sprintf(`
Todo Report: %s
Total Todos: %d
`, repo, len(todos))

	for _, auth := range byAuthorSl {
		authorBlock := fmt.Sprintf(`
*%s*
`, auth.AuthorMail)
		for _, td := range auth.Todos {

			authorBlock += fmt.Sprintf("%s:%d (%v) `%s`\n", td.Path, td.Line, td.Time.Format("2006-01-02"), td.Text)
		}

		message += authorBlock
	}

	slackMsgBody := struct {
		Text string `json:"text"`
	}{
		Text: message,
	}

	body, err := json.Marshal(slackMsgBody)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	if err := PostToWebhook(webhookUrl, bytes.NewReader(body)); err != nil {
		return fmt.Errorf("post to slack webhook %s: %w", webhookUrl, err)
	}

	return nil
}

func PostToWebhook(url string, body io.Reader) error {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("close response body: %s", err)
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic("failed to read response body: " + err.Error())
		}
		fmt.Println(string(errBody))
		return fmt.Errorf("http response error: %s", resp.Status)
	}
	return nil
}
