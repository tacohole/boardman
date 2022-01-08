package get

import (
	"log"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
)

// no queries just paginate
var getGameStatsCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "games-stats",
	Run:   getGameStats,
}

func init() {
	GetCmd.AddCommand(getGamesCmd)
}

func getGameStats(cmd *cobra.Command, args []string) {
	// loadDefaultVariables()

	s := schema.SingleGame{}

	err := schema.PrepareSeasonSchema()
	if err != nil {
		log.Fatalf("could not create games schema: %s", err)
	}

	for i := 1979; i < 2021; i++ {
		games, err := s.GetAllGameStats()
		if err != nil {
			log.Fatalf("can't get games: %s", err)
		}

		_, err = insertGameStats(games)
		if err != nil {
			log.Printf("can't insert games for season %d: %s", i, err)
		}
	}

}
