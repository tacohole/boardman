package get

import (
	"context"
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

	if err := insertPlayerRows(players); err != nil {
		log.Printf("Could not perform insert:, %s", err)
	}
	log.Printf("Inserted %d players", len(players))

}

func insertPlayerRows(p []internal.Player) error {
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

	_, err = tx.NamedExecContext(ctx, `INSERT INTO players (
		uuid,
		first_name,
		last_name,
		balldontlie_id, 
		team_id )
		VALUES (
			:uuid,
			:first_name,
			:last_name,
			:balldontlie_id,
			:team_id)`,
		p)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

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
		team_id uuid,
        CONSTRAINT fk_teams
        FOREIGN KEY(team_id)
        REFERENCES teams(uuid)
		);`

	db.MustExecContext(ctx, schema)

	return nil
}
