package fft

/*
#include "chuck_fft.h"
#include <stdlib.h>
*/
import "C"
import (
	"math"
	"unsafe"
)

const BUFFER_SIZE = 256
const BINS_SIZE = 256 / 4

type FilterType int

const (
	Block FilterType = iota + 1
	BoxFilter
	DoubleBoxFilter
)

type FFT struct {
	Bins []float64

	AmpScalar   int
	FilterMode  FilterType
	FilterRange int
	Alpha       float64
	Decay       int
}

func InitFFT(
	AmpScalar int,
	FilterMode FilterType,
	FilterRange int,
	Alpha float64,
	Decay int,
) *FFT {
	return &FFT{
		Bins:        make([]float64, BINS_SIZE),
		AmpScalar:   AmpScalar,
		FilterRange: FilterRange,
		FilterMode:  FilterMode,
		Alpha:       Alpha,
		Decay:       Decay,
	}
}

func (f *FFT) applyBlockFilter() {
	if f.FilterRange == 0 {
		return
	}

	for i := 0; i < BINS_SIZE; i += f.FilterRange {
		sum := 0.0
		for j := i; j < min(i+f.FilterRange, BINS_SIZE); j++ {
			sum += f.Bins[i]
		}

		avg := sum / float64(f.FilterRange)

		for j := i; j < min(i+f.FilterRange, BINS_SIZE); j++ {
			f.Bins[j] = avg
		}
	}
}

func (f *FFT) applyBoxFilter() {
	var tmp = make([]float64, BINS_SIZE)

	for i := 0; i < BINS_SIZE; i++ {
		sum := 0.

		start := max(i-f.FilterRange, 0)
		end := min(i+f.FilterRange, BINS_SIZE-1)
		for j := start; j <= end; j++ {
			sum += f.Bins[j]
		}

		avg := sum / float64((end-start)+1)

		tmp[i] = avg
	}

	f.Bins = tmp
}

func (f *FFT) filterFFT() {
	switch f.FilterMode {
	case Block:
		f.applyBlockFilter()
		break
	case BoxFilter:
		f.applyBoxFilter()
		break
	case DoubleBoxFilter:
		f.applyBoxFilter()
		f.applyBoxFilter()
		break
	}
}

func (f *FFT) ApplyFFT(samples []uint8) {
	var fft_tmp [BUFFER_SIZE]float64

	// shift by half of 256 to center the u8 on 0
	shift := 256.0 / float64(2)
	for i := 0; i < BUFFER_SIZE; i++ {
		fft_tmp[i] = (float64(samples[i]) - shift) * (float64(f.AmpScalar) / shift)
	}

	c_fft_tmp := convToCFloatArray(fft_tmp)

	C.rfft(c_fft_tmp, BUFFER_SIZE/2, 1)

	fft_tmp = convFromCFloatArray(c_fft_tmp)

	fft_tmp[0] = fft_tmp[2]

	// normalize
	for i := 0; i < BUFFER_SIZE; i++ {
		fft_tmp[i] *= 0.04 + (0.5 * (float64(i) / float64(BUFFER_SIZE)))
	}

	for i := 0; i < BUFFER_SIZE/2; i += 2 {
		var mag float64 = math.Sqrt(fft_tmp[i]*fft_tmp[i] + fft_tmp[i+1]*fft_tmp[i+1])

		mag = (0.7 * math.Log10(1.1*mag)) + (0.7 * mag)

		mag = clamp(mag, 0.0, 1.0)

		prevmag := f.Bins[i/2]

		if mag > prevmag {
			f.Bins[i/2] = mag
		}

		if mag < prevmag {
			f.Bins[i/2] = prevmag * (float64(f.Decay) / 100.0)

			if f.Bins[i/2] < 0.0001 {
				f.Bins[i/2] = 0
			}
		}
	}

	f.filterFFT()
}

func convToCFloatArray(a [BUFFER_SIZE]float64) *C.float {
	cslice := make([]C.float, len(a))
	for i, v := range a {
		cslice[i] = C.float(v)
	}
	return (*C.float)(unsafe.Pointer(&cslice[0]))
}

func convFromCFloatArray(ca *C.float) [BUFFER_SIZE]float64 {
	var a [BUFFER_SIZE]float64

	cSlice := (*[BUFFER_SIZE]C.float)(unsafe.Pointer(ca))[:]

	for i := 0; i < BUFFER_SIZE; i++ {
		a[i] = float64(cSlice[i])
	}

	return a
}

func clamp[T float64 | int](v, _min, _max T) T {
	return max(_min, min(_max, v))
}
