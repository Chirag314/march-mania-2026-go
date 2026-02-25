package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/Chirag314/march-mania-2026-go/internal/mm"
)

func main() {
	var dataDir, outDir string
	flag.StringVar(&dataDir, "data_dir", "data", "directory with Kaggle CSV files")
	flag.StringVar(&outDir, "out_dir", "artifacts", "output directory")
	flag.Parse()

	fmt.Println("Reading regular season...")
	reg, err := mm.ReadRegularSeasonCompact(dataDir)
	must(err)

	fmt.Println("Reading tourney results...")
	tour, err := mm.ReadTourneyCompact(dataDir)
	must(err)

	fmt.Println("Reading seeds...")
	seeds, err := mm.ReadSeeds(dataDir)
	must(err)

	fmt.Println("Reading Massey (optional)...")
	massey, _ := mm.ReadMassey(dataDir)

	fmt.Println("Building Elo...")
	eloEnd := mm.BuildEloEnd(reg, mm.DefaultEloConfig())

	fmt.Println("Aggregating team-season features...")
	agg := mm.BuildTeamSeasonAgg(reg, seeds, massey)
	mm.AttachEloEnd(agg, eloEnd)

	_, err = mm.WriteTeamSeasonAggCSV(outDir, agg)
	must(err)

	fmt.Println("Building train matchups from tourney...")
	train := mm.BuildTrainMatchupsFromTourney(tour)
	train = mm.JoinFeatures(train, agg)

	fmt.Println("Reading sample submission IDs...")
	ids, err := mm.ReadSampleSubmission(dataDir)
	must(err)

	fmt.Println("Building test matchups from sample submission...")
	test, err := mm.BuildTestMatchupsFromIDs(ids)
	must(err)
	test = mm.JoinFeatures(test, agg)

	_, err = mm.WriteMatchupsCSV(outDir, "features_train.csv", train)
	must(err)
	_, err = mm.WriteMatchupsCSV(outDir, "features_test.csv", test)
	must(err)

	fmt.Printf("Done.\n- %s\n- %s\n",
		filepath.Join(outDir, "features_train.csv"),
		filepath.Join(outDir, "features_test.csv"),
	)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
