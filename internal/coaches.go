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

type Coach struct {
	UUID        uuid.UUID `db:"uuid"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	IsAssistant bool      `db:"is_assistant"`
	TeamID      uuid.UUID `db:"team_id"`
	Season      int       `db:"season"`
	NBA_TeamID  string    `db:"nba_team_id"`
	NBA_ID      string    `db:"nba_id"`
}

func GetSeasonCoaches(season int) ([]Coach, error) {
	getUrl := NbaDataUrl + fmt.Sprint(season) + Coaches

	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var page NbaPage
	var c Coach
	var coaches []Coach
	if err = json.Unmarshal(r, &page); err != nil {
		return nil, err
	}

	for _, item := range page.League.Standard {
		c.UUID = uuid.New()
		c.FirstName = item.FirstName
		c.LastName = item.LastName
		c.IsAssistant = item.IsAssistant
		// problem specific to the data from this endpoint: isAssistant is reversed for a few older years
		if season <= 2017 {
			c.IsAssistant = !c.IsAssistant
		}
		c.NBA_ID = item.PersonID
		c.NBA_TeamID = item.TeamID
		c.Season = season
		coaches = append(coaches, c)
	}
	return coaches, nil
}

func AddTeamUUID(teams []Team, coach Coach) (*uuid.UUID, error) {

	for j := 0; j < len(teams); j++ {
		if coach.NBA_TeamID == teams[j].NBA_ID {
			return &teams[j].UUID, nil
		}
	}

	return nil, fmt.Errorf("no team UUID for coach %s", coach.NBA_ID)
}

func GetTeamCache() ([]Team, error) {
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

	q := `SELECT * FROM teams`
	dest := &[]Team{}

	if err = tx.SelectContext(ctx, dest, q); err != nil {
		return nil, err
	}

	return *dest, nil
}
