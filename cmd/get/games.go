package get

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
	"github.com/tacohole/boardman/util/config"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getGamesCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getGames,
}

func init() {
	getGamesCmd.Flags().StringVar(&writeTo, "output", "", "output type, options are JSON or csv")

	getGamesCmd.AddCommand(getPlayersCmd)

}

func getGames(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	g := schema.Game{}
	seasons := []int{1981: 2021}

	for _, season := range seasons {
		games, err := g.GetSeasonGames(season)
		if err != nil {
			log.Fatalf("can't get games: %s", err)
		}

		result, err := insertSeasonGames(games)
		if err != nil {
			log.Printf("can't get games for season %d: %s", season, err)
		}
		log.Print(fmt.Sprint(result))
	}

}

func insertSeasonGames(g []schema.Game) (*sql.Result, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.DbTimeout)
	defer cancel()

	tx := db.MustBegin()
	defer tx.Rollback()

	result, err := tx.NamedExecContext(ctx, `INSERT INTO games (
		id, 
			date,
			home_team, 
			home_team_score, 
			visitor_team, 
			visitor_team_score, 
			season, 
			postseason, 
			winner, 
			margin, 
		}
	) VALUES (
		:date,
		:home_team, 
		:home_team_score, 
		:visitor_team, 
		:visitor_team_score, 
		:season, 
		:postseason, 
		:winner, 
		:margin, )`,
		g)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &result, nil

}
