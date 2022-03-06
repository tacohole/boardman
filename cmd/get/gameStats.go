package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	internal "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

// no queries just paginate
var getGameStatsCmd = &cobra.Command{
	Short: "gets detailed box score stats on game stats by player for all games since 1979",
	Long:  "gets detailed box score stats on game stats by player for all games since 1979",
	Use:   "games-stats",
	Run:   getGameStats,
}

func init() {
	GetCmd.AddCommand(getGameStatsCmd)
}

func getGameStats(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	if err := dbutil.PrepareSchema(internal.GameStatsSchema); err != nil {
		log.Fatalf("could not create games schema: %s", err)
	}

	for i := 1979; i <= 2021; i++ {

		if err := getGameStatsPage(i); err != nil {
			log.Fatalf("can't get game stats page for %d: %s", i, err)
		}

	}

	// add our UUIDs to new table
	playerResult, err := updateGamesWithPlayerIds()
	if err != nil || playerResult < 1 {
		log.Fatalf("can't add player UUIDs to player_game_stats: %s", err)
	}

	teamResult, err := updateGamesWithTeamIds()
	if err != nil || teamResult < 1 {
		log.Fatalf("can't add team UUIDs to player_game_stats: %s", err)
	}

	gameResult, err := updateGamesWithGameIds()
	if err != nil || gameResult < 1 {
		log.Fatalf("can't add game UUIDs to player_game_stats: %s", err)
	}

	if err := cleanupDuplicateStats(); err != nil {
		log.Printf("can't remove duplicate values from games: %s", err)
	}

}

func cleanupDuplicateStats() error {
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

	q := `DELETE FROM player_game_stats a 
		USING player_game_stats b 
		WHERE a.uuid < b.uuid 
		AND a.balldontlie_id = b.balldontlie_id;`

	_, err = tx.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

func insertGameStatsPage(stats []internal.SingleGame) error {
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

	_, err = tx.NamedExecContext(ctx, `INSERT INTO player_game_stats (
		uuid,
		balldontlie_id,
		player_bdl_id,
		team_bdl_id,
		game_bdl_id,
		season,
		min,
		fgm,
		fga,
		fg3m,
		fg3a,
		ftm,
		fta,
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
		:game_bdl_id,
		:season,
		:min,
		:fgm,
		:fga,
		:ftm,
		:fta,
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

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func gameStatsPagetoStruct(page internal.Page) ([]internal.SingleGame, error) {
	var s internal.SingleGame
	var games []internal.SingleGame

	for _, d := range page.Data {
		s.UUID = uuid.New()
		s.BDL_ID = d.ID
		s.GameBDL_ID = d.Game.BDL_ID
		s.PlayerBDL_ID = d.Player.BDL_ID
		s.TeamBDL_ID = d.Team.BDL_ID
		s.Season = d.Game.Season
		s.AST = d.AST
		s.BLK = d.BLK
		s.DREB = d.DREB
		s.FG3A = d.FG3A
		s.FG3M = d.FG3M
		s.FG3_PCT = d.FG3_PCT
		s.FGA = d.FGA
		s.FGM = d.FGM
		s.FG_PCT = d.FG_PCT
		s.FTA = d.FTA
		s.FTM = d.FTM
		s.FT_PCT = d.FT_PCT
		s.Minutes = d.Minutes
		s.OREB = d.OREB
		s.PF = d.OREB
		s.PF = d.PF
		s.PTS = d.PTS
		s.REB = d.REB
		s.STL = d.STL
		s.TO = d.TO
		games = append(games, s)
	}

	return games, nil
}

func getGameStatsPage(season int) error {
	var page internal.Page
	var errorCount int

	for pageIndex := 0; pageIndex <= page.PageData.TotalPages; pageIndex++ {
		getUrl := internal.BDLUrl + internal.BDLStats + "?seasons[]=" + fmt.Sprint(season) + "&page=" + fmt.Sprint(pageIndex) + "&per_page=100"

		resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
		if err != nil {
			errorCount++
			if errorCount > 2 {
				return err
			} else {
				pageIndex--
				continue
			}
		}
		defer resp.Body.Close()

		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(r, &page); err != nil {
			return err
		}

		gameStats, err := gameStatsPagetoStruct(page)
		if err != nil {
			return fmt.Errorf("can't convert game stats page to struct: %s", err)
		}

		if err := insertGameStatsPage(gameStats); err != nil {
			return fmt.Errorf("can't insert stats: %s", err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func updateGamesWithPlayerIds() (int64, error) {
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

	stmt := `UPDATE player_game_stats 
			SET player_uuid = players.uuid
			FROM players
			WHERE player_game_stats.player_bdl_id = players.balldontlie_id;`

	result := tx.MustExecContext(ctx, stmt)
	tx.Commit()

	return result.RowsAffected()
}

func updateGamesWithGameIds() (int64, error) {
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

	stmt := `UPDATE player_game_stats
			SET game_uuid = games.uuid
			FROM games
			WHERE player_game_stats.game_bdl_id = games.balldontlie_id;`

	result := tx.MustExecContext(ctx, stmt)
	tx.Commit()

	return result.RowsAffected()
}

func updateGamesWithTeamIds() (int64, error) {
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

	stmt := `UPDATE player_game_stats 
			SET team_uuid = teams.uuid
			FROM teams 
			WHERE player_game_stats.team_bdl_id = teams.balldontlie_id;`

	result := tx.MustExecContext(ctx, stmt)
	tx.Commit()

	return result.RowsAffected()
}
