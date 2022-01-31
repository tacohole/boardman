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
	PSWins       int         `db:"postseason_wins"`
	PSLosses     int         `db:"postseason_losses"`
	WinPct       float32     `db:"wpct"`
	PlusMinus    int         `db:"plus_minus"`
	ConfRank     int         `db:"conf_rank"`
	OvrRank      int         `db:"ovr_rank"`
	MadePlayoffs bool        `db:"made_playoffs"`
	FGM          float32     `db:"fgm"`       // for player in team in season sum fgm on all games
	FGA          float32     `db:"fga"`       // for player in team in season sum fga on all games
	FTM          float32     `db:"ftm"`       // for player in team in season sum ftm on all games
	FTA          float32     `db:"fta"`       // for player in team in season sum fta on all games
	FG3M         float32     `db:"fg3m"`      // for player in team in season sum fg3m on all games
	FG3A         float32     `db:"fg3a"`      // for player in team in season sum fg3a on all games
	OREB         float32     `db:"oreb"`      // for player in team in season sum oreb on all games
	DREB         float32     `db:"dreb"`      // for player in team in season sum dreb on all games
	REB          float32     `db:"reb"`       // for player in team in season sum reb on all games
	AST          float32     `db:"ast"`       // for player in team in season sum ast on all games
	STL          float32     `db:"stl"`       // for player in team in season sum stl on all games
	BLK          float32     `db:"blk"`       // for player in team in season sum blk on all games
	TO           float32     `db:"turnovers"` // for player in team in season sum turnovers on all games
	PF           float32     `db:"pf"`        // for player in team in season sum pf on all games
	PTS          float32     `db:"pts"`       // for player in team in season sum pts on all games
	FG_PCT       float32     `db:"fg_pct"`    // for team in season fgm / (fga+fgm)
	FG3_PCT      float32     `db:"fg3_pct"`   // for team in season fg3m / (fg3a+fg3m)
	FT_PCT       float32     `db:"ft_pct"`    // for team in season ftm / (fta + ftm)
	Roster       []uuid.UUID `db:"roster"`    // select player_uuid from gamestats where season = season and teamuuid = team.uuid
	Coaches      []uuid.UUID `db:"coaches"`   // select coach.uuid from coaches where season = season and team_uuid = team.uuid
}

func PrepareTeamSeasonSchema() error {
	schema := `CREATE TABLE IF NOT EXISTS team_season(
	uuid UUID
	team_uuid UUID
	season INT
	wins INT
	losses INT
	postseason_wins INT
	postseason_losses INT
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

// get conf champ - get last game won within conf
// get league champ - get last game won in season
// sum stats by team-player, /82
// build roster
// build coaches
