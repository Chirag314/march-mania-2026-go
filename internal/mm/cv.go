package mm

import (
	"math/rand"
	"sort"
)

type Fold struct {
	TrainIdx  []int
	ValIdx    []int
	ValSeason int
}

func SeasonGroupFolds(rows []MatchupFeatureRow, k int, seed int64) []Fold {
	seasonsSet := map[int]struct{}{}
	for _, r := range rows {
		if r.HasLabel {
			seasonsSet[r.Season] = struct{}{}
		}
	}
	seasons := make([]int, 0, len(seasonsSet))
	for s := range seasonsSet {
		seasons = append(seasons, s)
	}
	sort.Ints(seasons)

	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(seasons), func(i, j int) { seasons[i], seasons[j] = seasons[j], seasons[i] })

	if k <= 1 {
		k = 5
	}
	buckets := make([][]int, k)
	for i, s := range seasons {
		buckets[i%k] = append(buckets[i%k], s)
	}

	folds := make([]Fold, 0, k)
	for fi := 0; fi < k; fi++ {
		valSeasons := map[int]struct{}{}
		for _, s := range buckets[fi] {
			valSeasons[s] = struct{}{}
		}

		var trainIdx, valIdx []int
		for i, r := range rows {
			if !r.HasLabel {
				continue
			}
			_, isVal := valSeasons[r.Season]
			if isVal {
				valIdx = append(valIdx, i)
			} else {
				trainIdx = append(trainIdx, i)
			}
		}

		valSeason := -1
		if len(buckets[fi]) > 0 {
			valSeason = buckets[fi][0]
		}
		folds = append(folds, Fold{TrainIdx: trainIdx, ValIdx: valIdx, ValSeason: valSeason})
	}
	return folds
}
