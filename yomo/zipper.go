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
		sfns: make(map[DataTag]*Connection),
	}, nil
}

type zipperTcpImpl struct {
	listener net.Listener
	addr     string
	sfns     map[DataTag]*Connection
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

	h, err := ReadFrame[HandshakeFrame](conn)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if err = WriteFrame(conn, &HandshakeResponseFrame{Error: ""}); err != nil {
		log.Printf("%v", err)
		return
	}

	switch h.ClientType {
	case TYPE_SOURCE:
		if err = z.processDatagram(conn); err != nil {
			log.Printf("%v", err)
		}
	case TYPE_STREAM:
		if err = z.processStream(conn, h.Tag, h.StreamArg); err != nil {
			log.Printf("%v", err)
		}
	case TYPE_SFN:
		z.mu.Lock()
		z.sfns[h.Tag] = &Connection{
			writer: conn,
			addr:   h.SFNAddr,
		}
		z.mu.Unlock()

		if err = z.processDatagram(conn); err != nil {
			log.Printf("%v", err)
		}
	default:
		log.Printf("Unsupported handshake client type: %v", h.ClientType)
	}
}

func (z *zipperTcpImpl) processDatagram(reader io.Reader) error {
	for {
		f, err := ReadFrame[DataFrame](reader)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		z.mu.Lock()
		sfnConn, ok := z.sfns[f.Tag]
		z.mu.Unlock()

		if ok {
			if err = WriteFrame(sfnConn.writer, f); err != nil {
				if errors.Is(err, net.ErrClosed) {
					z.mu.Lock()
					delete(z.sfns, f.Tag)
					z.mu.Unlock()
					return nil
				}
				return err
			}
		}
	}
}

func (z *zipperTcpImpl) processStream(stream io.ReadCloser, tag DataTag, arg []byte) error {
	z.mu.Lock()
	sfnConn, ok := z.sfns[tag]
	z.mu.Unlock()
	if !ok {
		return errors.New("no observed sfn")
	}

	conn, err := net.Dial("tcp", sfnConn.addr)
	if err != nil {
		z.mu.Lock()
		delete(z.sfns, tag)
		z.mu.Unlock()
		return err
	}

	h := &HandshakeFrame{
		ClientType: TYPE_STREAM,
		StreamArg:  arg,
	}
	if err = WriteFrame(conn, h); err != nil {
		return err
	}

	if err = PipeStream(stream, conn); err != nil {
		return err
	}

	return nil
}
