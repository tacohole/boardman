package internal

import (
	"context"

	"github.com/google/uuid"
	dbutil "github.com/tacohole/boardman/util/db"
)

type SingleGame struct {
	UUID         uuid.UUID `db:"uuid"`
	BDL_ID       int       `json:"id" db:"balldontlie_id"`
	GameUUID     uuid.UUID `db:"game_uuid"`
	GameBDL_ID   int       `db:"game_bdl_id"`
	PlayerUUID   uuid.UUID `db:"player_uuid"`
	PlayerBDL_ID int       `db:"player_bdl_id"`
	TeamUUID     uuid.UUID `db:"team_uuid"`
	TeamBDL_ID   int       `db:"team_bdl_id"`
	Minutes      string    `json:"min" db:"min"`
	FGM          float32   `json:"fgm" db:"fgm"`
	FGA          float32   `json:"fga" db:"fga"`
	FG3M         float32   `json:"fg3m" db:"fg3m"`
	FG3A         float32   `json:"fg3a" db:"fg3a"`
	OREB         float32   `json:"oreb" db:"oreb"`
	DREB         float32   `json:"dreb" db:"dreb"`
	REB          float32   `json:"reb" db:"reb"`
	AST          float32   `json:"ast" db:"ast"`
	STL          float32   `json:"stl" db:"stl"`
	BLK          float32   `json:"blk" db:"blk"`
	TO           float32   `json:"turnover" db:"turnovers"`
	PF           float32   `json:"pf" db:"pf"`
	PTS          float32   `json:"pts" db:"pts"`
	FG_PCT       float32   `json:"fg_pct" db:"fg_pct"`
	FG3_PCT      float32   `json:"fg3_pct" db:"fg3_pct"`
	FT_PCT       float32   `json:"ft_pct" db:"ft_pct"`
}

func PrepareGameStatsSchema() error {
	schema := `CREATE TABLE player_game_stats(
		uuid UUID PRIMARY KEY,
		balldontlie_id INT,
		player_uuid UUID,
		player_bdl_id INT,
		team_uuid UUID,
		team_bdl_id INT,
		game_uuid UUID,
		game_bdl_id INT,
		min TEXT,
		fgm NUMERIC,
		fga NUMERIC,
		fg3m NUMERIC,
		fg3a NUMERIC,
		oreb NUMERIC,
		dreb NUMERIC,
		reb NUMERIC,
		ast NUMERIC,
		stl NUMERIC,
		blk NUMERIC,
		turnovers NUMERIC,
		pf NUMERIC,
		pts NUMERIC,
		fg_pct NUMERIC,
		fg3_pct NUMERIC,
		ft_pct NUMERIC,
		CONSTRAINT fk_players
		FOREIGN KEY(player_uuid)
		REFERENCES players(uuid),
		CONSTRAINT fk_teams
		FOREIGN KEY(team_uuid)
		REFERENCES teams(uuid),
		CONSTRAINT fk_games
		FOREIGN KEY(game_uuid)
		REFERENCES games(uuid)
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
