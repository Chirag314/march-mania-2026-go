package mm

func ClipProb(p, lo, hi float64) float64 {
	if p < lo {
		return lo
	}
	if p > hi {
		return hi
	}
	return p
}
