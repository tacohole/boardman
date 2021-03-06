package get

import (
	"context"
	"log"

	internal "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"

	"github.com/spf13/cobra"
)

var getPlayersCmd = &cobra.Command{
	Short: "gets all NBA players since 1979 season",
	Long:  "gets player name, position, height/weight, and team for all players since 1979",
	Use:   "players",
	Run:   getPlayers,
}

func init() {
	GetCmd.AddCommand(getPlayersCmd)
}

func getPlayers(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	err := dbutil.PrepareSchema(internal.PlayerSchema, "nba_data")
	if err != nil {
		log.Fatalf("can't create players schema: %s", err)
	}

	p := internal.Player{}

	players, err := p.GetAllPlayers()
	if err != nil {
		log.Fatalf("can't get players: %s", err)
	}

	if err := insertPlayerRows(players); err != nil {
		log.Fatalf("Could not perform insert:, %s", err)
	}
	if verbose {
		log.Printf("Inserted %d players", len(players))
	}

	count, err := updatePlayersWithTeamUUIDs()
	remaining := count - int64(len(players))
	if err != nil || remaining > 0 {
		log.Printf("%d players not updated: %s", remaining, err)
	}

}

// can fail silently with 0 rows updated, so we return count
func updatePlayersWithTeamUUIDs() (int64, error) {
	db, err := dbutil.DbConn("nba_data")
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

	stmt := `UPDATE players
			SET team_uuid = teams.uuid
			FROM teams
			WHERE players.team_bdl_id = teams.balldontlie_id;`

	result := tx.MustExecContext(ctx, stmt)
	tx.Commit()

	return result.RowsAffected()
}

func insertPlayerRows(p []internal.Player) error {
	db, err := dbutil.DbConn("nba_data")
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

	_, err = tx.NamedExecContext(ctx, `INSERT INTO players (
		uuid,
		balldontlie_id,
		first_name,
		last_name,
		position,
		height_feet,
		height_in,
		weight,
		team_bdl_id)
		VALUES (
			:uuid,
			:balldontlie_id,
			:first_name,
			:last_name,
			:position,
			:height_feet,
			:height_in,
			:weight,
			:team_bdl_id);`,
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
