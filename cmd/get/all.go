package get

import (
	"encoding/json"
	"log"
	"time"

	"github.com/tacohole/boardman/internal/schema"
	"github.com/tacohole/boardman/util/csvutil"
	"github.com/tacohole/boardman/util/httputil"

	"github.com/spf13/cobra"
)

var getAllCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getAll,
}

func init() {
	getAllCmd.Flags().StringVar(&writeTo, "output", "", "output type, options are JSON or csv")

}

func getAll(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	sourceCache, err := getSourceCache()
	if err != nil {
		log.Fatalf("sources file unavailable: %s", err)
	}

	for _, source := range sourceCache {
		// get some data
		resp, err := httputil.MakeHttpRequest("GET", source.Url, nil, "")
		if err != nil {
			log.Printf("error:", err)
		}
		defer resp.Body.Close()

		var schema source.DbSchema

		err = json.Unmarshal(resp.Body, &schema)
		if err != nil {
			log.Printf("error", err)
		}

		if writeTo == "csv" {
			// write a csv
			fileName := (source.Name + time.Now())

			csvutil.WriteCsv(fileName, &schema)
		}

	}

	if writeTo == "JSON" {
		// write a JSON file
	}

}

func getSourceCache() ([]schema.Source, error) {
	var sourceCache []schema.Source

	sourceCache, err := csvutil.ReadCsv("~/sources.csv", schema.Source)
	if err != nil {
		return nil, err
	}

	return sourceCache, nil
}
