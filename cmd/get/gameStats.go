package get

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	internal "github.com/tacohole/boardman/internal"
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
	GetCmd.AddCommand(getGameStatsCmd)
}

func getGameStats(cmd *cobra.Command, args []string) {
	// loadDefaultVariables()

	err := internal.PrepareGameStatsSchema()
	if err != nil {
		log.Fatalf("could not create games schema: %s", err)
	}

	for i := 2020; i <= 2021; i++ {
		var page internal.Page

		for pageIndex := 0; pageIndex <= page.PageData.TotalPages; pageIndex++ {
			gameSlice, err := internal.GetGameStatsPage(i, pageIndex)
			if err != nil {
				log.Fatalf("%s", err)
			}

			if err = insertGameStats(gameSlice); err != nil {
				log.Fatalf("can't insert games for season %d: %s", i, err)
			}

			time.Sleep(1000 * time.Millisecond) // more 429 dodging
		}

	}

	// add our UUIDs to new table
	if err = updateGamesWithPlayerIds(); err != nil {
		log.Fatalf("can't add player UUIDs to player_game_stats: %s", err)
	}

	if err = updateGamesWithTeamIds(); err != nil {
		log.Fatalf("can't add team UUIDs to player_game_stats: %s", err)
	}

	if err = updateGamesWithGameIds(); err != nil {
		log.Fatalf("can't add game UUIDs to player_game_stats: %s", err)
	}

}

func insertGameStats(stats []internal.SingleGame) error {
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

	_, err = tx.NamedExecContext(ctx, `INSERT INTO player_game_stats (
		uuid,
		balldontlie_id,
		player_bdl_id,
		team_bdl_id,
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
		turnovers,
		pf,
		pts,
		fg_pct,
		fg3_pct,
		ft_pct) 
		VALUES(
		:uuid,
		:balldontlie_id,
		:player_bdl_id,
		:team_bdl_id,
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
		:turnovers,
		:pf,
		:pts,
		:fg_pct,
		:fg3_pct,
		:ft_pct)`,
		stats)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func updateGamesWithPlayerIds() error {
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

	stmt := `UPDATE player_game_stats 
			SET player_uuid = players.uuid
			FROM players
			WHERE player_game_stats.player_bdl_id = players.balldontlie_id;`

	tx.MustExecContext(ctx, stmt)

	return nil
}

func updateGamesWithGameIds() error {
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

	stmt := `UPDATE player_game_stats
			SET game_uuid = games.uuid
			FROM games
			WHERE player_game_stats.balldontlie_id = games.balldontlie_id;`

	tx.MustExecContext(ctx, stmt)

	return nil
}

func updateGamesWithTeamIds() error {
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

	stmt := `UPDATE player_game_stats 
			SET team_uuid = teams.uuid
			FROM teams 
			WHERE player_game_stats.team_bdl_id = teams.balldontlie_id;`

	tx.MustExecContext(ctx, stmt)

	return nil
}
