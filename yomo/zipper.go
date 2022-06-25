package yomo

import (
	"errors"
	"io"
	"log"
	"net"
	"net/url"
	"sync"
)

func NewZipper(addr string) (Zipper, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "tcp" {
		return nil, errors.New("Currently only support TCP stream")
	}

	return &zipperTcpImpl{
		addr: u.Host,
		sfns: make(map[DataTag]string),
	}, nil
}

type zipperTcpImpl struct {
	listener net.Listener
	addr     string
	sfns     map[DataTag]string
	mu       sync.Mutex
}

func (z *zipperTcpImpl) Close() error {
	if z.listener != nil {
		z.listener.Close()
	}
	return nil
}

func (z *zipperTcpImpl) Serve() error {
	listener, err := net.Listen("tcp", z.addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	z.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		go z.processConn(conn)
	}
}

func (z *zipperTcpImpl) processConn(conn net.Conn) {
	defer conn.Close()

	t, err := readHandshakeType(conn)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	switch t {
	case TYPE_STREAM:
		h, err := readHandshake[handshakeStream](conn)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		if err = writeHandshakeResponse(conn, ""); err != nil {
			log.Printf("%v", err)
			return
		}

		if err = z.processStream(conn, h.Tag, h.Arg); err != nil {
			log.Printf("%v", err)
			return
		}
	case TYPE_SFN:
		h, err := readHandshake[handshakeSfn](conn)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		if err = writeHandshakeResponse(conn, ""); err != nil {
			log.Printf("%v", err)
			return
		}

		conn.Close()

		z.mu.Lock()
		z.sfns[h.Tag] = h.Addr
		z.mu.Unlock()
	default:
		log.Printf("Unsupported client type: %v", t)
		return
	}
}

func (z *zipperTcpImpl) processStream(stream io.ReadCloser, tag DataTag, arg []byte) error {
	z.mu.Lock()
	sfnAddr, ok := z.sfns[tag]
	z.mu.Unlock()
	if !ok {
		return errors.New("no observed sfn")
	}

	conn, err := net.Dial("tcp", sfnAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	h := &handshakeStream{Arg: arg}
	if err = writeHandshake(conn, TYPE_STREAM, h); err != nil {
		return err
	}

	if err = PipeStream(stream, conn); err != nil {
		return err
	}

	return nil
}
