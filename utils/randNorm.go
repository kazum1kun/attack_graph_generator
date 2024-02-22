package utils

import (
	"math"
	"math/rand"
)

var mean, std, minLimit float64
var ready = false

func SetGaussianParams(_mean, _std, _minLimit float64) {
	mean = _mean
	std = _std
	minLimit = _minLimit
	ready = true
}

func RandNorm(r *rand.Rand) float64 {
	if !ready {
		println("RNG not ready, please set up parameters first!")
		return -1
	} else {
		temp := r.NormFloat64()*std + mean
		return math.Min(1, math.Max(minLimit, temp))
	}
}
