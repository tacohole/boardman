package internal

import (
	"encoding/json"
	"fmt"
	"io"

	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Player struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CurrentTeam Team   `json:"team"`
}

// get player by ID
func (p *Player) GetPlayerById(id string) (*Player, error) {
	getUrl := httpHelpers.BaseUrl + httpHelpers.Players + fmt.Sprint(p.ID)

	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(r, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// get all players
func (p *Player) GetAllPlayers() ([]Player, error) {
	allPlayers := []Player{}

	var page Page

	for pageIndex := 0; pageIndex < page.PageData.NextPageIndex; pageIndex++ {
		getUrl := httpHelpers.BaseUrl + httpHelpers.Players + "?page=" + fmt.Sprint(pageIndex)
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
			p.ID = d.ID
			p.FirstName = d.FirstName
			p.LastName = d.LastName
			p.CurrentTeam = d.CurrentTeam
			allPlayers = append(allPlayers, *p)
		}
	}
	return allPlayers, nil
}
