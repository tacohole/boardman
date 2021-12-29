package get

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	schema "github.com/tacohole/boardman/internal"
	"github.com/tacohole/boardman/util/config"
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

	// insert into database
	for _, player := range players {
		result, err := insertPlayerRow(player)
		if err != nil {
			log.Printf("Could not insert player %d, %s", player.ID, err)
		}
		log.Printf("Inserted %s", fmt.Sprint(result))
	}

}

func insertPlayerRow(p schema.Player) (*sql.Result, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.DbTimeout)
	defer cancel()

	tx := db.MustBegin()

	result, err := tx.NamedExecContext(ctx, "INSERT INTO players (id,first_name,last_name,team) VALUES (:id,:first_name,:last_name,:team)", p)
	if err != nil {
		return nil, err
	}

	return &result, nil

}
