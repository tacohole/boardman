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
	loadDefaultVariables()
	if err := prepareCoachesSchema(); err != nil {
		log.Fatalf("can't prepare coaches schema, %s", err)
	}

	coachCount := 0

	for i := 2015; i <= 2021; i++ {
		coaches, err := internal.GetSeasonCoaches(i)
		if err != nil {
			log.Fatalf("can't get coaches for season %d: %s", i, err)
		}

		if err = insertCoaches(coaches); err != nil {
			log.Printf("can't insert coaches for season %d: %s", i, err)
		}
		coachCount += len(coaches)
	}

	if verbose {
		log.Printf("inserted %d coaches", coachCount)
	}

	count, err := updateCoachesWithTeamIds()
	remaining := int64(coachCount) - count
	if err != nil || remaining > 1 {
		log.Printf("%d coaches not updated: %s", remaining, err)
	} else if verbose {
		log.Printf("%d coaches updated", remaining)
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
	tx.Commit()

	return result.RowsAffected()
}

func insertCoaches(c []internal.Coach) error {
	db, err := dbutil.DbConn()
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
