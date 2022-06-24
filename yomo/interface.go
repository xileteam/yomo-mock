package yomo

import (
	"context"
	"io"
)

type DataTag int64

type Handler func(in io.ReadCloser, arg []byte) (DataTag, io.ReadCloser, []byte)

type Source interface {
	io.Closer

	Connect() error

	NewStream(tag DataTag, arg []byte) (io.WriteCloser, error)
}

type SFN interface {
	io.Closer

	Connect() error

	Serve(ctx context.Context) error
}

type Zipper interface {
	Serve(ctx context.Context) error
}
