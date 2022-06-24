package main

import (
	"log"
	"os"

	"yomo-mock/yomo"
)

func main() {
	port := os.Getenv("YOMO_ZIPPER_PORT")
	if port == "" {
		port = "9000"
	}

	z, err := yomo.NewZipper("tcp://0.0.0.0:" + port)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = z.Serve(); err != nil {
		log.Fatalf("%v", err)
	}
}
