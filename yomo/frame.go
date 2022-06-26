package yomo

import (
	"encoding/json"
	"io"
)

const (
	TYPE_SOURCE = '0'
	TYPE_STREAM = '1'
	TYPE_SFN    = '2'
)

type HandshakeFrame struct {
	ClientType byte    `json:"client_type"`
	Tag        DataTag `json:"tag"`
	StreamArg  []byte  `json:"stream_arg"`
	SFNAddr    string  `json:"sfn_addr"`
}

type HandshakeResponseFrame struct {
	Error string `json:"error"`
}

type DataFrame struct {
	Tag  DataTag `json:"tag"`
	Data []byte  `json:"data"`
}

func ReadFrame[T any](reader io.Reader) (*T, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	length := bytesToInt(buf)

	buf = make([]byte, length)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}

	var data T
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func WriteFrame[T any](writer io.Writer, data T) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	lengthBuf := intToBytes(len(buf))
	if _, err := writer.Write(lengthBuf); err != nil {
		return err
	}

	if _, err := writer.Write(buf); err != nil {
		return err
	}

	return nil
}
