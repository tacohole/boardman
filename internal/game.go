package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	dbutil "github.com/tacohole/boardman/util/db"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

type Game struct {
	ID           uuid.UUID `db:"uuid"`
	BDL_ID       int       `json:"id" db:"balldontlie_id"`
	Date         string    `json:"date" db:"date"`
	HomeID       int       `db:"home_id"`
	HomeScore    int       `json:"home_team_score" db:"home_score"`
	VisitorID    int       `db:"visitor_id"`
	VisitorScore int       `json:"visitor_team_score" db:"visitor_score"`
	LeagueYear   int       `json:"season" db:"season"`
	IsPostseason bool      `json:"postseason" db:"is_postseason"`
	Winner       int       `json:"winner" db:"winner_id"`
	Margin       int       `json:"margin" db:"margin"`
}

// get all games for a season
func (g *Game) GetSeasonGames(season int) ([]Game, error) {
	allGames := []Game{}

	var page Page

	for pageIndex := 0; pageIndex <= page.PageData.TotalPages; pageIndex++ {
		getUrl := BDLUrl + BDLGames + "?seasons[]=" + fmt.Sprint(season) + "&amp;page=" + fmt.Sprint(pageIndex) + "per_page=100"
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
			g.HomeID = d.Home.BDL_ID
			g.HomeScore = d.HomeScore
			g.VisitorID = d.Visitor.BDL_ID
			g.VisitorScore = d.VisitorScore
			g.LeagueYear = d.LeagueYear
			g.IsPostseason = d.IsPostseason
			g.CalculateWinnerAndMargin()
			allGames = append(allGames, *g)
		}
		time.Sleep(2000 * time.Millisecond) // avoiding 429 from BDL
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

func PrepareSeasonSchema() error {
	db, err := dbutil.DbConn()
	if err != nil {
		return err
	}

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	schema := `CREATE TABLE games(
        uuid uuid PRIMARY KEY,
 		balldontlie_id INT,
        date DATE,
        home_id INT,
        visitor_id INT,
        home_score INT,
        visitor_score INT,
        season INT,
        winner_id INT,
        margin INT,
        is_postseason BOOL,
        CONSTRAINT fk_teams
           FOREIGN KEY(home_id)
           REFERENCES teams(balldontlie_id),
           FOREIGN KEY(visitor_id)
           REFERENCES teams(balldontlie_id),
           FOREIGN KEY(winner_id)
           REFERENCES teams(balldontlie_id)
		); `

	db.MustExecContext(ctx, schema)

	return nil
}
