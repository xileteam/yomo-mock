package yomo

import (
	"io"
)

type DataTag int64

type Handler func(in io.ReadCloser, arg []byte) (DataTag, io.ReadCloser, []byte)

type Client interface {
	io.Closer

	Connect() error
}

type Source interface {
	Client

	NewStream(tag DataTag, arg []byte) (io.WriteCloser, error)
}

type Sfn interface {
	Client
}

type zipper interface {
	Serve() error
}
