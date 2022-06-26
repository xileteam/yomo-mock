package main

import (
	"encoding/json"
	"log"
	"yomo-mock/noise"
	"yomo-mock/yomo"
)

func handler(req []byte) (yomo.DataTag, []byte) {
	var data noise.NoiseData
	if err := json.Unmarshal(req, &data); err != nil {
		log.Printf("%v", err)
		return yomo.TAG_NIL, nil
	}

	log.Printf("[sfn-2] ✅ Receive %v", data)

	return yomo.TAG_NIL, nil
}

func main() {
	sfn, err := yomo.NewDatagramSFN("tcp://localhost:9000")
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = sfn.Connect(noise.TAG_NOISE_2); err != nil {
		log.Fatalf("%v", err)
	}
	defer sfn.Close()

	if err = sfn.ServeDatagram(handler); err != nil {
		log.Fatalf("%v", err)
	}
}
