package get

import "github.com/spf13/cobra"

var getCoachesCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "coaches",
	Run:   getCoaches,
}

func getCoaches(cmd *cobra.Command, args []string) {
	// get endpoint
	// make structs in internal
	//
}
