package get

import (
	"encoding/json"
	"io"
	"log"
	"time"

	schema "github.com/tacohole/boardman/internal"
	csvHelpers "github.com/tacohole/boardman/util/csv"
	httpHelpers "github.com/tacohole/boardman/util/http"
	jsonHelpers "github.com/tacohole/boardman/util/json"

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

	var data []schema.DbSchema

	t := time.Now().Format(time.UnixDate)

	sourceCache, err := getSourceCache()
	if err != nil {
		log.Fatalf("sources file unavailable: %s", err)
	}

	for _, source := range sourceCache {
		// get some data
		resp, err := httpHelpers.MakeHttpRequest("GET", source.Url, nil, "")
		if err != nil {
			log.Printf("error:", err)
		}
		defer resp.Body.Close()

		r, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error: %s", err)
		}

		err = json.Unmarshal(r, &data)
		if err != nil {
			log.Printf("error: %s", err)
		}

		if writeTo == "csv" {
			// write a csv
			fileName := (source.Name + t + ".csv")
			err = csvHelpers.WriteCsv(fileName, data)
			if err != nil {
				log.Printf("Error writing csv %s: %s", fileName, err)
			}
		}

		if writeTo == "JSON" {
			// write a JSON file
			fileName := (source.Name + t + ".json")
			err = jsonHelpers.WriteJson(fileName, data)
			if err != nil {
				log.Printf("Error writing JSON file %s: %s", fileName, err)
			}
		}

	}

}

func getSourceCache() ([]schema.Source, error) {
	var sourceCache []schema.Source

	sourceCache, err := csvHelpers.ReadCsv("~/sources.csv")
	if err != nil {
		return nil, err
	}

	return sourceCache, nil
}
