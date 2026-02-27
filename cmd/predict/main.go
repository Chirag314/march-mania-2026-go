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
		p = mm.ClipProb(p, 0.02, 0.98)
		_ = w.Write([]string{r.ID, fmt.Sprintf("%.6f", p)})
	}

	fmt.Println("Wrote:", outPath)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
