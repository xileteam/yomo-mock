package main

import (
	"log"
	"net"
	"ys5-mock/utils"
	"ys5-mock/yomo"
)

func CrawlerHandler(arg string, stream yomo.Stream) {
	// 爬虫向外部建立连接, arg为请求地址
	conn, err := net.Dial("tcp", arg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// 转发请求
	go utils.PipeStream(stream, conn)

	// 转发响应
	go utils.PipeStream(conn, stream)
}

func main() {
	sfn := yomo.NewSfn()
	sfn.WithObserveDataTags(0x0A).
		WithStreamHandler(0x0B, CrawlerHandler).
		Connect("SFN", "unix:///tmp/yomo.sock")
}
