package mm

import "math"

func BrierScore(yTrue, yPred []float64) float64 {
	if len(yTrue) == 0 {
		return math.NaN()
	}
	var sum float64
	for i := range yTrue {
		d := yPred[i] - yTrue[i]
		sum += d * d
	}
	return sum / float64(len(yTrue))
}

func MeanStd(xs []float64) (mean float64, std float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	for _, x := range xs {
		mean += x
	}
	mean /= float64(len(xs))
	for _, x := range xs {
		d := x - mean
		std += d * d
	}
	std = math.Sqrt(std / float64(len(xs)))
	return mean, std
}
