package main

import (
	"io"
	"log"
	"net"
	"ys5-mock/utils"
	"ys5-mock/yomo"

	"golang.org/x/text/transform"
)

func CrawlerHandler(arg string, stream yomo.Stream) {
	// 爬虫向外部建立连接, arg为请求地址
	conn, err := net.Dial("tcp", arg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("Connected： ", arg)

	// 转发请求
	go utils.PipeStream(stream, conn)

	// 转发响应
	// 并将所有字符串转换为大写
	transformer := transform.NewReader(conn, &utils.Rot13Transformer{})
	go io.Copy(stream, transformer)

	// 原始转发响应 @wujunzhuo 的老代码
	// go utils.PipeStream(conn, stream)
}

func main() {
	sfn := yomo.NewSfn()
	sfn.WithObserveDataTags(0x0A).
		WithStreamHandler(0x0B, CrawlerHandler).
		Connect("SFN", "./yomo.sock")
}
