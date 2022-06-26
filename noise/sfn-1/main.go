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

	log.Printf("[sfn-1] ✅ Receive %v", data)

	data.Noise *= 10
	buf, err := json.Marshal(data)
	if err != nil {
		log.Printf("%v", err)
		return yomo.TAG_NIL, nil
	}

	return noise.TAG_NOISE_2, buf
}

func main() {
	sfn, err := yomo.NewDatagramSFN("tcp://localhost:9000")
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = sfn.Connect(noise.TAG_NOISE_1); err != nil {
		log.Fatalf("%v", err)
	}
	defer sfn.Close()

	if err = sfn.ServeDatagram(handler); err != nil {
		log.Fatalf("%v", err)
	}
}
