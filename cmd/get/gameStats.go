package get

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

// no queries just paginate
var getGameStatsCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "games-stats",
	Run:   getGameStats,
}

func init() {
	GetCmd.AddCommand(getGamesCmd)
}

func getGameStats(cmd *cobra.Command, args []string) {
	// loadDefaultVariables()

	s := schema.SingleGame{}

	err := schema.PrepareSeasonSchema()
	if err != nil {
		log.Fatalf("could not create games schema: %s", err)
	}

	for i := 1979; i < 2021; i++ {
		games, err := s.GetAllGameStats(i)
		if err != nil {
			log.Fatalf("can't get games: %s", err)
		}

		err = insertGameStats(games)
		if err != nil {
			log.Printf("can't insert games for season %d: %s", i, err)
		}
	}

}

func insertGameStats(stats []schema.SingleGame) error {
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

	result, err := db.NamedExecContext(ctx, `INSERT INTO player_game_stats (
		player_uuid,
		game_uuid,
		team_id,
		min,
		fgm,
		fga,
		fg3m,
		fg3a,
		oreb,
		dreb,
		reb,
		ast,
		stl,
		blk,
		to,
		pf,
		pts,
		fg_pct,
		fg3_pct,
		ft_pct,
	) VALUES (
		:player_uuid,
		:game_uuid,
		:team_id,
		:min,
		:fgm,
		:fga,
		:fg3m,
		:fg3a,
		:oreb,
		:dreb,
		:reb,
		:ast,
		:stl,
		:blk,
		:to,
		:pf,
		:pts,
		:fg_pct,
		:fg3_pct,
		:ft_pct, )`,
		stats)
	if err != nil {
		log.Printf("Insert failed, %s", result)
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
