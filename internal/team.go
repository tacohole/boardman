package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/google/uuid"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Team struct {
	UUID       uuid.UUID `db:"uuid"`
	BDL_ID     int       `db:"balldontlie_id"`
	NBA_ID     string    `db:"nba_id"`
	Name       string    `db:"name"`
	Abbrev     string    `db:"abbrev"`
	Conference string    `db:"conference"`
	Division   string    `db:"division"`
}

// get team by ID
func (t *Team) GetTeamById() (*Team, error) {
	getUrl := BDLUrl + BDLTeams + fmt.Sprint(t.BDL_ID)

	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(r, &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

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
		return nil, err
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

func AddNbaId(ids []NbaData, team Team) (string, error) {

	for j := 0; j < len(ids); j++ {
		if team.Abbrev == ids[j].Abbrev {
			return ids[j].TeamID, nil
		}
	}

	return "", fmt.Errorf("no NBA ID for team %s", team.Name)
}

// // get all teams in conf - move to Presti, can't query this endpoint
// func (t *Team) GetConfTeams(conf string) ([]Team, error) {
// 	confTeams := []Team{}

// 	getUrl := httpHelpers.BaseUrl + httpHelpers.Teams + "?conference=" + conf

// 	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl, []byte(""), "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	r, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var p Page

// 	err = json.Unmarshal(r, &p)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, d := range p.Data {
// 		t.ID = d.ID
// 		t.Abbrev = d.Abbrev
// 		t.Conference = d.Conference
// 		t.Name = d.Name
// 		t.Division = d.Division
// 		confTeams = append(confTeams, *t)
// 	}

// 	return confTeams, nil
// }

// // get all teams in div - move to Presti, can't query this endpoint
// func (t *Team) GetDivTeams() ([]Team, error) {
// 	divTeams := []Team{}

// 	getUrl := httpHelpers.BaseUrl + httpHelpers.Teams +

// 	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl, []byte(""), "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	r, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var p Page

// 	err = json.Unmarshal(r, &p)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, d := range p.Data {
// 		t.ID = d.ID
// 		t.Abbrev = d.Abbrev
// 		t.Conference = d.Conference
// 		t.Name = d.Name
// 		t.Division = d.Division
// 		divTeams = append(divTeams, *t)
// 	}
// 	return divTeams, nil
// }
