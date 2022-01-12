package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	dbutil "github.com/tacohole/boardman/util/db"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

type SingleGame struct {
	PlayerUUID uuid.UUID `db:"player_uuid"`
	GameID     uuid.UUID `db:"game_uuid"`
	TeamID     int       `db:"team_id"`
	Minutes    string    `json:"min" db:"min"`
	FGM        float32   `json:"fgm" db:"fgm"`
	FGA        float32   `json:"fga" db:"fga"`
	FG3M       float32   `json:"fg3m" db:"fg3m"`
	FG3A       float32   `json:"fg3a" db:"fg3a"`
	OREB       float32   `json:"oreb" db:"oreb"`
	DREB       float32   `json:"dreb" db:"dreb"`
	REB        float32   `json:"reb" db:"reb"`
	AST        float32   `json:"ast" db:"ast"`
	STL        float32   `json:"stl" db:"stl"`
	BLK        float32   `json:"blk" db:"blk"`
	TO         float32   `json:"turnover" db:"to"`
	PF         float32   `json:"pf" db:"pf"`
	PTS        float32   `json:"pts" db:"pts"`
	FG_PCT     float32   `json:"fg_pct" db:"fg_pct"`
	FG3_PCT    float32   `json:"fg3_pct" db:"fg3_pct"`
	FT_PCT     float32   `json:"ft_pct" db:"ft_pct"`
}

func (s *SingleGame) GetAllGameStats(season int) ([]SingleGame, error) {
	var games []SingleGame // init return value
	var page Page

	for pageIndex := 0; pageIndex <= page.PageData.TotalPages; pageIndex++ {
		getUrl := BDLUrl + BDLStats + "?seasons[]=" + fmt.Sprint(season)

		resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(r, &page)

		for _, d := range page.Data {
			playerId, err := GetUUIDFromBDLID(d.Player.BDL_ID)
			if err != nil {
				return nil, err
			}

			gameId, err := GetUUIDFromBDLID(d.Game.BDL_ID)
			if err != nil {
				return nil, err
			}

			s.AST = d.AST
			s.BLK = d.BLK
			s.DREB = d.DREB
			s.FG3A = d.FG3A
			s.FG3M = d.FG3M
			s.FG3_PCT = d.FG3_PCT
			s.FGA = d.FGA
			s.FGM = d.FGM
			s.FT_PCT = d.FT_PCT
			s.GameID = *gameId
			s.Minutes = d.Minutes
			s.OREB = d.OREB
			s.PF = d.OREB
			s.PF = d.PF
			s.PTS = d.PTS
			s.PlayerUUID = *playerId
			s.REB = d.REB
			s.STL = d.STL
			s.TO = d.TO
			s.TeamID = d.Team.BDL_ID
		}
	}

	return games, nil
}

func PrepareGameStatsSchema() error {
	schema := `CREATE TABLE player_game_stats(
		player_uuid UUID,
		game_uuid UUID,
		team_id INT,
		min NUMERIC(4,2),
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
		FOREIGN KEY(player_uuid)
		REFERENCES players(uuid),
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
