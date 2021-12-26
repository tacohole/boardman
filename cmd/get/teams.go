package get

import (
	"log"

	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
)

var getTeamsCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getTeamData,
}

func init() {
	getTeamsCmd.Flags().StringVar(&writeTo, "output", "", "output type")

	getTeamsCmd.MarkFlagRequired("writeTo")

}

func getTeamData(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	team := schema.Team{}

	_, err := team.GetAllTeams()
	if err != nil {
		log.Fatalf("can't get teams: %s", err)
	}

	// insert into database

}
