package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Game struct {
	ID           uuid.UUID `db:"uuid"`
	BDL_ID       int       `json:"id" db:"balldontlie_id"`
	Date         string    `json:"date" db:"date"`
	HomeID       uuid.UUID `db:"home_id"`
	HomeScore    int       `json:"home_team_score" db:"home_score"`
	VisitorID    uuid.UUID `db:"visitor_id"`
	VisitorScore int       `json:"visitor_team_score" db:"visitor_score"`
	Season       int       `json:"season" db:"season"`
	IsPostseason bool      `json:"postseason" db:"is_postseason"`
	Winner       uuid.UUID `db:"winner_id"`
	Margin       int       `db:"margin"`
}

const (
	GameSchema = `CREATE TABLE IF NOT EXISTS games(
	uuid uuid PRIMARY KEY,
	balldontlie_id INT,
	date DATE,
	home_id UUID,
	visitor_id UUID,
	home_score INT,
	visitor_score INT,
	season INT,
	winner_id UUID,
	margin INT,
	is_postseason BOOL,
	CONSTRAINT fk_teams
	   FOREIGN KEY(home_id)
	   REFERENCES teams(uuid),
	   FOREIGN KEY(visitor_id)
	   REFERENCES teams(uuid),
	   FOREIGN KEY(winner_id)
	   REFERENCES teams(uuid)
	); `
)

// get all games for a season
func (g *Game) GetSeasonGames(season int) ([]Game, error) {
	allGames := []Game{}

	teamCache, err := GetTeamCache()
	if err != nil {
		return nil, err
	}
	fmt.Print(teamCache)

	var page Page
	var errorCount int

	for pageIndex := 0; pageIndex <= page.PageData.TotalPages; pageIndex++ {
		getUrl := BDLUrl + BDLGames + "?seasons[]=" + fmt.Sprint(season) + "&amp;page=" + fmt.Sprint(pageIndex) + "per_page=100"
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

		r, err := io.ReadAll(resp.Body)
		if err != nil {
			errorCount++
			if errorCount > 2 {
				return nil, err
			} else {
				pageIndex--
				continue
			}
		}

		err = json.Unmarshal(r, &page)
		if err != nil {
			return nil, err
		}
		for _, d := range page.Data {
			homeId, err := AddTeamUUID(d.Home.BDL_ID, teamCache)
			if err != nil {
				log.Print(d.Home)
				return nil, err
			}
			visitorId, err := AddTeamUUID(d.Visitor.BDL_ID, teamCache)
			if err != nil {
				log.Print(d.Visitor)
				return nil, err
			}

			g.ID = uuid.New()
			g.BDL_ID = d.ID
			g.Date = d.Date
			g.HomeID = *homeId
			g.HomeScore = d.HomeScore
			g.VisitorID = *visitorId
			g.VisitorScore = d.VisitorScore
			g.Season = d.Season
			g.IsPostseason = d.IsPostseason
			g.CalculateWinnerAndMargin()
			allGames = append(allGames, *g)
		}
		time.Sleep(1000 * time.Millisecond) // avoiding 429 from BDL
	}
	return allGames, nil
}

// calculate winner
func (g *Game) CalculateWinnerAndMargin() {

	if g.HomeScore < g.VisitorScore {
		g.Winner = g.VisitorID
		g.Margin = g.VisitorScore - g.HomeScore
	} else {
		g.Winner = g.HomeID
		g.Margin = g.HomeScore - g.VisitorScore
	}
}

// this is a different method for syncing the table's UUID but we're leaving it
// because we are not storing the BDL team ID in this table
func AddTeamUUID(bdlId int, teamCache []Team) (*uuid.UUID, error) {
	for j := 0; j < len(teamCache); j++ {
		if bdlId == teamCache[j].BDL_ID {
			return &teamCache[j].UUID, nil
		}
	}

	return nil, fmt.Errorf("no team UUID for team %d", bdlId)
}
