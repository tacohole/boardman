package internal

import (
	"github.com/google/uuid"
)

type PlayerSeason struct {
	PlayerUUID  uuid.UUID `db:"uuid"`
	BDL_ID      int       `db:"balldontlie_id"`
	LeagueYear  int       `db:"season"`
	GamesPlayed int       `db:"games_played"`
	Minutes     string    `db:"avg_min"`
	FGM         float32   `db:"fgm"`
	FGA         float32   `db:"fga"`
	FG3M        float32   `db:"fg3m"`
	FG3A        float32   `db:"fg3a"`
	OREB        float32   `db:"oreb"`
	DREB        float32   `db:"dreb"`
	REB         float32   `db:"reb"`
	AST         float32   `db:"ast"`
	STL         float32   `db:"stl"`
	BLK         float32   `db:"blk"`
	TO          float32   `db:"turnovers"`
	PF          float32   `db:"pf"`
	PTS         float32   `db:"pts"`
	FG_PCT      float32   `db:"fg_pct"`
	FG3_PCT     float32   `db:"fg3_pct"`
	FT_PCT      float32   `db:"ft_pct"`
}

const (
	PlayerSeasonSchema = `CREATE TABLE IF NOT EXISTS player_season(
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
)
