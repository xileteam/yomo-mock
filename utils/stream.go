package utils

import (
	"errors"
	"io"
	"log"
	"net"
	"ys5-mock/yomo"
)

func PipeStream(src yomo.Stream, dst yomo.Stream) {
	defer src.Close()
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			log.Fatalf("%v", err)
		}
	}
}
