package bpm

import "time"

const minBeatInterval = 100 * time.Millisecond

type BPM struct {
	Bpm        float64
	LastEnergy float64
	LastBeat   time.Time
	HasBeat    bool
}

func (bpm *BPM) UpdateBPM(samples []uint8) {
	// convert from u8 to f64
	samplesShifted := make([]float64, len(samples))
	shift := 256.0 / float64(2)
	for i := 0; i < len(samples); i++ {
		samplesShifted[i] = float64(samples[i]) - shift
	}

	energy := 0.0
	for _, s := range samplesShifted {
		energy += s * s
	}
	energy = energy / float64(len(samples))

	// smooth out with last reduce big jumps
	alpha := 0.4
	energy = alpha*energy + (1-alpha)*(bpm.LastEnergy)

	ratio := energy / (bpm.LastEnergy + 1e-8)
	bpm.LastEnergy = energy

	// if the ratio has not been increased by 50%
	if ratio <= 1.5 {
		bpm.HasBeat = false
		return
	}

	now := time.Now()
	lastBeatTime := now.Sub(bpm.LastBeat)

	// if the last sudden spike in energy is not t least minBeatInterval old
	if lastBeatTime < minBeatInterval {
		return
	}

	bpmEstimate := 60 / lastBeatTime.Seconds()
	// again smooth out to reduce sudden spikes
	bpm.Bpm = 0.4*bpm.Bpm + 0.6*bpmEstimate

	bpm.LastBeat = now
	bpm.HasBeat = true
}
