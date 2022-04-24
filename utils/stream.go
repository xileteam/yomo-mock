package utils

import (
	"errors"
	"io"
	"log"
	"net"
	"ys5-mock/yomo"

	"golang.org/x/text/transform"
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

type Rot13Transformer struct{ transform.NopResetter }

func (Rot13Transformer) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	log.Println("Rot13 before", string(src))
	for i := 0; i < len(src); i++ {
		if src[i] >= 'a' && src[i] <= 'z' {
			dst[i] = ((src[i] - 'a' + 13) % 26) + 'a'
		} else if src[i] >= 'A' && src[i] <= 'Z' {
			dst[i] = ((src[i] - 'A' + 13) % 26) + 'A'
		} else {
			dst[i] = src[i]
		}
	}
	log.Println("Rot13 after", string(dst))
	return len(src), len(src), nil
}
