package internal

import (
	"github.com/google/uuid"
)

type SingleGame struct {
	UUID         uuid.UUID `db:"uuid"`
	BDL_ID       int       `db:"balldontlie_id"`
	GameUUID     uuid.UUID `db:"game_uuid"`
	GameBDL_ID   int       `db:"game_bdl_id"`
	PlayerUUID   uuid.UUID `db:"player_uuid"`
	PlayerBDL_ID int       `db:"player_bdl_id"`
	TeamUUID     uuid.UUID `db:"team_uuid"`
	TeamBDL_ID   int       `db:"team_bdl_id"`
	Season       int       `db:"season"`
	Minutes      string    `db:"min"`
	FGM          float32   `db:"fgm"`
	FGA          float32   `db:"fga"`
	FTM          float32   `db:"ftm"`
	FTA          float32   `db:"fta"`
	FG3M         float32   `db:"fg3m"`
	FG3A         float32   `db:"fg3a"`
	OREB         float32   `db:"oreb"`
	DREB         float32   `db:"dreb"`
	REB          float32   `db:"reb"`
	AST          float32   `db:"ast"`
	STL          float32   `db:"stl"`
	BLK          float32   `db:"blk"`
	TO           float32   `db:"turnovers"`
	PF           float32   `db:"pf"`
	PTS          float32   `db:"pts"`
	FG_PCT       float32   `db:"fg_pct"`
	FG3_PCT      float32   `db:"fg3_pct"`
	FT_PCT       float32   `db:"ft_pct"`
}

const (
	GameStatsSchema = `CREATE TABLE IF NOT EXISTS player_game_stats(
		uuid UUID PRIMARY KEY,
		balldontlie_id INT UNIQUE,
		player_uuid UUID,
		player_bdl_id INT,
		team_uuid UUID,
		team_bdl_id INT,
		game_uuid UUID,
		game_bdl_id INT,
		season INT,
		min TEXT,
		fgm NUMERIC(4),
		fga NUMERIC(4),
		ftm NUMERIC(4),
		fta NUMERIC(4),
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
		FOREIGN KEY(player_uuid)
		REFERENCES players(uuid),
		CONSTRAINT fk_teams
		FOREIGN KEY(team_uuid)
		REFERENCES teams(uuid),
		CONSTRAINT fk_games
		FOREIGN KEY(game_uuid)
		REFERENCES games(uuid)
	);`
)
