package get

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getCoachesCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "coaches",
	Run:   getCoaches,
}

func getCoaches(cmd *cobra.Command, args []string) {
	// get endpoint http://data.nba.net/prod/v1/{year}/coaches.json
	// make structs in internal
	//
	if err := prepareCoachesSchema(); err != nil {
		log.Fatalf("can't prepare coaches schema, %s", err)
	}
}

func prepareCoachesSchema() error {
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

	schema := `CREATE TABLE coaches(
        uuid UUID PRIMARY KEY,
 		first_name TEXT,
		last_name TEXT,
		is_assistant BOOL,
		team_id uuid,
		nba_team_id TEXT,
		nba_id TEXT
		CONSTRAINT fk_teams
		FOREIGN KEY(team_id)
		REFERENCES teams(uuid)
		); `

	db.MustExecContext(ctx, schema)

	return nil
}
