package hues

import (
	"math"
	"vbz/fft"
)

type Hues struct {
	PrevFHues []float64
	PrevBHues []float64
}

func InitHues() *Hues {
	var h Hues
	h.PrevFHues = make([]float64, fft.BINS_SIZE)
	h.PrevBHues = make([]float64, fft.BINS_SIZE)
	for i := 0; i < fft.BINS_SIZE; i++ {
		h.PrevFHues[i] = float64(i) / fft.BINS_SIZE
		h.PrevBHues[i] = float64(i+3) / fft.BINS_SIZE
	}
	return &h
}

func (h *Hues) UpdateHues(hueRate float64, bpm float64) {
	// TODO ? bpm scaling the hueRate is a bit distracting.
	// maybe add a option to enable/disable the bpm scaling
	hueRateP := math.Pow(hueRate+(bpm*1e-4), 0.99)

	for i := 0; i < fft.BINS_SIZE; i++ {
		h.PrevBHues[i] = math.Mod(1+h.PrevBHues[i]-hueRate, 1)
		h.PrevFHues[i] = math.Mod(1+h.PrevFHues[i]-hueRateP*0.7, 1)
	}
}
