package get

import (
	"boardman/internal/schema"
	"boardman/util/csvutil"

	"github.com/spf13/cobra"
)

var getAllCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getAll,
}

func init() {
	getAllCmd.Flags().StringVar(&writeTo, "output", "", "output type")

}

func getAll(cmd *cobra.Command, args []string) {
	loadDefaultVariables()
}

func getSourceCache() ([]schema.Source, error) {
	var sourceCache []schema.Source
	var headerString string

	sourceCache, err := csvutil.ReadCsv("~/sources.csv", headerString)
	if err != nil {
		return nil, err
	}

	return sourceCache, nil
}
