package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/google/uuid"
	dbutil "github.com/tacohole/boardman/util/db"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Team struct {
	UUID       uuid.UUID `db:"uuid"`
	BDL_ID     int       `json:"id" db:"balldontlie_id"`
	NBA_ID     string    `json:"teamId" db:"nba_id"`
	Name       string    `json:"full_name" db:"name"`
	Abbrev     string    `json:"abbreviation" db:"abbrev"`
	Conference string    `json:"conference" db:"conference"`
	Division   string    `json:"divsion" db:"division"`
}

const (
	TeamSchema = `CREATE TABLE IF NOT EXISTS teams(
	uuid uuid PRIMARY KEY,
	balldontlie_id INT UNIQUE,
	nba_id TEXT,
	name TEXT,
	abbrev TEXT,
	conference TEXT,
	division TEXT
	); `
)

// get all teams
func (t *Team) GetAllTeams() ([]Team, error) {
	allTeams := []Team{}

	getUrl := BDLUrl + BDLTeams

	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var p Page

	err = json.Unmarshal(r, &p)
	if err != nil {
		return nil, err
	}
	for _, d := range p.Data {
		t.UUID = uuid.New()
		t.BDL_ID = d.ID
		t.Abbrev = d.Abbrev
		t.Conference = d.Conference
		t.Name = d.Name
		t.Division = d.Division
		allTeams = append(allTeams, *t)
	}

	return allTeams, nil
}

func GetNbaIds() ([]NbaData, error) {
	getUrl := NbaDataUrl + fmt.Sprint(2021) + Teams

	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
	if err != nil {
		if resp.StatusCode == 429 {
			fmt.Printf("hit a rate limit, nite nite")
			time.Sleep(3000)
			return nil, err
		} else {
			return nil, err
		}
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var page NbaPage
	var teams []NbaData
	var t NbaData
	if err = json.Unmarshal(r, &page); err != nil {
		return nil, err
	}

	for _, item := range page.League.Standard {
		t.TeamID = item.TeamID
		t.Name = item.Name
		t.Abbrev = item.Abbrev
		teams = append(teams, t)
	}
	return teams, nil
}

func AddNbaId(ids []NbaData, team Team) (*string, error) {

	for j := 0; j < len(ids); j++ {
		if team.Abbrev == ids[j].Abbrev {
			return &ids[j].TeamID, nil
		}
	}

	return nil, fmt.Errorf("no NBA ID for team %s", team.Name)
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

	q := `SELECT uuid,balldontlie_id FROM teams`
	dest := &[]Team{}

	if err = tx.SelectContext(ctx, dest, q); err != nil {
		return nil, err
	}

	return *dest, nil
}
