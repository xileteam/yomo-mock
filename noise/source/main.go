package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
	"yomo-mock/noise"
	"yomo-mock/yomo"
)

func main() {
	source, err := yomo.NewSource("tcp://localhost:9000")
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = source.Connect(); err != nil {
		log.Fatalf("%v", err)
	}
	defer source.Close()

	seed := rand.NewSource(time.Now().UnixNano())

	for {
		data := &noise.NoiseData{
			Time:  time.Now().Unix(),
			Noise: rand.New(seed).Float32() * 200,
		}

		buf, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if err = source.SendDatagram(noise.TAG_NOISE_1, buf); err != nil {
			log.Fatalf("%v", err)
		}

		log.Printf("[source] ✅ Emit %v", data)

		time.Sleep(1 * time.Second)
	}
}
