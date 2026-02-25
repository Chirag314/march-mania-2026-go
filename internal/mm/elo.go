package mm

import "math"

type EloConfig struct {
	K     float64
	Start float64
}

func DefaultEloConfig() EloConfig {
	return EloConfig{K: 20.0, Start: 1500.0}
}

func expectedScore(ra, rb float64) float64 {
	return 1.0 / (1.0 + math.Pow(10.0, (rb-ra)/400.0))
}

func BuildEloEnd(rows []RegularSeasonCompactRow, cfg EloConfig) map[[2]int]float64 {
	rating := make(map[[2]int]float64)

	get := func(season, team int) float64 {
		k := [2]int{season, team}
		v, ok := rating[k]
		if !ok {
			v = cfg.Start
			rating[k] = v
		}
		return v
	}
	set := func(season, team int, v float64) {
		rating[[2]int{season, team}] = v
	}

	for _, g := range rows {
		ra := get(g.Season, g.WTeamID)
		rb := get(g.Season, g.LTeamID)

		ea := expectedScore(ra, rb)
		eb := 1.0 - ea

		ra2 := ra + cfg.K*(1.0-ea)
		rb2 := rb + cfg.K*(0.0-eb)

		set(g.Season, g.WTeamID, ra2)
		set(g.Season, g.LTeamID, rb2)
	}
	return rating
}
