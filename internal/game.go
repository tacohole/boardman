package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Game struct {
	ID           uuid.UUID  `db:"uuid"`
	BDL_ID       int        `json:"id" db:"balldontlie_id"`
	Date         *time.Time `json:"date" db:"date"`
	Home         Team       `json:"home_team"`
	HomeID       int        `db:"home_id"`
	HomeScore    int        `json:"home_team_score" db:"home_score"`
	Visitor      Team       `json:"visitor_team"`
	VisitorID    int        `db:"visitor_id"`
	VisitorScore int        `json:"visitor_team_score" db:"visitor_score"`
	LeagueYear   int        `json:"season" db:"season"`
	IsPostseason bool       `json:"postseason" db:"is_postseason"`
	Winner       Team       `json:"winner" db:"winner_id"`
	Margin       int        `json:"margin" db:"margin"`
}

// get all games for a season
func (g *Game) GetSeasonGames(season int) ([]Game, error) {
	allGames := []Game{}

	var page Page

	for pageIndex := 0; pageIndex <= page.PageData.TotalPages; pageIndex++ {
		getUrl := httpHelpers.BaseUrl + httpHelpers.Games + "?seasons[]=" + fmt.Sprint(season) + "&page=" + fmt.Sprint(pageIndex) + "per_page=100"
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
			g.ID = uuid.New()
			g.BDL_ID = d.ID
			g.Date = d.Date
			g.HomeID = d.Home.ID
			g.HomeScore = d.HomeScore
			g.VisitorID = d.Visitor.ID
			g.VisitorScore = d.VisitorScore
			g.LeagueYear = d.LeagueYear
			g.IsPostseason = d.IsPostseason
			g.CalculateWinnerAndMargin()
			allGames = append(allGames, *g)
		}
	}
	return allGames, nil
}

// calculate winner
func (g *Game) CalculateWinnerAndMargin() {

	if g.HomeScore < g.VisitorScore {
		g.Winner = g.Visitor
		g.Margin = g.VisitorScore - g.HomeScore
	} else {
		g.Winner = g.Home
		g.Margin = g.HomeScore - g.VisitorScore
	}
}
