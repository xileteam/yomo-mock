package yomo

import (
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
	listener   net.Listener
}

func (s *sfnTcpImpl) Close() error {
	if s.listener != nil {
		s.listener.Close()
	}
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

func (s *sfnTcpImpl) Serve() error {
	listener, err := net.Listen("tcp", "0.0.0.0:"+s.port)
	if err != nil {
		return err
	}
	defer listener.Close()
	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
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
	defer stream.Close()

	if tag == TAG_NIL || stream == nil {
		return
	}

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
