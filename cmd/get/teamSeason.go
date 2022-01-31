package get

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getTeamSeasonCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "team-season",
	Run:   getTeamSeasons,
}

func init() {
	GetCmd.AddCommand(getTeamSeasonCmd)
}

func getTeamSeasons(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	teams, err := internal.GetTeamCache()
	if err != nil {
		log.Fatalf("can't get team cache: %s", err)
	}

	ts := internal.TeamSeason{}

	for i := 1979; i <= 2021; i++ {
		for _, team := range teams {

			ts.UUID = uuid.New()
			ts.TeamUUID = team.UUID
			ts.Season = i

			// game values first
			w, err := sumTeamWins(ts)
			if err != nil {
				log.Fatalf("can't sum team wins: %s", err)
			}
			ts.Wins = *w

			l, err := sumTeamLosses(ts)
			if err != nil {
				log.Fatalf("can't sum team wins: %s", err)
			}
			ts.Losses = *l

			// calc wpct for team
			ts.WinPct = float32(ts.Wins) / (float32(ts.Wins) + float32(ts.Losses))

			playoffs, err := setMadePlayoffs(ts)
			if err != nil {
				log.Fatalf("can't set made playoffs: %s", err)
			}
			ts.MadePlayoffs = *playoffs

			if *playoffs {
				psWins, err := sumPsWins(ts)
				if err != nil {
					log.Fatalf("can't sum playoff wins: %s", err)
				}
				ts.PSWins = *psWins

				psLosses, err := sumPsLosses(ts)
				if err != nil {
					log.Fatalf("can't sum playoff losses: %s", err)
				}
				ts.PSLosses = *psLosses

			}

			plus, err := calculatePlus(ts)
			if err != nil {
				log.Fatalf("can't calculate plus/minus: %s", err)
			}
			minus, err := calculateMinus(ts)
			if err != nil {
				log.Fatalf("can't calculate plus/minus: %s", err)
			}
			ts.PlusMinus = (*plus - *minus)

			// gameStats values here

			fgm, err := sumFgm(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.FGM = *fgm

			fga, err := sumFga(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.FGA = *fga

			ftm, err := sumFtm(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.FTM = *ftm

			fta, err := sumFta(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.FTA = *fta

			fg3m, err := sumFg3m(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.FG3M = *fg3m

			fg3a, err := sumFg3a(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.FG3A = *fg3a

			oreb, err := sumOreb(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.OREB = *oreb

			dreb, err := sumDreb(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.DREB = *dreb
			ts.REB = (*dreb + *oreb)

			ast, err := sumAst(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.AST = *ast

			blk, err := sumBlk(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.BLK = *blk

			stl, err := sumStl(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.STL = *stl

			to, err := sumTurnovers(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.TO = *to

			pf, err := sumPf(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.PF = *pf

			pts, err := sumPts(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.PTS = *pts

			ts.FG_PCT = (ts.FGM/ts.FGA + ts.FGM)
			ts.FT_PCT = (ts.FTM/ts.FTA + ts.FTM)
			ts.FG3_PCT = (ts.FG3M/ts.FG3A + ts.FG3M)

			if err = insertTeamSeasonRecord(ts); err != nil {
				log.Printf("can't insert team record: %s", err)
			}

		}
	}

}

func sumTeamWins(ts internal.TeamSeason) (*int, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// return count of games where winner = team && is_postseason = false
	rows, err := db.NamedQueryContext(ctx,
		`SELECT COUNT(*) 
		FROM games
		WHERE season = :season
		AND winner_id = :team_uuid 
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var wins int // init return value
	if err = rows.Scan(wins); err != nil {
		return nil, err
	}

	return &wins, nil
}

// return count of games where winner != team && is_postseason = false
func sumTeamLosses(ts internal.TeamSeason) (*int, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// return count of games where winner = team && is_postseason = false
	rows, err := db.NamedQueryContext(ctx,
		`SELECT COUNT(*) 
		FROM games
		WHERE season = :season
		AND (visitor_id = :team_uuid OR home_id = :team_uuid)
		AND winner_id != :team_uuid 
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var losses int // init return value
	if err = rows.Scan(losses); err != nil {
		return nil, err
	}

	return &losses, nil
}

// for team in season where is_postseason = true
func setMadePlayoffs(ts internal.TeamSeason) (*bool, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	var playoffGames int
	var madePlayoffs bool

	rows, err := db.NamedQueryContext(ctx,
		`SELECT COUNT(*)
		FROM games
		WHERE season = :season
		AND (visitor_id = :team_uuid OR home_id = :team_uuid)
		AND is_postseason = 't'`,
		ts)
	if err != nil {
		return nil, err
	}
	if err = rows.Scan(playoffGames); err != nil {
		return nil, err
	}

	if playoffGames > 0 {
		madePlayoffs = true
	} else {
		madePlayoffs = false
	}

	return &madePlayoffs, nil
}

func sumPsWins(ts internal.TeamSeason) (*int, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// return count of games where winner = team && is_postseason = false
	rows, err := db.NamedQueryContext(ctx,
		`SELECT COUNT(*)
		FROM games
		WHERE season = :season
		AND winner_id = :team_uuid 
		AND is_postseason = 't'`,
		ts)
	if err != nil {
		return nil, err
	}

	var wins int // init return value
	if err = rows.Scan(wins); err != nil {
		return nil, err
	}

	return &wins, nil
}

func sumPsLosses(ts internal.TeamSeason) (*int, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// return count of games where winner = team && is_postseason = false
	rows, err := db.NamedQueryContext(ctx,
		`SELECT COUNT(*)
		FROM games
		WHERE season = :season
		AND (visitor_id = :team_uuid OR home_id = :team_uuid)
		AND winner_id != :team_uuid 
		AND is_postseason = 't'`,
		ts)
	if err != nil {
		return nil, err
	}

	var losses int // init return value
	if err = rows.Scan(losses); err != nil {
		return nil, err
	}

	return &losses, nil
}

func calculatePlus(ts internal.TeamSeason) (*int, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// return margin where winner = team && is_postseason = false
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(margin)
		FROM games
		WHERE season = :season
		AND winner_id = :team_uuid 
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var plus int // init return value
	if err = rows.Scan(plus); err != nil {
		return nil, err
	}

	return &plus, nil
}

func calculateMinus(ts internal.TeamSeason) (*int, error) {
	db, err := dbutil.DbConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// return margin where winner = team && is_postseason = false
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(margin)
		FROM games
		WHERE season = :season
		AND (home_id=:team_uuid OR visitor_id=team_uuid)
		AND winner_id != :team_uuid 
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var minus int // init return value
	if err = rows.Scan(minus); err != nil {
		return nil, err
	}

	return &minus, nil
}

// for player in team in season sum fgm on all games
func sumFgm(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum fga on all games
func sumFga(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum ftm on all games
func sumFtm(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum fta on all games
func sumFta(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum fg3m on all games
func sumFg3m(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum fg3a on all games
func sumFg3a(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum Oreb on all games
func sumOreb(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum dreb on all games
func sumDreb(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum ast on all games
func sumAst(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum blk on all games
func sumBlk(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum stl on all games
func sumStl(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum pf on all games
func sumPf(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum turnovers on all games
func sumTurnovers(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

// for player in team in season sum pts on all games
func sumPts(ts internal.TeamSeason) (*float32, error) {
	return nil, nil
}

func insertTeamSeasonRecord(ts internal.TeamSeason) error {
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

	query := `INSERT INTO team_season(
		uuid
		team_uuid
		season
		wins
		losses
		postseason_wins
		postseason_losses
		wpct
		plus_minus
		conf_rank
		ovr_rank
		made_playoffs
		fgm
		fga
		ftm
		fta
		fg3m
		fg3a
		oreb
		dreb
		reb
		ast
		stl
		blk
		turnovers
		pf
		pts
		fg_pct
		fg3_pct
		ft_pct
		roster
		coaches,
		VALUES
		:uuid
		:team_uuid
		:season
		:wins
		:losses
		:postseason_wins
		:postseason_losses
		:wpct
		:plus_minus
		:conf_rank
		:ovr_rank
		:made_playoffs
		:fgm
		:fga
		:ftm
		:fta
		:fg3m
		:fg3a
		:oreb
		:dreb
		:reb
		:ast
		:stl
		:blk
		:turnovers
		:pf
		:pts
		:fg_pct
		:fg3_pct
		:ft_pct
		:roster
		:coaches);`

	_, err = tx.NamedExecContext(ctx, query, ts)
	if err != nil {
		return err
	}

	return nil
}

// calculate +/-: sum win margin, sum loss margin, subtract win from loss
// get conf rank - sort by wpct within conference
// get overall rank - sort by wpct
