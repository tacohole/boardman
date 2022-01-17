package get

import (
	"context"
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
	// loadDefaultVariables()

	g := schema.Game{}

	err := schema.PrepareSeasonSchema()
	if err != nil {
		log.Fatalf("could not create games schema: %s", err)
	}

	for i := 1979; i < 2021; i++ {
		games, err := g.GetSeasonGames(i)
		if err != nil {
			log.Fatalf("can't get games: %s", err)
		}

		if err = insertSeasonGames(games); err != nil {
			log.Printf("can't insert games for season %d: %s", i, err)
		}
	}

}

func insertSeasonGames(g []schema.Game) error {
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

	tx := db.MustBegin()
	defer tx.Rollback()

	_, err = tx.NamedExecContext(ctx, `INSERT INTO games (
		uuid,
		balldontlie_id,
		date,
		home_id, 
		home_score, 
		visitor_id, 
		visitor_score, 
		season, 
		is_postseason, 
		winner_id, 
		margin 
	) VALUES (
		:uuid,
		:balldontlie_id,
		:date,
		:home_id, 
		:home_score, 
		:visitor_id, 
		:visitor_score, 
		:season, 
		:is_postseason, 
		:winner_id, 
		:margin )`,
		g)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}
