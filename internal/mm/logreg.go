package mm

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

type LogRegModel struct {
	FeatureNames []string  `json:"feature_names"`
	Weights      []float64 `json:"weights"` // bias at index 0
}

func sigmoid(x float64) float64 {
	if x < -50 {
		return 1e-22
	}
	if x > 50 {
		return 1 - 1e-22
	}
	return 1.0 / (1.0 + math.Exp(-x))
}

func (m *LogRegModel) PredictProba(x []float64) float64 {
	z := m.Weights[0]
	for i := 0; i < len(x) && i+1 < len(m.Weights); i++ {
		z += m.Weights[i+1] * x[i]
	}
	return sigmoid(z)
}

type TrainConfig struct {
	LR     float64
	Epochs int
	L2     float64
}

func DefaultTrainConfig() TrainConfig {
	return TrainConfig{LR: 0.05, Epochs: 400, L2: 1e-4}
}

func TrainLogReg(X [][]float64, y []float64, featureNames []string, cfg TrainConfig) *LogRegModel {
	if len(X) == 0 {
		return &LogRegModel{FeatureNames: featureNames, Weights: make([]float64, 1+len(featureNames))}
	}
	d := len(X[0])
	w := make([]float64, d+1)

	n := float64(len(X))
	for epoch := 0; epoch < cfg.Epochs; epoch++ {
		grad := make([]float64, d+1)

		for i := 0; i < len(X); i++ {
			p := sigmoid(dotBias(w, X[i]))
			err := p - y[i]
			grad[0] += err
			for j := 0; j < d; j++ {
				grad[j+1] += err * X[i][j]
			}
		}

		grad[0] /= n
		for j := 1; j < d+1; j++ {
			grad[j] = grad[j]/n + cfg.L2*w[j]
		}

		for j := 0; j < d+1; j++ {
			w[j] -= cfg.LR * grad[j]
		}
	}

	return &LogRegModel{FeatureNames: featureNames, Weights: w}
}

func dotBias(w []float64, x []float64) float64 {
	z := w[0]
	for i := 0; i < len(x) && i+1 < len(w); i++ {
		z += w[i+1] * x[i]
	}
	return z
}

func SaveModelJSON(path string, m *LogRegModel) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func LoadModelJSON(path string) (*LogRegModel, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m LogRegModel
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if len(m.Weights) != 1+len(m.FeatureNames) {
		return nil, fmt.Errorf("bad model: weights=%d features=%d", len(m.Weights), len(m.FeatureNames))
	}
	return &m, nil
}
