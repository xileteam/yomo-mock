package yomo

import (
	"errors"
	"io"
	"log"
	"net"
)

func PipeStream(src Stream, dst Stream) {
	defer src.Close()
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			log.Fatalf("%v", err)
		}
	}
}
