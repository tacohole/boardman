package get

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getGamesCmd = &cobra.Command{
	Short: "gets basic info on all games since 1979",
	Long:  "gets teams, date, score, and season stage for all games since 1979",
	Use:   "games",
	Run:   getGames,
}

func init() {
	GetCmd.AddCommand(getGamesCmd)
}

func getGames(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	g := internal.Game{}

	if err := dbutil.PrepareSchema(internal.GameSchema, "nba_data"); err != nil {
		log.Fatalf("could not create games schema: %s", err)
	}

	for i := 1979; i <= 2021; i++ {
		games, err := g.GetSeasonGames(i)
		if err != nil {
			log.Fatalf("can't get games: %s", err)
		}

		if err = insertSeasonGames(games); err != nil {
			log.Fatalf("can't insert games for season %d: %s", i, err)
		}
	}

	if err := cleanupDuplicateGames(); err != nil {
		log.Printf("can't remove duplicate values from games: %s", err)
	}

}

func cleanupDuplicateGames() error {
	db, err := dbutil.DbConn("nba_data")
	if err != nil {
		return err
	}
	defer db.Close()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	tx := db.MustBegin()
	defer tx.Rollback()

	q := `DELETE FROM games a 
		USING games b 
		WHERE a.uuid < b.uuid 
		AND a.balldontlie_id = b.balldontlie_id;`

	_, err = tx.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

func insertSeasonGames(g []internal.Game) error {
	db, err := dbutil.DbConn("nba_data")
	if err != nil {
		return err
	}
	defer db.Close()

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
