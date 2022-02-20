package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	dbutil "github.com/tacohole/boardman/util/db"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Player struct {
	UUID       uuid.UUID `db:"uuid"`
	BDL_ID     int       `db:"balldontlie_id"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	Position   string    `db:"position"`
	HeightFeet int       `db:"height_feet"`
	HeightIn   int       `db:"height_in"`
	Weight     int       `db:"weight"`
	TeamUUID   uuid.UUID `db:"team_uuid"`
	TeamBDL_ID int       `db:"team_bdl_id"`
}

const (
	PlayerSchema = `CREATE TABLE IF NOT EXISTS players(
	uuid uuid PRIMARY KEY,
	balldontlie_id INT,
	first_name TEXT,
	last_name TEXT,
	position TEXT,
	height_feet NUMERIC,
	height_in NUMERIC,
	weight NUMERIC,
	team_uuid uuid,
	team_bdl_id INT,
	CONSTRAINT fk_teams
	FOREIGN KEY(team_uuid)
	REFERENCES teams(uuid)
	);`
)

func GetPlayerIdCache() ([]Player, error) {
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

	tx := db.MustBegin()
	defer tx.Rollback()

	q := `SELECT uuid FROM players`
	dest := &[]Player{}

	if err = tx.SelectContext(ctx, dest, q); err != nil {
		return nil, err
	}

	return *dest, nil
}

// get all players
func (p *Player) GetAllPlayers() ([]Player, error) {
	allPlayers := []Player{}

	var page Page

	for pageIndex := 0; pageIndex <= page.PageData.TotalPages; pageIndex++ {
		getUrl := BDLUrl + BDLPlayers + "/?page=" + fmt.Sprint(pageIndex) + "&per_page=100"
		resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		r, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(r, &page)
		if err != nil {
			return nil, err
		}
		for _, d := range page.Data {
			p.UUID = uuid.New()
			p.BDL_ID = d.ID
			p.FirstName = d.FirstName
			p.LastName = d.LastName
			p.Position = d.Position
			p.HeightFeet = d.HeightFeet
			p.HeightIn = d.HeightIn
			p.Weight = d.Weight
			p.TeamBDL_ID = d.Team.BDL_ID
			allPlayers = append(allPlayers, *p)
		}
		// avoiding a 429
		time.Sleep(1000 * time.Millisecond)
	}
	return allPlayers, nil
}
