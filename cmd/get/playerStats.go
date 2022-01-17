package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

var getPlayerStatsCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "player-stats",
	Run:   getPlayerStats,
}

func init() {
	GetCmd.AddCommand(getPlayerStatsCmd)
}

func getPlayerStats(cmd *cobra.Command, args []string) {

	err := schema.PreparePlayerStatsSchema()
	if err != nil {
		log.Fatalf("can't create schema for stats: %s", err)
	}

	playerCache, err := getPlayerIdCache()
	if err != nil {
		log.Fatalf("can't get player cache: %s", err)
	}

	for i := 1979; i < 2021; i++ {
		for _, player := range *playerCache {
			stats, err := getPlayerSeasonAverages(i, player)
			if err != nil {
				log.Printf("can't get stats for %d", player.BDL_ID)
				continue // don't insert
			}
			err = insertPlayerSeasonAverages(stats)
			if err != nil {
				log.Fatalf("could not insert stats for %d", player.BDL_ID)
			}
		}
	}

}

func insertPlayerSeasonAverages(stats *schema.PlayerYear) error {
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

	result, err := db.NamedExecContext(ctx, `INSERT INTO player_season_avgs (
		player_id,
		league_year,
		avg_min,
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
		ft_pct
	) VALUES (
		:player_id,
		:league_year,
		:avg_min,
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
		:ft_pct )`,
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

func getPlayerSeasonAverages(season int, player schema.Player) (*schema.PlayerYear, error) {
	getUrl := schema.BDLUrl + schema.BDLStats + "?seasons[]=" + fmt.Sprint(season) + "&player_ids[]=" + fmt.Sprint(player.BDL_ID)
	resp, err := httpHelpers.MakeHttpRequest("GET", getUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var page schema.Page
	var playerYear schema.PlayerYear

	err = json.Unmarshal(r, &page)
	if err != nil {
		return nil, err
	}
	for _, d := range page.Data {
		playerYear.PlayerID = player.UUID
		playerYear.LeagueYear = d.LeagueYear
		playerYear.GamesPlayed = d.GamesPlayed
		playerYear.Minutes = d.Minutes
		playerYear.FGM = d.FGM
		playerYear.FGA = d.FGA
		playerYear.FG3M = d.FG3M
		playerYear.FG3A = d.FG3A
		playerYear.OREB = d.OREB
		playerYear.DREB = d.DREB
		playerYear.REB = d.REB
		playerYear.AST = d.AST
		playerYear.STL = d.STL
		playerYear.BLK = d.BLK
		playerYear.TO = d.TO
		playerYear.PF = d.PF
		playerYear.PTS = d.PTS
		playerYear.FG_PCT = d.FG_PCT
		playerYear.FG3_PCT = d.FG3_PCT
		playerYear.FT_PCT = d.FT_PCT
	}
	time.Sleep(1 * time.Second) // avoiding 429 from BDL
	return &playerYear, nil
}

func getPlayerIdCache() (*[]schema.Player, error) {
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

	p := []schema.Player{}

	result, err := tx.NamedExecContext(ctx, `SELECT uuid,balldontlie_id FROM players`, p)
	if err != nil {
		log.Printf("Select failed, %s", result)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &p, nil
}
