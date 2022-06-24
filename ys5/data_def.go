package ys5

const (
	DATATAG_CRAWLER = 0x0A
	DATATAG_SINK    = 0x0B
)

type ArgCrawler struct {
	Tid  string `json:"tid"`
	Addr string `json:"addr"`
}

type ArgSink struct {
	Tid string `json:"tid"`
}
