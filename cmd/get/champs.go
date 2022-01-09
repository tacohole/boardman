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
