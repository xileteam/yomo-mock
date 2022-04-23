package yomo

import (
	"bytes"
	"log"
	"net"
)

type SourceImpl struct{}

func NewSource() Source {
	return &SourceImpl{}
}

func (s *SourceImpl) Connect(name string, zipperAddr string) error {
	return nil
}

func (s *SourceImpl) NewStream(tag DataTag, arg string) (Stream, error) {
	conn, err := net.Dial("unix", "/tmp/yomo.sock")
	if err != nil {
		log.Fatalf("%v", err)
	}

	buf := bytes.NewBufferString(arg)
	for {
		if buf.Len() >= 32 {
			break
		}
		buf.WriteString(" ")
	}
	handshake := buf.Bytes()

	if _, err := conn.Write(handshake); err != nil {
		log.Fatalf("%v", err)
	}

	return conn, nil
}
