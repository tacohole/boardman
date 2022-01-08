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

func PreparePlayerStatsSchema() error {
	schema := `CREATE TABLE player_season_avgs(
		player_id UUID,
		season INT,
		avg_min NUMERIC(4,2),
		fgm NUMERIC(5,2),
		fga NUMERIC(5,2),
		fg3m NUMERIC(5,2),
		fg3a NUMERIC(5,2),
		oreb NUMERIC(5,2),
		dreb NUMERIC(5,2),
		reb NUMERIC(5,2),
		ast NUMERIC(5,2),
		stl NUMERIC(5,2),
		blk NUMERIC(5,2),
		to NUMERIC(5,2),
		pf NUMERIC(4,2),
		pts NUMERIC(5,2),
		fg_pct NUMERIC(4,3),
		fg3_pct NUMERIC(4,3),
		ft_pct NUMERIC(4,3),
		CONSTRAINT fk_players
		FOREIGN KEY(player_id)
		REFERENCES players(uuid)
	);`

	db, err := dbutil.DbConn()
	if err != nil {
		return err
	}

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	db.MustExecContext(ctx, schema)

	return nil
}
