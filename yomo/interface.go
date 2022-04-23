package yomo

import (
	"io"
)

type DataTag int64

type Stream io.ReadWriteCloser

type StreamHandler func(arg string, stream Stream)

type Client interface {
	Connect(name string, zipperAddr string) error
}

type Source interface {
	Client

	NewStream(tag DataTag, arg string) (Stream, error)
}

type Sfn interface {
	Client

	WithObserveDataTags(tags ...DataTag) Sfn
	WithStreamHandler(tag DataTag, handler StreamHandler) Sfn
}
