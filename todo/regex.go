package todo

import "regexp"

//Match any line that contains the word t-o-d-o, case-insensitive

const todoLineRegex = "(?i)\\btodo\\b"

var todoLineMatcher = regexp.MustCompile(todoLineRegex)
