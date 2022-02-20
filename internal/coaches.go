package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Coach struct {
	UUID        uuid.UUID `db:"uuid"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	IsAssistant bool      `db:"is_assistant"`
	TeamID      uuid.UUID `db:"team_uuid"`
	Season      int       `db:"season"`
	NBA_TeamID  string    `db:"nba_team_id"`
	NBA_ID      string    `db:"nba_id"`
}

const (
	CoachesSchema = `CREATE TABLE coaches(
	uuid UUID PRIMARY KEY,
	first_name TEXT,
	last_name TEXT,
	is_assistant BOOL,
	team_uuid UUID,
	season INT,
	nba_team_id TEXT,
	nba_id TEXT,
	CONSTRAINT fk_teams
	FOREIGN KEY(team_uuid)
	REFERENCES teams(uuid));`
)

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
		// some crummy data here: isAssistant is reversed for a few older years
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
