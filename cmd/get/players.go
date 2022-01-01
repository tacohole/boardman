package get

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	schema "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"

	"github.com/spf13/cobra"
)

var getPlayersCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "players",
	Run:   getPlayers,
}

func init() {

	GetCmd.AddCommand(getPlayersCmd)

}

func getPlayers(cmd *cobra.Command, args []string) {
	loadDefaultVariables()
	godotenv.Load(".env")

	p := schema.Player{}

	players, err := p.GetAllPlayers()
	if err != nil {
		log.Fatalf("can't get players: %s", err)
	}

	result, err := insertPlayerRows(players)
	if err != nil {
		log.Printf("Could not perform insert:, %s", err)
	}
	log.Printf("Inserted %s", fmt.Sprint(result))

}

func insertPlayerRows(p []schema.Player) (*sql.Result, error) {
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

	result, err := tx.NamedExecContext(ctx, `INSERT INTO players (
		id, first_name, last_name, balldontlie_id, team_id )
		VALUES (
			:id,
			:first_name,
			:last_name,
			:balldontlie_id,
			:team_id)`,
		p)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &result, nil

}
