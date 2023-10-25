package output

import (
	"fmt"
	"github.com/mpittkin/go-todo/todo"
)

func ToConsole(todos []todo.Todo) {
	for _, todo := range todos {
		fmt.Println(todo)
	}

	fmt.Printf("Found %d todos\n", len(todos))
}
