package mm

type RegularSeasonCompactRow struct {
	Season  int
	DayNum  int
	WTeamID int
	WScore  int
	LTeamID int
	LScore  int
	WLoc    string
	NumOT   int
}

type TourneyCompactRow struct {
	Season  int
	DayNum  int
	WTeamID int
	WScore  int
	LTeamID int
	LScore  int
	WLoc    string
	NumOT   int
}

type SeedRow struct {
	Season int
	TeamID int
	Seed   int
}

type MasseyRow struct {
	Season     int
	TeamID     int
	RankingDay int
	System     string
	Ordinal    int
}

type TeamSeasonAgg struct {
	Season int
	TeamID int

	Games  int
	Wins   int
	Losses int

	PointsFor     float64
	PointsAgainst float64
	MarginSum     float64

	WinPct    float64
	AvgPF     float64
	AvgPA     float64
	AvgMargin float64

	EloEnd        float64
	Seed          float64
	MasseyOrdinal float64
}

type MatchupFeatureRow struct {
	ID     string
	Season int
	TeamA  int
	TeamB  int

	DSeed      float64
	DElo       float64
	DWinPct    float64
	DAvgMargin float64
	DAvgPF     float64
	DAvgPA     float64
	DMasseyOrd float64

	Label    float64
	HasLabel bool
}
