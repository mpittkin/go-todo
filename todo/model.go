package todo

import "time"

type Todo struct {
	Path string
	Line int
	Text string
	BlameInfo
}

type BlameInfo struct {
	Name string
	Mail string
	Time time.Time
}
