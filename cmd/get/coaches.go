package get

import "github.com/spf13/cobra"

var getCoachesCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "coaches",
	Run:   getCoaches,
}

func getCoaches(cmd *cobra.Command, args []string) {
	// get endpoint http://data.nba.net/prod/v1/{year}/coaches.json
	// make structs in internal
	//
}
