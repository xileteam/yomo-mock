package main

import (
	"log"
	"net"
	"ys5-mock/yomo"
)

func main() {
	source := yomo.NewSource()
	source.Connect("Source", "localhost:9000")

	// 监听socks5端口
	server, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for {
		// 新请求到达
		conn, err := server.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}

		// socks5认证、解析
		if err := Auth(conn); err != nil {
			log.Fatalf("%v", err)
		}
		addr, err := Request(conn)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// 向zipper创建新流
		stream, err := source.NewStream(0x0A, addr)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// 转发请求
		go yomo.PipeStream(conn, stream)

		// 转发响应
		go yomo.PipeStream(stream, conn)
	}
}
