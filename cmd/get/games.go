package get

import (
	"log"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
)

var getGamesCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getGames,
}

func init() {
	getGamesCmd.Flags().StringVar(&writeTo, "output", "", "output type, options are JSON or csv")

	getGamesCmd.AddCommand(getPlayersCmd)

}

func getGames(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	g := schema.Game{}
	seasons := []int{1981: 2021}

	for _, season := range seasons {
		_, err := g.GetSeasonGames(season)
		if err != nil {
			log.Fatalf("can't get games: %s", err)
		}

		// insert into database
	}

}
