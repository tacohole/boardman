package get

import (
	"log"

	"github.com/spf13/cobra"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getPaidCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "paid",
	Run:   getPaid,
}

func init() {
	GetCmd.AddCommand(getPaidCmd)
}

func getPaid(cmd *cobra.Command, args []string) {
	loadDefaultVariables()
	const InitDb = `CREATE DATABASE nba_data;`

	if err := dbutil.PrepareSchema(InitDb, ""); err != nil {
		log.Fatalf("can't create database, %s", err)
	}

	getTeamData(getTeamsCmd, []string{})
	getPlayers(getPlayersCmd, []string{})
	getCoaches(getCoachesCmd, []string{})
	getGames(getGamesCmd, []string{})
	getGameStats(getGameStatsCmd, []string{})
}
