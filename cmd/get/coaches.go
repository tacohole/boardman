package get

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getCoachesCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "coaches",
	Run:   getCoaches,
}

func init() {
	GetCmd.AddCommand(getCoachesCmd)
}

func getCoaches(cmd *cobra.Command, args []string) {

	if err := prepareCoachesSchema(); err != nil {
		log.Fatalf("can't prepare coaches schema, %s", err)
	}

	for i := 2015; i <= 2021; i++ {
		coaches, err := internal.GetSeasonCoaches(i)
		if err != nil {
			log.Fatalf("can't get coaches for season %d: %s", i, err)
		}

		if err = insertCoaches(coaches); err != nil {
			log.Printf("can't insert coaches for season %d: %s", i, err)
		}
	}

	// silently erroring still
	coachesCount, err := updateCoachesWithTeamIds()
	if err != nil || coachesCount < 1 {
		log.Fatalf("can't add team uuids to coaches table: %s", err)
	}

}

func updateCoachesWithTeamIds() (int64, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return 0, err
	}

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	tx := db.MustBegin()
	defer tx.Rollback()

	stmt := `UPDATE coaches 
			SET team_uuid = teams.uuid
			FROM teams 
			WHERE coaches.nba_team_id = teams.nba_id;`

	result := tx.MustExecContext(ctx, stmt)

	return result.RowsAffected()
}

func insertCoaches(c []internal.Coach) error {
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

	result, err := tx.NamedExecContext(ctx, `INSERT INTO coaches (
		uuid,
		first_name,
		last_name,
		is_assistant, 
		nba_id,
		season,
		nba_team_id 
	) VALUES (
		:uuid,
		:first_name,
		:last_name,
		:is_assistant, 
		:nba_id,
		:season,
		:nba_team_id )`,
		c)
	if err != nil {
		log.Printf("Insert failed, %s", result)
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
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
		team_uuid UUID,
		season INT,
		nba_team_id TEXT,
		nba_id TEXT,
		CONSTRAINT fk_teams
		FOREIGN KEY(team_uuid)
		REFERENCES teams(uuid));`

	db.MustExecContext(ctx, schema)

	return nil
}
