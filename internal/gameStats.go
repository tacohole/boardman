package internal

import "github.com/google/uuid"

type SingleGame struct {
	PlayerID uuid.UUID `db:"player_uuid"`
	GameID   uuid.UUID `db:"game_uuid"`
	TeamID   int       `db:"team_id"`
	Minutes  string    `json:"min" db:"min"`
	FGM      int       `json:"fgm" db:"fgm"`
	FGA      int       `json:"fga" db:"fga"`
	FG3M     int       `json:"fg3m" db:"fg3m"`
	FG3A     int       `json:"fg3a" db:"fg3a"`
	OREB     int       `json:"oreb" db:"oreb"`
	DREB     int       `json:"dreb" db:"dreb"`
	REB      int       `json:"reb" db:"reb"`
	AST      int       `json:"ast" db:"ast"`
	STL      int       `json:"stl" db:"stl"`
	BLK      int       `json:"blk" db:"blk"`
	TO       int       `json:"turnover" db:"to"`
	PF       int       `json:"pf" db:"pf"`
	PTS      int       `json:"pts" db:"pts"`
	FG_PCT   int       `json:"fg_pct" db:"fg_pct"`
	FG3_PCT  int       `json:"fg3_pct" db:"fg3_pct"`
	FT_PCT   int       `json:"ft_pct" db:"ft_pct"`
}
