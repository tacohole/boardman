package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Player struct {
	UUID       uuid.UUID `db:"uuid"`
	BDL_ID     int       `json:"id" db:"balldontlie_id"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	Position   string    `json:"position" db:"position"`
	HeightFeet int       `json:"height_feet" db:"height_feet"`
	HeightIn   int       `json:"height_inches" db:"height_in"`
	Weight     int       `json:"weight_pounds" db:"weight"`
	TeamUUID   uuid.UUID `db:"team_uuid"`
	TeamBDL_ID int       `json:"team" db:"team_bdl_id"`
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
		time.Sleep(1000 * time.Millisecond)
	}
	return allPlayers, nil
}
