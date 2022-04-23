package yomo

import (
	"bytes"
	"log"
	"net"
	"net/url"
)

type SourceImpl struct {
	u *url.URL
}

func NewSource() Source {
	return &SourceImpl{}
}

func (s *SourceImpl) Connect(name string, zipperAddr string) error {
	u, err := url.Parse(zipperAddr)
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.u = u
	return nil
}

func (s *SourceImpl) NewStream(tag DataTag, arg string) (Stream, error) {
	conn, err := net.Dial(s.u.Scheme, s.u.Path)
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
