package get

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getStatsCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "stats",
	Run:   getStats,
}

func init() {
	GetCmd.AddCommand(getStatsCmd)
}

func getStats(cmd *cobra.Command, args []string) {

	err := prepareStatsSchema()
	if err != nil {
		log.Fatalf("can't create schema for stats: %s", err)
	}

	stats, err := getPlayerSeasonAverages(2021, 1)
	if err != nil {
		log.Printf("can't get stats for ")
	}
	err = insertPlayerSeasonAverages(stats)

}

func insertPlayerSeasonAverages(stats *schema.PlayerYear) error {

	return nil
}

func getPlayerSeasonAverages(season int, playerID int) (*schema.PlayerYear, error) {

	return nil, nil
}

func prepareStatsSchema() error {
	schema := `CREATE TABLE player_stats(
		id UUID,
		player_id UUID,
		season INT,
		avg_min NUMERIC(4,2),
		fgm NUMERIC(5,2),
		fga NUMERIC(5,2),
		fg3m NUMERIC(5,2),
		fg3a NUMERIC(5,2),
		oreb NUMERIC(5,2),
		dreb NUMERIC(5,2),
		reb NUMERIC(5,2),
		ast NUMERIC(5,2),
		stl NUMERIC(5,2),
		blk NUMERIC(5,2),
		to NUMERIC(5,2),
		pf NUMERIC(4,2),
		pts NUMERIC(5,2),
		fg_pct NUMERIC(4,3),
		fg3_pct NUMERIC(4,3),
		ft_pct NUMERIC(4,3),
		CONSTRAINT fk_players
		FOREIGN KEY(player_id)
		REFERENCES players(uuid)
	);`

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

	db.MustExecContext(ctx, schema)

	return nil
}
