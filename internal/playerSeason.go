package internal

import (
	"context"

	"github.com/google/uuid"
	dbutil "github.com/tacohole/boardman/util/db"
)

type PlayerSeason struct {
	PlayerUUID  uuid.UUID `db:"uuid"`
	BDL_ID      int       `json:"player_id" db:"balldontlie_id"`
	LeagueYear  int       `json:"season" db:"season"`
	GamesPlayed int       `json:"games_played" db:"games_played"`
	Minutes     string    `json:"min" db:"avg_min"`
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
	TO          float32   `json:"to" db:"turnovers"`
	PF          float32   `json:"pf" db:"pf"`
	PTS         float32   `json:"pts" db:"pts"`
	FG_PCT      float32   `json:"fg_pct" db:"fg_pct"`
	FG3_PCT     float32   `json:"fg3_pct" db:"fg3_pct"`
	FT_PCT      float32   `json:"ft_pct" db:"ft_pct"`
}

func PreparePlayerSeasonSchema() error {
	schema := `CREATE TABLE IF NOT EXISTS player_season(
		uuid UUID,
		balldontlie_id INT,
		season INT,
		avg_min TEXT,
		fgm NUMERIC(4),
		fga NUMERIC(4),
		fg3m NUMERIC(4),
		fg3a NUMERIC(4),
		oreb NUMERIC(4),
		dreb NUMERIC(4),
		reb NUMERIC(4),
		ast NUMERIC(4),
		stl NUMERIC(4),
		blk NUMERIC(4),
		turnovers NUMERIC(4),
		pf NUMERIC(4),
		pts NUMERIC(4),
		fg_pct NUMERIC(4),
		fg3_pct NUMERIC(4),
		ft_pct NUMERIC(4),
		CONSTRAINT fk_players
		FOREIGN KEY(uuid)
		REFERENCES players(uuid));`

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
