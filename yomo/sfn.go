package yomo

import (
	"context"
	"errors"
	"log"
	"net"
	"net/url"
	"os"
)

func NewSFN(zipperAddr string, tag DataTag, handler Handler) (SFN, error) {
	u, err := url.Parse(zipperAddr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "tcp" {
		return nil, errors.New("Currently only support TCP stream")
	}

	host := os.Getenv("YOMO_SFN_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("YOMO_SFN_PORT")
	if port == "" {
		port = "12000"
	}

	return &sfnTcpImpl{
		host:       host,
		port:       port,
		zipperAddr: u.Host,
		tag:        tag,
		handler:    handler,
	}, nil
}

type sfnTcpImpl struct {
	host       string
	port       string
	zipperAddr string
	tag        DataTag
	handler    Handler
}

func (s *sfnTcpImpl) Close() error {
	return nil
}

func (s *sfnTcpImpl) Connect() error {
	conn, err := net.Dial("tcp", s.zipperAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	h := &handshakeSfn{
		Addr: s.host + ":" + s.port,
		Tag:  s.tag,
	}

	if err = writeHandshake(conn, TYPE_SFN, h); err != nil {
		return err
	}

	if err = readHandshakeResponse(conn); err != nil {
		return err
	}

	return nil
}

func (s *sfnTcpImpl) Serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", "0.0.0.0:"+s.port)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("SFN Started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		if _, err = readHandshakeType(conn); err != nil {
			return err
		}

		h, err := readHandshake[handshakeStream](conn)
		if err != nil {
			return err
		}

		if err = writeHandshakeResponse(conn, ""); err != nil {
			return err
		}

		go s.process(conn, h.Arg)
	}
}

func (s *sfnTcpImpl) process(src net.Conn, arg []byte) {
	defer src.Close()

	tag, stream, arg := s.handler(src, arg)
	if stream == nil {
		return
	}
	defer stream.Close()

	conn, err := net.Dial("tcp", s.zipperAddr)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer conn.Close()

	h := &handshakeStream{
		Tag: tag,
		Arg: arg,
	}
	if err = writeHandshake(conn, TYPE_STREAM, h); err != nil {
		log.Printf("%v", err)
		return
	}

	if err = readHandshakeResponse(conn); err != nil {
		log.Printf("%v", err)
		return
	}

	if err = PipeStream(stream, conn); err != nil {
		log.Printf("%v", err)
		return
	}
}
