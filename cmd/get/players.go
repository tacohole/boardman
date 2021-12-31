package get

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	schema "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"

	"github.com/spf13/cobra"
)

var getPlayersCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getPlayers,
}

func init() {
	getPlayersCmd.Flags().StringVar(&writeTo, "output", "", "output type, options are JSON or csv")

	getCmd.AddCommand(getPlayersCmd)

}

func getPlayers(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

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
		id, 
		first_name, 
		last_name, 
		team, ) 
		VALUES (
			:id,
			:first_name,
			:last_name,
			:team)`,
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
