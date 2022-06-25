package yomo

import (
	"io"
)

type DataTag string

const (
	TAG_NIL DataTag = ""
)

type Handler func(in io.ReadCloser, arg []byte) (DataTag, io.ReadCloser, []byte)

type Source interface {
	io.Closer

	Connect() error

	NewStream(tag DataTag, arg []byte) (io.WriteCloser, error)
}

type SFN interface {
	io.Closer

	Connect() error

	Serve() error
}

type Zipper interface {
	io.Closer

	Serve() error
}
