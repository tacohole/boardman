package get

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
	"github.com/tacohole/boardman/util/config"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getTeamsCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getTeamData,
}

func init() {
	getTeamsCmd.Flags().StringVar(&writeTo, "output", "", "output type")

	getTeamsCmd.MarkFlagRequired("writeTo")

}

func getTeamData(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	team := schema.Team{}

	teams, err := team.GetAllTeams()
	if err != nil {
		log.Fatalf("can't get teams: %s", err)
	}

	// insert into database
	for _, team := range teams {
		result, err := insertTeamRow(team)
		if err != nil {
			log.Printf("Error inserting team: %s", err)
			continue
		}
		log.Printf("Inserted %s", fmt.Sprint(result))
	}

}

func insertTeamRow(t schema.Team) (*sql.Result, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.DbTimeout)
	defer cancel()

	tx := db.MustBegin()

	result, err := tx.NamedExecContext(ctx, "INSERT INTO teams (id,name,abbrev,conference,division) VALUES (:id,:name,:abbrev,:conference,:division)", t)
	if err != nil {
		return nil, err
	}

	return &result, nil

}
