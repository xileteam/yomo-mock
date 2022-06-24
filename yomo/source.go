package yomo

import (
	"errors"
	"io"
	"log"
	"net"
	"net/url"
)

func NewSource(zipperAddr string) (Source, error) {
	u, err := url.Parse(zipperAddr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "tcp" {
		return nil, errors.New("Currently only support TCP stream")
	}

	return &SourceTcpImpl{zipperAddr: u.Host}, nil
}

type SourceTcpImpl struct {
	zipperAddr string
}

func (s *SourceTcpImpl) Close() error {
	return nil
}

func (s *SourceTcpImpl) Connect() error {
	log.Println("Source Started")
	return nil
}

func (s *SourceTcpImpl) NewStream(tag DataTag, arg []byte) (io.WriteCloser, error) {
	conn, err := net.Dial("tcp", s.zipperAddr)
	if err != nil {
		return nil, err
	}

	h := &handshakeStream{
		Tag: tag,
		Arg: arg,
	}
	if err = writeHandshake(conn, TYPE_STREAM, h); err != nil {
		conn.Close()
		return nil, err
	}

	if err = readHandshakeResponse(conn); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}
