package get

import (
	"log"

	schema "github.com/tacohole/boardman/internal"

	"github.com/spf13/cobra"
)

var getPlayersCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getPlayers,
}

func init() {
	getPlayersCmd.Flags().StringVar(&writeTo, "output", "", "output type, options are JSON or csv")

	getCmd.AddCommand(getPlayersCmd)

}

func getPlayers(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	p := schema.Player{}

	_, err := p.GetAllPlayers()
	if err != nil {
		log.Fatalf("can't get players: %s", err)
	}

	// insert into database

}
