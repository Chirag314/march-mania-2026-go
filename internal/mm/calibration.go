package mm

import "math"

func ClipProb(p, lo, hi float64) float64 {
	if p < lo {
		return lo
	}
	if p > hi {
		return hi
	}
	return p
}

// TemperatureScale applies logit(p)/t then sigmoid.
// t > 1 makes probabilities LESS extreme (usually better for Brier).
func TemperatureScale(p, t float64) float64 {
	if t <= 0 {
		return p
	}
	const eps = 1e-15
	if p < eps {
		p = eps
	}
	if p > 1-eps {
		p = 1 - eps
	}
	logit := math.Log(p / (1 - p))
	logit /= t
	return 1 / (1 + math.Exp(-logit))
}
