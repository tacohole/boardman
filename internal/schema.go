package internal

type Season struct {
	LeagueYear int
	Champion   Team
	WConfChamp Team
	EConfChamp Team
	MVP        Player
}

// get champ
// get conf champ
// get mvp

type TeamYear struct {
	TeamCache    Team
	Season       Season
	Wins         int
	Losses       int
	WinPct       int
	ConfRank     int
	OvrRank      int
	MadePlayoffs bool
	Roster       []Player
	Coach        string
}

// sum wins
// sum losses
// sum wpct
// get conf rank
// build roster

type PlayerYear struct {
	Player  Player
	Season  Season
	Minutes string  `json:"avg_min" db:"avg_min"`
	FGM     float32 `json:"fgm" db:"fgm"`
	FGA     float32 `json:"fga" db:"fga"`
	FG3M    float32 `json:"fg3m" db:"fg3m"`
	FG3A    float32 `json:"fg3a" db:"fg3a"`
	OREB    float32 `json:"oreb" db:"oreb"`
	DREB    float32 `json:"dreb" db:"dreb"`
	REB     float32 `json:"reb" db:"reb"`
	AST     float32 `json:"ast" db:"ast"`
	STL     float32 `json:"stl" db:"stl"`
	BLK     float32 `json:"blk" db:"blk"`
	TO      float32 `json:"to" db:"to"`
	PF      float32 `json:"pf" db:"pf"`
	PTS     float32 `json:"pts" db:"pts"`
	FG_PCT  float32 `json:"fg_pct" db:"fg_pct"`
	FG3_PCT float32 `json:"fg3_pct" db:"fg3_pct"`
	FT_PCT  float32 `json:"ft_pct" db:"ft_pct"`
}
