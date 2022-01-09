package get

import (
	"github.com/spf13/cobra"
)

// no queries just paginate
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
	// loadDefaultVariables()
	getTeamData(getTeamsCmd, []string{})
	getPlayers(getPlayersCmd, []string{})
	getGames(getGamesCmd, []string{})
	getGameStats(getGameStatsCmd, []string{})
	getPlayerStats(getPlayerStatsCmd, []string{})
	getCoaches(getCoachesCmd, []string{})
	getAwardWinners(getAwardWinnersCmd, []string{})

}