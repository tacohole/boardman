package internal

import (
	"context"

	"github.com/google/uuid"
	dbutil "github.com/tacohole/boardman/util/db"
)

type TeamSeason struct {
	UUID         uuid.UUID   `db:"uuid"`
	TeamUUID     uuid.UUID   `db:"team_uuid"`
	Season       int         `db:"season"`
	Wins         int         `db:"wins"`
	Losses       int         `db:"losses"`
	WinPct       float32     `db:"wpct"`
	PlusMinus    int         `db:"plus_minus"`
	ConfRank     int         `db:"conf_rank"`
	OvrRank      int         `db:"ovr_rank"`
	MadePlayoffs bool        `db:"made_playoffs"`
	FGM          float32     `db:"fgm"`
	FGA          float32     `db:"fga"`
	FTM          float32     `db:"ftm"`
	FTA          float32     `db:"fta"`
	FG3M         float32     `db:"fg3m"`
	FG3A         float32     `db:"fg3a"`
	OREB         float32     `db:"oreb"`
	DREB         float32     `db:"dreb"`
	REB          float32     `db:"reb"`
	AST          float32     `db:"ast"`
	STL          float32     `db:"stl"`
	BLK          float32     `db:"blk"`
	TO           float32     `db:"turnovers"`
	PF           float32     `db:"pf"`
	PTS          float32     `db:"pts"`
	FG_PCT       float32     `db:"fg_pct"`
	FG3_PCT      float32     `db:"fg3_pct"`
	FT_PCT       float32     `db:"ft_pct"`
	Roster       []uuid.UUID `db:"roster"`
	Coaches      []uuid.UUID `db:"coaches"`
}

func PrepareTeamSeasonSchema() error {
	schema := `CREATE TABLE IF NOT EXISTS team_season(
	uuid UUID
	team_uuid UUID
	season INT
	wins INT
	losses INT
	wpct NUMERIC
	plus_minus INT
	conf_rank INT
	ovr_rank INT
	made_playoffs BOOL
	fgm NUMERIC(4)
	fga NUMERIC(4)
	ftm NUMERIC(4)
	fta NUMERIC(4)
	fg3m NUMERIC(4)
	fg3a NUMERIC(4)
	oreb NUMERIC(4)
	dreb NUMERIC(4)
	reb NUMERIC(4)
	ast NUMERIC(4)
	stl NUMERIC(4)
	blk NUMERIC(4)
	turnovers NUMERIC(4)
	pf NUMERIC(4)
	pts NUMERIC(4)
	fg_pct NUMERIC(4)
	fg3_pct NUMERIC(4)
	ft_pct NUMERIC(4)
	roster UUID[]
	coaches UUID[]
	CONSTRAINT fk_teams
	FOREIGN KEY(team_uuid)
	REFERENCES teams(uuid));`

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

// sum wins for all teams
// sum losses for all teams
// sum wpct for all teams
// calculate +/-: sum win margin, sum loss margin, subtract win from loss
// get conf rank - sort by wpct within conference
// get conf champ - get last game won within conf
// get league champ - get last game won in season
// get is_postseason for each team/year
// sum stats by team-player, /82
// build roster
// build coaches
