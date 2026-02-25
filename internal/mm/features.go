package mm

import (
	"os"
	"path/filepath"
)

func BuildTeamSeasonAgg(
	regular []RegularSeasonCompactRow,
	seeds []SeedRow,
	massey []MasseyRow,
) map[[2]int]*TeamSeasonAgg {

	agg := make(map[[2]int]*TeamSeasonAgg)

	get := func(season, team int) *TeamSeasonAgg {
		k := [2]int{season, team}
		a, ok := agg[k]
		if !ok {
			a = &TeamSeasonAgg{Season: season, TeamID: team}
			agg[k] = a
		}
		return a
	}

	for _, g := range regular {
		w := get(g.Season, g.WTeamID)
		l := get(g.Season, g.LTeamID)

		w.Games++
		w.Wins++
		l.Games++
		l.Losses++

		w.PointsFor += float64(g.WScore)
		w.PointsAgainst += float64(g.LScore)
		w.MarginSum += float64(g.WScore - g.LScore)

		l.PointsFor += float64(g.LScore)
		l.PointsAgainst += float64(g.WScore)
		l.MarginSum += float64(g.LScore - g.WScore)
	}

	for _, a := range agg {
		if a.Games > 0 {
			a.WinPct = float64(a.Wins) / float64(a.Games)
			a.AvgPF = a.PointsFor / float64(a.Games)
			a.AvgPA = a.PointsAgainst / float64(a.Games)
			a.AvgMargin = a.MarginSum / float64(a.Games)
		}
	}

	for _, s := range seeds {
		a := get(s.Season, s.TeamID)
		a.Seed = float64(s.Seed)
	}

	// simple Massey: last day ordinal over all systems
	lastDay := make(map[[2]int]int)
	lastOrd := make(map[[2]int]float64)
	for _, m := range massey {
		k := [2]int{m.Season, m.TeamID}
		if m.RankingDay >= lastDay[k] {
			lastDay[k] = m.RankingDay
			lastOrd[k] = float64(m.Ordinal)
		}
	}
	for k, ord := range lastOrd {
		if a, ok := agg[k]; ok {
			a.MasseyOrdinal = ord
		} else {
			agg[k] = &TeamSeasonAgg{Season: k[0], TeamID: k[1], MasseyOrdinal: ord}
		}
	}

	return agg
}

func AttachEloEnd(agg map[[2]int]*TeamSeasonAgg, eloEnd map[[2]int]float64) {
	for k, r := range eloEnd {
		a, ok := agg[k]
		if !ok {
			a = &TeamSeasonAgg{Season: k[0], TeamID: k[1]}
			agg[k] = a
		}
		a.EloEnd = r
	}
}

func WriteTeamSeasonAggCSV(outDir string, agg map[[2]int]*TeamSeasonAgg) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(outDir, "team_season_features.csv")
	w, err := NewCSVWriter(path, []string{
		"Season", "TeamID",
		"Games", "Wins", "Losses",
		"WinPct", "AvgPF", "AvgPA", "AvgMargin",
		"EloEnd", "Seed", "MasseyOrdinal",
	})
	if err != nil {
		return "", err
	}
	defer w.Close()

	for _, a := range agg {
		w.WriteRow([]string{
			fmtInt(a.Season),
			fmtInt(a.TeamID),
			fmtInt(a.Games),
			fmtInt(a.Wins),
			fmtInt(a.Losses),
			fmtF(a.WinPct),
			fmtF(a.AvgPF),
			fmtF(a.AvgPA),
			fmtF(a.AvgMargin),
			fmtF(a.EloEnd),
			fmtF(a.Seed),
			fmtF(a.MasseyOrdinal),
		})
	}
	return path, nil
}
