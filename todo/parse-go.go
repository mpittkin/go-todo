package todo

import (
	"go/parser"
	"go/token"
	"log"
)

func ParseGo(path string) ([]Todo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var todos []Todo
	for _, c := range f.Comments {
		for _, sc := range c.List {
			isTodo := todoLineMatcher.MatchString(sc.Text)
			if err != nil {
				log.Fatal(err)
			}
			if isTodo {
				pos := fset.Position(sc.Pos())
				todos = append(todos, Todo{
					Path: path,
					Line: pos.Line,
					Text: sc.Text,
				})
			}
		}
	}

	for i := 0; i < len(todos); i++ {
		blame, err := Blame(path, todos[i].Line)
		if err != nil {
			return nil, err
		}
		todos[i].BlameInfo = blame
	}

	return todos, nil
}
