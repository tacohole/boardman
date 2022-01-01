package get

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getGamesCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "games",
	Run:   getGames,
}

func init() {
	GetCmd.AddCommand(getGamesCmd)
}

func getGames(cmd *cobra.Command, args []string) {
	//loadDefaultVariables()

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

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	tx := db.MustBegin()
	defer tx.Rollback()

	result, err := tx.NamedExecContext(ctx, `INSERT INTO games (
		id, 
		date,
		home_id, 
		home_team_score, 
		visitor_id, 
		visitor_team_score, 
		season, 
		postseason, 
		winner, 
		margin, 
	) VALUES (
		:date,
		:home_id, 
		:home_team_score, 
		:visitor_id, 
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
