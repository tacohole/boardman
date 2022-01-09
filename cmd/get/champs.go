package get

import "github.com/spf13/cobra"

var getAwardWinnersCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "champs",
	Run:   getAwardWinners,
}

func init() {
	GetCmd.AddCommand(getAwardWinnersCmd)
}

func getAwardWinners(cmd *cobra.Command, args []string) {

}

// init schema?

// call to champs/playoff results http://data.nba.net/prod/v1/2016/playoffsBracket.json
// call to MVP
// insert MVP
// call to ROY
// insert ROY
// call to COY
// insert COY
