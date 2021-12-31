package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Game struct {
	ID           int        `json:"id"`
	Date         *time.Time `json:"date"`
	Home         Team       `json:"home_team"`
	HomeID       int
	HomeScore    int  `json:"home_team_score"`
	Visitor      Team `json:"visitor_team"`
	VisitorID    int
	VisitorScore int    `json:"visitor_team_score"`
	Season       Season `json:"season"`
	IsPostseason bool   `json:"postseason"`
	Winner       Team   `json:"winner"`
	Margin       int    `json:"margin"`
}

// get all games for a season
func (g *Game) GetSeasonGames(season int) ([]Game, error) {
	allGames := []Game{}

	var page Page

	for pageIndex := 0; pageIndex < page.PageData.NextPageIndex; pageIndex++ {
		getUrl := httpHelpers.BaseUrl + httpHelpers.Games + "?seasons[]=" + fmt.Sprint(season) + "&page=" + fmt.Sprint(pageIndex)
		resp, err := httpHelpers.MakeHttpRequest("GET", getUrl, []byte(""), "")
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
			g.ID = d.ID
			g.Date = d.Date
			g.HomeID = d.Home.ID
			g.HomeScore = d.HomeScore
			g.VisitorID = d.Visitor.ID
			g.VisitorScore = d.VisitorScore
			g.Season = d.Season
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

// "CREATE TABLE games(
// 	game_id INT,
// 	date DATE,
// 	home_team_id INT,
// 	visitor_team_id INT,
// 	home_team_score INT,
// 	visitor_team_score INT,
// 	season INT,
// 	winner_id INT,
// 	margin INT,
// 	is_postseason BOOL,
// 	CONSTRAINT fk_teams
// 	   FOREIGN KEY(home_team_id)
// 	   REFERENCES teams(id),
// 	CONSTRAINT fk_teams
// 	   FOREIGN KEY(visitor_team_id)
// 	   REFERENCES teams(id),
// 	CONSTRAINT fk_teams
// 	   FOREIGN KEY(winner_id)
// 	   REFERENCES teams(id),
// );"
