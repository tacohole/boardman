package get

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getPlayerSeasonCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "player-season",
	Run:   getPlayerSeasons,
}

func init() {
	GetCmd.AddCommand(getPlayerSeasonCmd)
}

func getPlayerSeasons(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	if err := dbutil.PrepareSchema(internal.PlayerSeasonSchema); err != nil {
		log.Fatalf("can't prepare player_season schema: %s", err)
	}

	playerCache, err := internal.GetPlayerIdCache()
	if err != nil {
		log.Fatalf("can't get player cache: %s", err)
	}

	ps := internal.PlayerSeason{}

	for _, player := range playerCache {
		for i := 1979; i <= 2021; i++ {
			ps.PlayerUUID = player.UUID
			ps.Season = i
			ps.BDL_ID = player.BDL_ID

			// find values
			gp, err := getGamesPlayed(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.GamesPlayed = *gp

			mins, err := getAvgMinutes(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.Minutes = *mins

			fgm, err := getFgm(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.FGM = *fgm

			fga, err := getFga(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.FGA = *fga
			// fg_pct
			ps.FG_PCT = ps.FGM / (ps.FGM + ps.FGA)

			ftm, err := getFtm(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.FTM = *ftm

			fta, err := getFta(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.FTA = *fta

			ps.FT_PCT = ps.FTM / (ps.FTM + ps.FTA)

			fg3m, err := getFg3m(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.FG3M = *fg3m

			fg3a, err := getFg3a(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.FG3A = *fg3a

			ps.FG3_PCT = ps.FG3M / (ps.FG3M + ps.FG3A)

			oreb, err := getOreb(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.OREB = *oreb

			dreb, err := getDreb(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.DREB = *dreb

			ps.REB = (ps.OREB + ps.DREB)

			ast, err := getAst(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.AST = *ast

			stl, err := getStl(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.STL = *stl

			blk, err := getBlk(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.BLK = *blk

			tos, err := getTos(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.TO = *tos

			pf, err := getPf(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.PF = *pf

			pts, err := getPts(ps)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ps.PTS = *pts

			// fg3_pct
			// ft_pct

			if err := insertPlayerSeasonValues(ps); err != nil {
				log.Fatalf("can't insert player_season values: %s", err)
			}

		}
	}

}

func getGamesPlayed(ps internal.PlayerSeason) (*int, error) {
	return nil, nil
}

func getAvgMinutes(ps internal.PlayerSeason) (*string, error) {
	return nil, nil
}

func getFgm(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getFga(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getFg3m(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getFtm(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getFta(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getFg3a(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getOreb(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getDreb(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

// sum rebs

func getAst(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getStl(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getBlk(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getTos(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getPf(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func getPts(ps internal.PlayerSeason) (*float32, error) {
	return nil, nil
}

func insertPlayerSeasonValues(ps internal.PlayerSeason) error {
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

	q := `INSERT INTO player_season(
		uuid,
		balldontlie_id,
		season,
		avg_min,
		fgm,
		fga,
		ftm,
		fta,
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
		VALUES (
		:uuid,
		:balldontlie_id,
		:season,
		:avg_min,
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
		:ft_pct)`

	_, err = tx.NamedExecContext(ctx, q, ps)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
