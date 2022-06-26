package yomo

import "io"

type Connection struct {
	writer io.Writer
	addr   string
}
