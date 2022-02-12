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
	Short: "calculates team stats for all seasons in database",
	Long:  "calculates regular season stats for all teams/seasons, along with rosters and coaches",
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

			// build rosters
			roster, err := buildRoster(ts)
			if err != nil {
				log.Fatalf("%s", err)
			}
			ts.Roster = roster

			// build coaches only recent seasons
			if ts.Season >= 2015 {
				coaches, err := buildCoaches(ts)
				if err != nil {
					log.Fatalf("%s", err)
				}
				ts.Coaches = coaches
			}

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

	// return sum fgm where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(fgm)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var fgm float32 // init return value
	if err = rows.Scan(fgm); err != nil {
		return nil, err
	}

	return &fgm, nil
}

// for player in team in season sum fga on all games
func sumFga(ts internal.TeamSeason) (*float32, error) {
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

	// return sum fga where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(fga)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var fga float32 // init return value
	if err = rows.Scan(fga); err != nil {
		return nil, err
	}

	return &fga, nil
}

// for player in team in season sum ftm on all games
func sumFtm(ts internal.TeamSeason) (*float32, error) {
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

	// return sum ftm where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(ftm)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var ftm float32 // init return value
	if err = rows.Scan(ftm); err != nil {
		return nil, err
	}

	return &ftm, nil
}

// for player in team in season sum fta on all games
func sumFta(ts internal.TeamSeason) (*float32, error) {
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

	// return sum fta where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(fta)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var fta float32 // init return value
	if err = rows.Scan(fta); err != nil {
		return nil, err
	}

	return &fta, nil
}

// for player in team in season sum fg3m on all games
func sumFg3m(ts internal.TeamSeason) (*float32, error) {
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

	// return sum fg3m where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(fg3m)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var fg3m float32 // init return value
	if err = rows.Scan(fg3m); err != nil {
		return nil, err
	}

	return &fg3m, nil
}

// for player in team in season sum fg3a on all games
func sumFg3a(ts internal.TeamSeason) (*float32, error) {
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

	// return sum fgm where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(fg3a)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var fg3a float32 // init return value
	if err = rows.Scan(fg3a); err != nil {
		return nil, err
	}

	return &fg3a, nil
}

// for player in team in season sum Oreb on all games
func sumOreb(ts internal.TeamSeason) (*float32, error) {
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

	// return sum oreb where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(oreb)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var oreb float32 // init return value
	if err = rows.Scan(oreb); err != nil {
		return nil, err
	}

	return &oreb, nil
}

// for player in team in season sum dreb on all games
func sumDreb(ts internal.TeamSeason) (*float32, error) {
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

	// return sum fgm where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(dreb)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var dreb float32 // init return value
	if err = rows.Scan(dreb); err != nil {
		return nil, err
	}

	return &dreb, nil
}

// for player in team in season sum ast on all games
func sumAst(ts internal.TeamSeason) (*float32, error) {
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

	// return sum ast where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(ast)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var ast float32 // init return value
	if err = rows.Scan(ast); err != nil {
		return nil, err
	}

	return &ast, nil
}

// for player in team in season sum blk on all games
func sumBlk(ts internal.TeamSeason) (*float32, error) {
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

	// return sum blk where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(blk)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var blk float32 // init return value
	if err = rows.Scan(blk); err != nil {
		return nil, err
	}

	return &blk, nil
}

// for player in team in season sum stl on all games
func sumStl(ts internal.TeamSeason) (*float32, error) {
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

	// return sum stl where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(stl)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var stl float32 // init return value
	if err = rows.Scan(stl); err != nil {
		return nil, err
	}

	return &stl, nil
}

// for player in team in season sum pf on all games
func sumPf(ts internal.TeamSeason) (*float32, error) {
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

	// return sum pf where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(pf)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var pf float32 // init return value
	if err = rows.Scan(pf); err != nil {
		return nil, err
	}

	return &pf, nil
}

// for player in team in season sum turnovers on all games
func sumTurnovers(ts internal.TeamSeason) (*float32, error) {
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

	// return sum turnovers where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(turnovers)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var to float32 // init return value
	if err = rows.Scan(to); err != nil {
		return nil, err
	}

	return &to, nil
}

// for player in team in season sum pts on all games
func sumPts(ts internal.TeamSeason) (*float32, error) {
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

	// return sum pts where team_uuid = ts.team_uuid and game_uuid in :season
	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT SUM(pts)
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid 
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var pts float32 // init return value
	if err = rows.Scan(pts); err != nil {
		return nil, err
	}

	return &pts, nil
}

// select player_uuid from gamestats where season = season and teamuuid = team.uuid
func buildRoster(ts internal.TeamSeason) ([]uuid.UUID, error) {
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

	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT DISTINCT player_uuid
		FROM player_game_stats
		INNER JOIN games ON player_game_stats.game_uuid = games.uuid
		WHERE team_uuid=:team_uuid
		AND  games.season = :season
		AND is_postseason = 'f'`,
		ts)
	if err != nil {
		return nil, err
	}

	var roster []uuid.UUID
	if err = rows.Scan(roster); err != nil {
		return nil, err
	}

	return roster, nil
}

// select coach.uuid from coaches where season = season and team_uuid = team.uuid
func buildCoaches(ts internal.TeamSeason) ([]uuid.UUID, error) {
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

	// a048985d-0531-46ae-b8d3-121595957f9c - Hawks
	rows, err := db.NamedQueryContext(ctx,
		`SELECT uuid
		FROM coaches
		WHERE team_uuid=:team_uuid
		AND  season = :season`,
		ts)
	if err != nil {
		return nil, err
	}

	var coaches []uuid.UUID
	if err = rows.Scan(coaches); err != nil {
		return nil, err
	}

	return coaches, nil
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

// get conf rank - sort by wpct within conference
// get overall rank - sort by wpct
