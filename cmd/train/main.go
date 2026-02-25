package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/Chirag314/march-mania-2026-go/internal/mm"
)

func main() {
	var artDir, outDir string
	var k int
	var seed int64

	flag.StringVar(&artDir, "art_dir", "artifacts", "artifacts directory")
	flag.StringVar(&outDir, "out_dir", "artifacts", "output directory")
	flag.IntVar(&k, "k", 5, "number of CV folds (season-grouped)")
	flag.Int64Var(&seed, "seed", 42, "random seed for season assignment")
	flag.Parse()

	trainPath := filepath.Join(artDir, "features_train.csv")
	fmt.Println("Reading:", trainPath)

	rows, err := mm.ReadMatchupsCSV(trainPath)
	must(err)

	featureNames := []string{"DSeed", "DElo", "DWinPct", "DAvgMargin", "DAvgPF", "DAvgPA", "DMasseyOrd"}

	X, y := toXY(rows)

	folds := mm.SeasonGroupFolds(rows, k, seed)
	fmt.Printf("CV: %d folds (season-grouped)\n", len(folds))

	var foldScores []float64
	for fi, fold := range folds {
		Xtr, ytr := subsetXY(X, y, fold.TrainIdx)
		Xva, yva := subsetXY(X, y, fold.ValIdx)

		model := mm.TrainLogReg(Xtr, ytr, featureNames, mm.DefaultTrainConfig())

		pred := make([]float64, len(Xva))
		for i := range Xva {
			pred[i] = model.PredictProba(Xva[i])
		}

		score := mm.BrierScore(yva, pred)
		foldScores = append(foldScores, score)
		fmt.Printf("Fold %d: n_val=%d Brier=%.6f\n", fi+1, len(fold.ValIdx), score)
	}

	mean, std := mm.MeanStd(foldScores)
	fmt.Printf("CV Brier mean=%.6f std=%.6f\n", mean, std)

	finalModel := mm.TrainLogReg(X, y, featureNames, mm.DefaultTrainConfig())
	modelPath := filepath.Join(outDir, "model.json")
	must(mm.SaveModelJSON(modelPath, finalModel))
	fmt.Println("Saved model:", modelPath)
}

func toXY(rows []mm.MatchupFeatureRow) (X [][]float64, y []float64) {
	for _, r := range rows {
		if !r.HasLabel {
			continue
		}
		x := []float64{r.DSeed, r.DElo, r.DWinPct, r.DAvgMargin, r.DAvgPF, r.DAvgPA, r.DMasseyOrd}
		X = append(X, x)
		y = append(y, r.Label)
	}
	return
}

func subsetXY(X [][]float64, y []float64, idx []int) ([][]float64, []float64) {
	outX := make([][]float64, 0, len(idx))
	outY := make([]float64, 0, len(idx))
	for _, i := range idx {
		outX = append(outX, X[i])
		outY = append(outY, y[i])
	}
	return outX, outY
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
