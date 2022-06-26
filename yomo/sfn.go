package yomo

import (
	"errors"
	"log"
	"net"
	"net/url"
	"os"
)

func NewDatagramSFN(zipperAddr string) (DatagramSFN, error) {
	return newSFN(zipperAddr)
}

func NewStreamSFN(zipperAddr string) (StreamSFN, error) {
	return newSFN(zipperAddr)
}

func newSFN(zipperAddr string) (*sfnTcpImpl, error) {
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
	}, nil
}

type sfnTcpImpl struct {
	host         string
	port         string
	zipperAddr   string
	listener     net.Listener
	datagramConn net.Conn
}

func (s *sfnTcpImpl) Close() error {
	if s.datagramConn != nil {
		s.datagramConn.Close()
	}
	if s.listener != nil {
		s.listener.Close()
	}
	return nil
}

func (s *sfnTcpImpl) Connect(tag DataTag) error {
	conn, err := net.Dial("tcp", s.zipperAddr)
	if err != nil {
		return err
	}

	h := &HandshakeFrame{
		ClientType: TYPE_SFN,
		Tag:        tag,
		SFNAddr:    s.host + ":" + s.port,
	}
	if err = WriteFrame(conn, h); err != nil {
		conn.Close()
		return err
	}

	resp, err := ReadFrame[HandshakeResponseFrame](conn)
	if err != nil {
		conn.Close()
		return err
	} else if resp.Error != "" {
		conn.Close()
		return errors.New(resp.Error)
	}

	s.datagramConn = conn

	return nil
}

func (s *sfnTcpImpl) ServeDatagram(handler DatagramHandler) error {
	for {
		f, err := ReadFrame[DataFrame](s.datagramConn)
		if err != nil {
			return err
		}

		tag, resp := handler(f.Data)
		if tag != TAG_NIL && len(resp) > 0 {
			data := &DataFrame{
				Tag:  tag,
				Data: resp,
			}

			if err := WriteFrame(s.datagramConn, data); err != nil {
				return err
			}
		}
	}
}

func (s *sfnTcpImpl) ServeStream(handler StreamHandler) error {
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

		h, err := ReadFrame[HandshakeFrame](conn)
		if err != nil {
			return err
		}

		if err = WriteFrame(conn, &HandshakeResponseFrame{Error: ""}); err != nil {
			return err
		}

		go s.processStream(conn, h.StreamArg, handler)
	}
}

func (s *sfnTcpImpl) processStream(src net.Conn, arg []byte, handler StreamHandler) {
	defer src.Close()

	tag, stream, arg := handler(src, arg)
	if stream == nil {
		return
	}
	defer stream.Close()
	if tag == TAG_NIL {
		return
	}

	conn, err := net.Dial("tcp", s.zipperAddr)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer conn.Close()

	h := &HandshakeFrame{
		ClientType: TYPE_STREAM,
		Tag:        tag,
		StreamArg:  arg,
	}
	if err = WriteFrame(conn, h); err != nil {
		log.Printf("%v", err)
		return
	}

	resp, err := ReadFrame[HandshakeResponseFrame](conn)
	if err != nil {
		log.Printf("%v", err)
		return
	} else if resp.Error != "" {
		log.Printf("%s", resp.Error)
		return
	}

	if err = PipeStream(stream, conn); err != nil {
		log.Printf("%v", err)
		return
	}
}
