package get

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getTeamsCmd = &cobra.Command{
	Short: "test",
	Long:  "",
	Use:   "teams",
	Run:   getTeamData,
}

func init() {
	GetCmd.AddCommand(getTeamsCmd)
}

func getTeamData(cmd *cobra.Command, args []string) {
	loadDefaultVariables()
	godotenv.Load(".env")

	team := schema.Team{}

	teams, err := team.GetAllTeams()
	if err != nil {
		log.Fatalf("can't get teams: %s", err)
	}

	result, err := insertTeams(teams)
	if err != nil {
		log.Printf("Error inserting team: %s", err)
	}
	log.Printf("Inserted %s", fmt.Sprint(result))

}

func insertTeams(t []schema.Team) (*sql.Result, error) {
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

	result, err := tx.NamedExecContext(ctx, `INSERT INTO teams (
		id, 
		name, 
		abbrev, 
		conference, 
		division) 
		VALUES (
			:id, 
			:name,
			:abbrev, 
			:conference, 
			:division)`,
		t)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return &result, nil

}
