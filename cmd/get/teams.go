package get

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	internal "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getTeamsCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "teams",
	Run:   getTeamData,
}

func init() {
	GetCmd.AddCommand(getTeamsCmd)
}

func getTeamData(cmd *cobra.Command, args []string) {
	//	loadDefaultVariables()

	err := prepareTeamsSchema()
	if err != nil {
		log.Fatalf("can't create teams schema: %s", err)
	}

	team := internal.Team{}

	teams, err := team.GetAllTeams()
	if err != nil {
		log.Fatalf("can't get teams: %s", err)
	}

	nbaIds, err := internal.GetNbaIds()
	if err != nil {
		log.Fatalf("can't get NBA teamIDs: %s", err)
	}
	teamsWithIds := []internal.Team{}

	for _, team := range teams {
		team.NBA_ID, err = internal.AddNbaId(nbaIds, team)
		if err != nil {
			log.Fatalf("can't add NBA ids to teams: %s", err)
		}
		teamsWithIds = append(teamsWithIds, team)
	}

	if err = insertTeams(teamsWithIds); err != nil {
		log.Printf("Error inserting team: %s", err)
	}

}

func insertTeams(t []internal.Team) error {
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

	_, err = tx.NamedExecContext(ctx, `INSERT INTO teams (
		uuid,
		balldontlie_id,
		nba_id,
		name, 
		abbrev, 
		conference, 
		division) 
		VALUES (
			:uuid,
			:balldontlie_id,
			:nba_id,
			:name,
			:abbrev, 
			:conference, 
			:division)`,
		t)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil

}

func prepareTeamsSchema() error {
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

	schema := `CREATE TABLE teams(
        uuid uuid PRIMARY KEY,
 		balldontlie_id INT UNIQUE,
		nba_id TEXT,
        name TEXT,
		abbrev TEXT,
		conference TEXT,
		division TEXT
		); `

	db.MustExecContext(ctx, schema)

	return nil
}
