package get

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	internal "github.com/tacohole/boardman/internal"
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
	// loadDefaultVariables()

	err := preparePlayersSchema()
	if err != nil {
		log.Fatalf("can't create players schema: %s", err)
	}

	p := internal.Player{}

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

func insertPlayerRows(p []internal.Player) (*sql.Result, error) {
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

func preparePlayersSchema() error {
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

	schema := `CREATE TABLE players(
        uuid uuid PRIMARY KEY,
 		balldontlie_id INT,
        first_name TEXT,
		last_name TEXT,
		current_team_id INT,
        CONSTRAINT fk_teams
           FOREIGN KEY(current_team_id)
           REFERENCES teams(id)
		); `

	db.MustExecContext(ctx, schema)

	return nil
}
