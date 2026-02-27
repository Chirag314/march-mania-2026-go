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

func findFile(dataDir string, candidates ...string) (string, error) {
	for _, c := range candidates {
		p := filepath.Join(dataDir, c)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	entries, _ := os.ReadDir(dataDir)
	for _, e := range entries {
		name := e.Name()
		for _, c := range candidates {
			if strings.EqualFold(name, c) {
				return filepath.Join(dataDir, name), nil
			}
		}
	}
	return "", fmt.Errorf("file not found in %s: %v", dataDir, candidates)
}

func atoi(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty int")
	}
	return strconv.Atoi(s)
}
func atof(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty float")
	}
	return strconv.ParseFloat(s, 64)
}

func indexMap(header []string) map[string]int {
	m := make(map[string]int, len(header))
	for i, h := range header {
		m[strings.TrimSpace(h)] = i
	}
	return m
}

func getStr(rec []string, col map[string]int, name string) string {
	i, ok := col[name]
	if !ok || i >= len(rec) {
		return ""
	}
	return strings.TrimSpace(rec[i])
}

func mustInt(rec []string, col map[string]int, name string) int {
	v := getStr(rec, col, name)
	x, err := atoi(v)
	if err != nil {
		panic(fmt.Sprintf("parse int %s=%q: %v", name, v, err))
	}
	return x
}

func getIntDefault(rec []string, col map[string]int, name string, def int) int {
	v := getStr(rec, col, name)
	if v == "" {
		return def
	}
	x, err := atoi(v)
	if err != nil {
		return def
	}
	return x
}

func mustFloat(rec []string, col map[string]int, name string) float64 {
	v := getStr(rec, col, name)
	x, err := atof(v)
	if err != nil {
		panic(fmt.Sprintf("parse float %s=%q: %v", name, v, err))
	}
	return x
}

func ReadRegularSeasonCompact(dataDir string) ([]RegularSeasonCompactRow, error) {
	path, err := findFile(dataDir,
		"MRegularSeasonCompactResults.csv",
		"WRegularSeasonCompactResults.csv",
	)
	if err != nil {
		return nil, err
	}
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

	var out []RegularSeasonCompactRow
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		row := RegularSeasonCompactRow{}
		row.Season = mustInt(rec, col, "Season")
		row.DayNum = mustInt(rec, col, "DayNum")
		row.WTeamID = mustInt(rec, col, "WTeamID")
		row.WScore = mustInt(rec, col, "WScore")
		row.LTeamID = mustInt(rec, col, "LTeamID")
		row.LScore = mustInt(rec, col, "LScore")
		row.WLoc = getStr(rec, col, "WLoc")
		row.NumOT = getIntDefault(rec, col, "NumOT", 0)

		out = append(out, row)
	}
	return out, nil
}

func ReadTourneyCompact(dataDir string) ([]TourneyCompactRow, error) {
	path, err := findFile(dataDir,
		"MNCAATourneyCompactResults.csv",
		"WNCAATourneyCompactResults.csv",
	)
	if err != nil {
		return nil, err
	}
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

	var out []TourneyCompactRow
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		row := TourneyCompactRow{}
		row.Season = mustInt(rec, col, "Season")
		row.DayNum = mustInt(rec, col, "DayNum")
		row.WTeamID = mustInt(rec, col, "WTeamID")
		row.WScore = mustInt(rec, col, "WScore")
		row.LTeamID = mustInt(rec, col, "LTeamID")
		row.LScore = mustInt(rec, col, "LScore")
		row.WLoc = getStr(rec, col, "WLoc")
		row.NumOT = getIntDefault(rec, col, "NumOT", 0)

		out = append(out, row)
	}
	return out, nil
}

func ReadSeeds(dataDir string) ([]SeedRow, error) {
	path, err := findFile(dataDir,
		"MNCAATourneySeeds.csv",
		"WNCAATourneySeeds.csv",
	)
	if err != nil {
		return nil, err
	}
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

	var out []SeedRow
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		seedStr := getStr(rec, col, "Seed")
		seedNum := parseSeedNumeric(seedStr)

		out = append(out, SeedRow{
			Season: mustInt(rec, col, "Season"),
			TeamID: mustInt(rec, col, "TeamID"),
			Seed:   seedNum,
		})
	}
	return out, nil
}

func ReadMassey(dataDir string) ([]MasseyRow, error) {
	path, err := findFile(dataDir, "MMasseyOrdinals.csv", "WMasseyOrdinals.csv")
	if err != nil {
		return nil, nil // optional
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, nil
	}
	col := indexMap(header)

	var out []MasseyRow
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil
		}
		out = append(out, MasseyRow{
			Season:     mustInt(rec, col, "Season"),
			TeamID:     mustInt(rec, col, "TeamID"),
			RankingDay: mustInt(rec, col, "RankingDayNum"),
			System:     getStr(rec, col, "SystemName"),
			Ordinal:    mustInt(rec, col, "OrdinalRank"),
		})
	}
	return out, nil
}

func ReadSampleSubmission(dataDir string) ([]string, error) {
	path, err := findFile(dataDir,
		"SampleSubmissionStage1.csv",
		"SampleSubmissionStage2.csv",
		"MSampleSubmissionStage1.csv",
		"MSampleSubmissionStage2.csv",
	)
	if err != nil {
		return nil, err
	}
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

	var out []string
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		id := getStr(rec, col, "ID")
		if id != "" {
			out = append(out, id)
		}
	}
	return out, nil
}

func parseSeedNumeric(seed string) int {
	seed = strings.TrimSpace(seed)
	if seed == "" {
		return 0
	}
	var digits strings.Builder
	for _, ch := range seed {
		if ch >= '0' && ch <= '9' {
			digits.WriteRune(ch)
		}
	}
	if digits.Len() == 0 {
		return 0
	}
	n, err := strconv.Atoi(digits.String())
	if err != nil {
		return 0
	}
	return n
}
