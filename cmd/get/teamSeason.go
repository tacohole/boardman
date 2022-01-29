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

			w, err := sumTeamWins(team, i)
			if err != nil {
				log.Fatalf("can't sum team wins: %s", err)
			}
			ts.Wins = *w

			l, err := sumTeamLosses(team, i)
			if err != nil {
				log.Fatalf("can't sum team wins: %s", err)
			}
			ts.Losses = *l

			// calc wpct for team
			ts.WinPct = float32(ts.Wins) / (float32(ts.Wins) + float32(ts.Losses))

			playoffs, err := setMadePlayoffs(team, i)
			if err != nil {
				log.Fatalf("can't set made playoffs: %s", err)
			}
			ts.MadePlayoffs = *playoffs

			if *playoffs {
				psWins, err := sumPsWins(team, i)
				if err != nil {
					log.Fatalf("can't sum playoff wins: %s", err)
				}
				ts.PSWins = *psWins

				psLosses, err := sumPsLosses(team, i)
				if err != nil {
					log.Fatalf("can't sum playoff losses: %s", err)
				}
				ts.PSLosses = *psLosses

			}

			if err = insertTeamSeasonRecord(ts); err != nil {
				log.Printf("can't insert team record: %s", err)
			}

		}
	}

}

func sumTeamWins(t internal.Team, season int) (*int, error) {

	// lookup all games for team in season
	// return count of games where winner = team && is_postseason = false

	return nil, nil
}

func sumTeamLosses(t internal.Team, season int) (*int, error) {
	// lookup all games for team in season
	// return count of games where winner != team && is_postseason = false

	return nil, nil
}

func setMadePlayoffs(t internal.Team, season int) (*bool, error) {
	// for team in season where is_postseason = true

	return nil, nil
}

func sumPsWins(t internal.Team, season int) (*int, error) {
	return nil, nil
}

func sumPsLosses(t internal.Team, season int) (*int, error) {
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
