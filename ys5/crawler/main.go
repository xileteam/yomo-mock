package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"yomo-mock/yomo"
	"yomo-mock/ys5"
)

func crawlerHandler(in io.ReadCloser, arg []byte) (yomo.DataTag, io.ReadCloser, []byte) {
	var argCrawler ys5.ArgCrawler
	if err := json.Unmarshal(arg, &argCrawler); err != nil {
		log.Printf("%v", err)
		return yomo.TAG_NIL, nil, nil
	}

	// 爬虫向外部建立连接
	conn, err := net.Dial("tcp", argCrawler.Addr)
	if err != nil {
		log.Printf("%v", err)
		return yomo.TAG_NIL, nil, nil
	}
	log.Printf("[%s] ++ %s", argCrawler.Tid, argCrawler.Addr)

	// 转发请求
	go yomo.PipeStream(in, conn)

	argSink := &ys5.ArgSink{Tid: argCrawler.Tid}
	if arg, err = json.Marshal(argSink); err != nil {
		log.Printf("%v", err)
		return yomo.TAG_NIL, nil, nil
	}

	return argCrawler.SinkTag, conn, arg
}

func main() {
	sfn, err := yomo.NewSFN(
		"tcp://localhost:9000",
		ys5.TAG_CRAWLER,
		crawlerHandler,
	)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = sfn.Connect(); err != nil {
		log.Fatalf("%v", err)
	}
	defer sfn.Close()

	if err = sfn.Serve(); err != nil {
		log.Fatalf("%v", err)
	}
}
