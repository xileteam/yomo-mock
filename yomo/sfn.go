package yomo

import (
	"log"
	"net"
	"net/url"
	"strings"
)

type SfnImpl struct {
	observeTag DataTag
	handler    StreamHandler
}

func NewSfn() Sfn {
	return &SfnImpl{}
}

func (s *SfnImpl) Connect(name string, zipperAddr string) error {
	u, err := url.Parse(zipperAddr)
	if err != nil {
		log.Fatalf("%v", err)
	}

	listener, err := net.Listen(u.Scheme, u.Path)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}

		go func() {
			handshake := make([]byte, 32)
			if _, err := conn.Read(handshake); err != nil {
				log.Fatalf("%v", err)
			}

			arg := strings.TrimSpace(string(handshake))
			s.handler(arg, conn)
		}()
	}
}

func (s *SfnImpl) WithObserveDataTags(tags ...DataTag) Sfn {
	s.observeTag = tags[0]
	return s
}

func (s *SfnImpl) WithStreamHandler(tag DataTag, handler StreamHandler) Sfn {
	s.handler = handler
	return s
}
