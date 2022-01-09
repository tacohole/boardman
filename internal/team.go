package internal

import (
	"encoding/json"
	"fmt"
	"io"

	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Team struct {
	ID         int    `json:"id"`
	Name       string `json:"full_name"`
	Abbrev     string `json:"abbreviation"`
	Conference string `json:"conference"`
	Division   string `json:"division"`
}

// get team by ID
func (t *Team) GetTeamById() (*Team, error) {
	getUrl := BDLUrl + BDLTeams + fmt.Sprint(t.ID)

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
		t.ID = d.ID
		t.Abbrev = d.Abbrev
		t.Conference = d.Conference
		t.Name = d.Name
		t.Division = d.Division
		allTeams = append(allTeams, *t)
	}

	return allTeams, nil
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
