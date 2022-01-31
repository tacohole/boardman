package get

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tacohole/boardman/internal"
	dbutil "github.com/tacohole/boardman/util/db"
)

var getPlayerSeasonCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Use:   "player-season",
	Run:   getPlayerSeasons,
}

func init() {
	GetCmd.AddCommand(getPlayerSeasonCmd)
}

func getPlayerSeasons(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	if err := dbutil.PrepareSchema(internal.PlayerSeasonSchema); err != nil {
		log.Fatalf("can't prepare player_season schema: %s", err)
	}

	if err := insertPlayerSeasonValues(); err != nil {
		log.Fatalf("can't insert player_season values: %s", err)
	}
}

func insertPlayerSeasonValues() error {
	return nil
}
