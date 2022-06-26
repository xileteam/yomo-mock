package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
	"yomo-mock/yomo"
	"yomo-mock/ys5"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var conns sync.Map

func sinkHandler(in io.ReadCloser, arg []byte) (yomo.DataTag, io.ReadCloser, []byte) {
	var argSink ys5.ArgSink

	if err := json.Unmarshal(arg, &argSink); err != nil {
		log.Printf("%v", err)
		return yomo.TAG_NIL, nil, nil
	}

	log.Printf("[%s] ++", argSink.Tid)

	if v, ok := conns.Load(argSink.Tid); ok {
		conn := v.(net.Conn)

		// 转发响应
		yomo.PipeStream(in, conn)

		conn.Close()
		conns.Delete(argSink.Tid)
	}

	return yomo.TAG_NIL, nil, nil
}

func main() {
	// 创建source
	source, err := yomo.NewSource("tcp://localhost:9000")
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = source.Connect(); err != nil {
		log.Fatalf("%v", err)
	}
	defer source.Close()

	gid, err := gonanoid.New()
	if err != nil {
		log.Fatalf("%v", err)
	}
	sinkTag := yomo.DataTag(gid)

	sink, err := yomo.NewStreamSFN("tcp://localhost:9000")
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = sink.Connect(sinkTag); err != nil {
		log.Fatalf("%v", err)
	}
	defer sink.Close()

	go sink.ServeStream(sinkHandler)

	// 监听socks5端口
	server, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatalf("server listen %v", err)
	}

	for {
		// 新请求到达
		conn, err := server.Accept()
		if err != nil {
			log.Printf("%v", err)
			break
		}

		// socks5认证
		if err := ys5.Auth(conn); err != nil {
			log.Printf("%v", err)
			conn.Close()
			continue
		}

		// socks5解析
		addr, err := ys5.Request(conn)
		if err != nil {
			log.Printf("%v", err)
			conn.Close()
			continue
		}

		tid, err := gonanoid.New()
		if err != nil {
			log.Printf("%v", err)
			conn.Close()
			continue
		}

		arg := &ys5.ArgCrawler{
			Tid:     tid,
			Addr:    addr,
			SinkTag: sinkTag,
		}

		buf, err := json.Marshal(arg)
		if err != nil {
			log.Printf("%v", err)
			conn.Close()
			continue
		}

		// 向zipper创建新流
		stream, err := source.NewStream(ys5.TAG_CRAWLER, buf)
		if err != nil {
			log.Printf("%v", err)
			conn.Close()
			continue
		}

		conns.Store(tid, conn)

		// 转发请求
		go yomo.PipeStream(conn, stream)
	}
}
