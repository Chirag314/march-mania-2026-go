package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Chirag314/march-mania-2026-go/internal/mm"
)

func main() {
	var artDir, outDir string
	var minSeason int
	var temp float64
	flag.Float64Var(&temp, "temp", 1.0, "temperature scaling (t>1 softens probs)")
	flag.IntVar(&minSeason, "min_season", 1985, "minimum season to include in training/CV")
	flag.StringVar(&artDir, "art_dir", "artifacts", "artifacts directory")
	flag.StringVar(&outDir, "out_dir", "submissions", "output directory")
	flag.Parse()

	modelPath := filepath.Join(artDir, "model.json")
	model, err := mm.LoadModelJSON(modelPath)
	must(err)
	fmt.Println("Loaded model:", modelPath)

	testPath := filepath.Join(artDir, "features_test.csv")
	rows, err := mm.ReadMatchupsCSV(testPath)
	must(err)

	rows = filterMinSeasonLabeled(rows, minSeason)
	fmt.Println("Train/CV rows after min_season filter:", len(rows))

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(err)
	}

	outPath := filepath.Join(outDir, "submission.csv")
	f, err := os.Create(outPath)
	must(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	_ = w.Write([]string{"ID", "Pred"})
	for _, r := range rows {
		x := []float64{r.DSeed, r.DElo, r.DWinPct, r.DAvgMargin, r.DAvgPF, r.DAvgPA, r.DMasseyOrd}
		p := model.PredictProba(x)
		p = mm.TemperatureScale(p, temp) // <- new
		p = mm.ClipProb(p, 0.02, 0.98)   // keep clipping
		_ = w.Write([]string{r.ID, fmt.Sprintf("%.6f", p)})
	}

	fmt.Println("Wrote:", outPath)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func filterMinSeasonLabeled(rows []mm.MatchupFeatureRow, minSeason int) []mm.MatchupFeatureRow {
	out := make([]mm.MatchupFeatureRow, 0, len(rows))
	for _, r := range rows {
		if !r.HasLabel {
			continue
		}
		if r.Season < minSeason {
			continue
		}
		out = append(out, r)
	}
	return out
}
