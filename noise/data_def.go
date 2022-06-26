package noise

type NoiseData struct {
	Time  int64   `json:"time"`
	Noise float32 `json:"noise"`
}

const (
	TAG_NOISE_1 = "noise-1"
	TAG_NOISE_2 = "noise-2"
)
