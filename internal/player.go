package internal

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"

	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Player struct {
	UUID      uuid.UUID `db:"uuid"`
	BDL_ID    int       `json:"id" db:"balldontlie_id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Position  string    `json:"position" db:"position"`
	HeightFt  int       `json: "height_feet" db:"height_feet"`
	HeightIn  int       `json:"height_inches" db:"height_in"`
	Weight    int       `json:"weight_pounds" db:"weight"`
	TeamID    int       `json:"team" db:"team_id"`
}

// get player by ID
func GetPlayerById(id int) (*Player, error) {
	getUrl := httpHelpers.BaseUrl + httpHelpers.Players + fmt.Sprint(id)

	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	p := Player{}

	err = json.Unmarshal(r, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// get all players
func (p *Player) GetAllPlayers() ([]Player, error) {
	allPlayers := []Player{}

	var page Page

	for pageIndex := 0; pageIndex <= page.PageData.TotalPages; pageIndex++ {
		getUrl := httpHelpers.BaseUrl + httpHelpers.Players + "/?page=" + fmt.Sprint(pageIndex) + "&per_page=100"
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
			p.FirstName = d.FirstName
			p.LastName = d.LastName
			p.BDL_ID = d.ID
			p.TeamID = d.Team.ID
			allPlayers = append(allPlayers, *p)
		}
	}
	return allPlayers, nil
}
