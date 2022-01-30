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

			// gameStats values here

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
