package yomo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

func intToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func bytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func PipeStream(src io.ReadCloser, dst io.WriteCloser) error {
	defer src.Close()
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			return err
		}
	}

	return nil
}
