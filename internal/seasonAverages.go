package internal

import (
	"context"

	"github.com/google/uuid"
	dbutil "github.com/tacohole/boardman/util/db"
)

type PlayerYear struct {
	PlayerID    uuid.UUID `db:"uuid"`
	BDL_ID      int       `json:"player_id" db:"balldontlie_id"`
	LeagueYear  int       `json:"season" db:"league_year"`
	GamesPlayed int       `json:"games_played" db:"games_played"`
	Minutes     string    `json:"avg_min" db:"avg_min"`
	FGM         float32   `json:"fgm" db:"fgm"`
	FGA         float32   `json:"fga" db:"fga"`
	FG3M        float32   `json:"fg3m" db:"fg3m"`
	FG3A        float32   `json:"fg3a" db:"fg3a"`
	OREB        float32   `json:"oreb" db:"oreb"`
	DREB        float32   `json:"dreb" db:"dreb"`
	REB         float32   `json:"reb" db:"reb"`
	AST         float32   `json:"ast" db:"ast"`
	STL         float32   `json:"stl" db:"stl"`
	BLK         float32   `json:"blk" db:"blk"`
	TO          float32   `json:"to" db:"to"`
	PF          float32   `json:"pf" db:"pf"`
	PTS         float32   `json:"pts" db:"pts"`
	FG_PCT      float32   `json:"fg_pct" db:"fg_pct"`
	FG3_PCT     float32   `json:"fg3_pct" db:"fg3_pct"`
	FT_PCT      float32   `json:"ft_pct" db:"ft_pct"`
}

// lookup our UUID off BDL_ID
func (py *PlayerYear) getUUIDFromBDLID() (*uuid.UUID, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}
	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	_, err = db.NamedExecContext(ctx, "SELECT id FROM players WHERE players(balldontlie_id)=:balldontlie_id", py)
	if err != nil {
		return nil, err
	}

	return &py.PlayerID, nil
}
