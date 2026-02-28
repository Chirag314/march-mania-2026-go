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
	var temp float64

	flag.StringVar(&artDir, "art_dir", "artifacts", "artifacts directory")
	flag.StringVar(&outDir, "out_dir", "submissions", "output directory")
	flag.Float64Var(&temp, "temp", 1.0, "temperature scaling (t>1 softens probs)")
	flag.Parse()

	modelPath := filepath.Join(artDir, "model.json")
	model, err := mm.LoadModelJSON(modelPath)
	must(err)
	fmt.Println("Loaded model:", modelPath)

	testPath := filepath.Join(artDir, "features_test.csv")
	fmt.Println("testPath:", testPath)

	rows, err := mm.ReadMatchupsCSV(testPath)
	must(err)
	fmt.Println("rows read:", len(rows))
	if len(rows) == 0 {
		panic("no rows read from features_test.csv (check ReadMatchupsCSV and path)")
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(err)
	}

	outPath := filepath.Join(outDir, "submission.csv")
	f, err := os.Create(outPath)
	must(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	must(w.Write([]string{"ID", "Pred"}))

	for _, r := range rows {
		x := []float64{
			r.DSeed, r.DElo, r.DWinPct, r.DAvgMargin, r.DAvgPF, r.DAvgPA, r.DMasseyOrd,
		}
		p := model.PredictProba(x)
		p = mm.TemperatureScale(p, temp)
		p = mm.ClipProb(p, 0.02, 0.98)

		must(w.Write([]string{r.ID, fmt.Sprintf("%.6f", p)}))
	}

	fmt.Println("Wrote:", outPath)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
