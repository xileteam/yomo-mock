package yomo

import (
	"io"
)

type DataTag string

const (
	TAG_NIL DataTag = ""
)

type DatagramHandler func(req []byte) (DataTag, []byte)

type StreamHandler func(in io.ReadCloser, arg []byte) (DataTag, io.ReadCloser, []byte)

type Source interface {
	io.Closer

	Connect() error

	SendDatagram(tag DataTag, data []byte) error

	NewStream(tag DataTag, arg []byte) (io.WriteCloser, error)
}

type DatagramSFN interface {
	io.Closer

	Connect(tag DataTag) error

	ServeDatagram(handler DatagramHandler) error
}

type StreamSFN interface {
	io.Closer

	Connect(tag DataTag) error

	ServeStream(handler StreamHandler) error
}

type Zipper interface {
	io.Closer

	Serve() error
}
