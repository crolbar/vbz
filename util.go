package main

import (
	"bytes"
	"encoding/binary"
	"time"
)

func byteToU8(data []byte) ([]uint8, error) {
	buf := bytes.NewReader(data)

	var result []uint8

	for {
		var u uint8
		err := binary.Read(buf, binary.LittleEndian, &u)
		if err != nil {
			break
		}
		result = append(result, uint8(u))
	}

	return result, nil
}

func (v *VBZ) getPeakLowAmp() float64 {
	peakLow := 0.0
	for i, mag := range v.fft.Bins {
		if i > 5 {
			break
		}
		if mag > float64(peakLow) {
			peakLow = mag
		}
	}
	return peakLow
}

const minBeatInterval = 100 * time.Millisecond

func (v *VBZ) getBPM(samples []uint8) {
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
	energy = alpha*energy + (1-alpha)*(v.bpm.lastEnergy)

	ratio := energy / (v.bpm.lastEnergy + 1e-8)
	v.bpm.lastEnergy = energy

	// if the ratio has not been increased by 50%
	if ratio <= 1.5 {
		v.bpm.hasBeat = false
		return
	}

	now := time.Now()
	lastBeatTime := now.Sub(v.bpm.lastBeat)

	// if the last sudden spike in energy is not t least minBeatInterval old
	if lastBeatTime < minBeatInterval {
		return
	}

	bpmEstimate := 60 / lastBeatTime.Seconds()
	// again smooth out to reduce sudden spikes
	v.bpm.bpm = 0.4*v.bpm.bpm + 0.6*bpmEstimate

	v.bpm.lastBeat = now
	v.bpm.hasBeat = true
}
