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

	playerCache, err := getPlayerIdCache()
	if err != nil {
		log.Fatalf("can't get player cache: %s", err)
	}

	for i := 1979; i < 2021; i++ {
		for _, player := range *playerCache {
			stats, err := getPlayerSeasonAverages(i, player)
			if err != nil {
				log.Printf("can't get stats for %d", player.BDL_ID)
			}
			err = insertPlayerSeasonAverages(stats)
			if err != nil {
				log.Printf("could not insert stats for %d", player.BDL_ID)
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
	) VALUES (
		:uuid,
		:balldontlie_id,
		:date,
		:home_id, 
		:home_score, 
		:visitor_id, 
		:visitor_score, 
		:season, 
		:is_postseason, 
		:winner_id, 
		:margin )`,
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
	getUrl := httpHelpers.BaseUrl + httpHelpers.Stats + "?seasons[]=" + fmt.Sprint(season) + "&player_ids[]=" + fmt.Sprint(player.BDL_ID)
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
		playerYear.PlayerID = player.ID
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

func prepareStatsSchema() error {
	schema := `CREATE TABLE player_season_avgs(
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
