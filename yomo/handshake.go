package yomo

import (
	"encoding/json"
	"errors"
	"io"
)

const (
	TYPE_STREAM = '1'
	TYPE_SFN    = '2'
)

type handshakeSfn struct {
	Addr string  `json:"addr"`
	Tag  DataTag `json:"tag"`
}

type handshakeStream struct {
	Tag DataTag `json:"tag"`
	Arg []byte  `json:"arg"`
}

func readHandshakeType(reader io.Reader) (byte, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func readHandshake[T any](reader io.Reader) (*T, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	length := bytesToInt(buf)

	buf = make([]byte, length)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}

	var h T
	if err := json.Unmarshal(buf, &h); err != nil {
		return nil, err
	}

	return &h, nil
}

func readHandshakeResponse(reader io.Reader) error {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return err
	}
	length := bytesToInt(buf)

	if length > 0 {
		buf = make([]byte, length)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return err
		}

		return errors.New(string(buf))
	}

	return nil
}

func writeHandshake[T any](writer io.Writer, handshakeType byte, h *T) error {
	if _, err := writer.Write([]byte{handshakeType}); err != nil {
		return err
	}

	data, err := json.Marshal(h)
	if err != nil {
		return err
	}

	lengthBuf := intToBytes(len(data))
	if _, err := writer.Write(lengthBuf); err != nil {
		return err
	}

	if _, err = writer.Write(data); err != nil {
		return err
	}

	return nil
}

func writeHandshakeResponse(writer io.Writer, resp string) error {
	lengthBuf := intToBytes(len(resp))
	if _, err := writer.Write(lengthBuf); err != nil {
		return err
	}

	if _, err := writer.Write([]byte(resp)); err != nil {
		return err
	}

	return nil
}
