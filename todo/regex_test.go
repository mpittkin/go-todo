package todo

import (
	"regexp"
	"testing"
)

func TestMatch(t *testing.T) {
	var isTodo = regexp.MustCompile(todoLineRegex)

	tests := []struct {
		str       string
		wantMatch bool
	}{
		{"todo: do stuff", true},
		{"TODO: do stuff", true},
		{"ToDo - do stuff", true},
		{"we should todo later", true},
		{"look at the rest of our todos", false},
		{"regular old normal comment line", false},
	}

	for _, tt := range tests {
		gotMatch := isTodo.MatchString(tt.str)
		if gotMatch != tt.wantMatch {
			t.Errorf("\"%s\" wanted match %t got %t", tt.str, tt.wantMatch, gotMatch)
		}
	}
}
