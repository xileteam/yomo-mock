package yomo

import (
	"errors"
	"io"
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

	return &sourceTcpImpl{zipperAddr: u.Host}, nil
}

type sourceTcpImpl struct {
	zipperAddr   string
	datagramConn net.Conn
}

func (s *sourceTcpImpl) Close() error {
	if s.datagramConn != nil {
		s.datagramConn.Close()
	}
	return nil
}

func (s *sourceTcpImpl) SendDatagram(tag DataTag, data []byte) error {
	return WriteFrame(s.datagramConn, &DataFrame{
		Tag:  tag,
		Data: data,
	})
}

func (s *sourceTcpImpl) Connect() error {
	conn, err := net.Dial("tcp", s.zipperAddr)
	if err != nil {
		return err
	}

	h := &HandshakeFrame{ClientType: TYPE_SOURCE}
	if err = WriteFrame(conn, h); err != nil {
		conn.Close()
		return err
	}
	s.datagramConn = conn

	resp, err := ReadFrame[HandshakeResponseFrame](conn)
	if err != nil {
		conn.Close()
		return err
	} else if resp.Error != "" {
		conn.Close()
		return errors.New(resp.Error)
	}

	return nil
}

func (s *sourceTcpImpl) NewStream(tag DataTag, arg []byte) (io.WriteCloser, error) {
	conn, err := net.Dial("tcp", s.zipperAddr)
	if err != nil {
		return nil, err
	}

	h := &HandshakeFrame{
		ClientType: TYPE_STREAM,
		Tag:        tag,
		StreamArg:  arg,
	}
	if err = WriteFrame(conn, h); err != nil {
		conn.Close()
		return nil, err
	}

	resp, err := ReadFrame[HandshakeResponseFrame](conn)
	if err != nil {
		conn.Close()
		return nil, err
	} else if resp.Error != "" {
		conn.Close()
		return nil, errors.New(resp.Error)
	}

	return conn, nil
}
