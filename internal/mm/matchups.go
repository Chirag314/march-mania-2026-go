package mm

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseMatchupID(id string) (season int, teamA int, teamB int, err error) {
	parts := strings.Split(id, "_")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("bad ID format: %q", id)
	}
	season, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, err
	}
	teamA, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, err
	}
	teamB, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, err
	}
	return season, teamA, teamB, nil
}

func BuildTrainMatchupsFromTourney(tourney []TourneyCompactRow) []MatchupFeatureRow {
	var out []MatchupFeatureRow
	for _, g := range tourney {
		id1 := fmt.Sprintf("%d_%d_%d", g.Season, g.WTeamID, g.LTeamID)
		out = append(out, MatchupFeatureRow{
			ID: id1, Season: g.Season, TeamA: g.WTeamID, TeamB: g.LTeamID,
			Label: 1.0, HasLabel: true,
		})
		id2 := fmt.Sprintf("%d_%d_%d", g.Season, g.LTeamID, g.WTeamID)
		out = append(out, MatchupFeatureRow{
			ID: id2, Season: g.Season, TeamA: g.LTeamID, TeamB: g.WTeamID,
			Label: 0.0, HasLabel: true,
		})
	}
	return out
}

func BuildTestMatchupsFromIDs(ids []string) ([]MatchupFeatureRow, error) {
	out := make([]MatchupFeatureRow, 0, len(ids))
	for _, id := range ids {
		season, a, b, err := ParseMatchupID(id)
		if err != nil {
			return nil, err
		}
		out = append(out, MatchupFeatureRow{ID: id, Season: season, TeamA: a, TeamB: b, HasLabel: false})
	}
	return out, nil
}

func JoinFeatures(matchups []MatchupFeatureRow, agg map[[2]int]*TeamSeasonAgg) []MatchupFeatureRow {
	get := func(season, team int) *TeamSeasonAgg {
		a, ok := agg[[2]int{season, team}]
		if !ok {
			return &TeamSeasonAgg{Season: season, TeamID: team}
		}
		return a
	}

	for i := range matchups {
		m := &matchups[i]
		a := get(m.Season, m.TeamA)
		b := get(m.Season, m.TeamB)

		m.DSeed = a.Seed - b.Seed
		m.DElo = a.EloEnd - b.EloEnd
		m.DWinPct = a.WinPct - b.WinPct
		m.DAvgMargin = a.AvgMargin - b.AvgMargin
		m.DAvgPF = a.AvgPF - b.AvgPF
		m.DAvgPA = a.AvgPA - b.AvgPA
		m.DMasseyOrd = a.MasseyOrdinal - b.MasseyOrdinal
	}
	return matchups
}

func WriteMatchupsCSV(outDir, name string, rows []MatchupFeatureRow) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(outDir, name)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"ID", "Season", "TeamA", "TeamB",
		"DSeed", "DElo", "DWinPct", "DAvgMargin", "DAvgPF", "DAvgPA", "DMasseyOrd",
		"Label", "HasLabel",
	}
	if err := w.Write(header); err != nil {
		return "", err
	}
	for _, r := range rows {
		rec := []string{
			r.ID,
			strconv.Itoa(r.Season),
			strconv.Itoa(r.TeamA),
			strconv.Itoa(r.TeamB),
			fmtF(r.DSeed),
			fmtF(r.DElo),
			fmtF(r.DWinPct),
			fmtF(r.DAvgMargin),
			fmtF(r.DAvgPF),
			fmtF(r.DAvgPA),
			fmtF(r.DMasseyOrd),
			fmtF(r.Label),
			strconv.FormatBool(r.HasLabel),
		}
		if err := w.Write(rec); err != nil {
			return "", err
		}
	}
	return path, nil
}

func ReadMatchupsCSV(path string) ([]MatchupFeatureRow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	col := indexMap(header)

	// helper: check if a column exists
	_, hasLabelCol := col["Label"]
	_, hasHasLabelCol := col["HasLabel"]

	var out []MatchupFeatureRow
	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		row := MatchupFeatureRow{
			ID:         getStr(rec, col, "ID"),
			Season:     mustInt(rec, col, "Season"),
			TeamA:      mustInt(rec, col, "TeamA"),
			TeamB:      mustInt(rec, col, "TeamB"),
			DSeed:      mustFloat(rec, col, "DSeed"),
			DElo:       mustFloat(rec, col, "DElo"),
			DWinPct:    mustFloat(rec, col, "DWinPct"),
			DAvgMargin: mustFloat(rec, col, "DAvgMargin"),
			DAvgPF:     mustFloat(rec, col, "DAvgPF"),
			DAvgPA:     mustFloat(rec, col, "DAvgPA"),
			DMasseyOrd: mustFloat(rec, col, "DMasseyOrd"),
			Label:      0,
			HasLabel:   false,
		}

		// HasLabel is optional / may be false in test
		hs := strings.ToLower(getStr(rec, col, "HasLabel"))
		if hs == "true" {
			row.HasLabel = true
		}

		// Label optional: only parse if non-empty
		ls := getStr(rec, col, "Label")
		if strings.TrimSpace(ls) != "" {
			if v, err := atof(ls); err == nil {
				row.Label = v
			}
		}

		out = append(out, row)
	}

	return out, nil
}
