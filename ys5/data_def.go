package ys5

import "yomo-mock/yomo"

const (
	TAG_CRAWLER = "crawler"
)

type ArgCrawler struct {
	Tid     string       `json:"tid"`
	Addr    string       `json:"addr"`
	SinkTag yomo.DataTag `json:"sink_tag"`
}

type ArgSink struct {
	Tid string `json:"tid"`
}
