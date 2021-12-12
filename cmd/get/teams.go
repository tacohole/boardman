package get

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/spf13/cobra"
	schema "github.com/tacohole/boardman/internal"
	httpHelpers "github.com/tacohole/boardman/util/http"
)

var getTeamDataCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getTeamData,
}

func init() {
	getTeamDataCmd.Flags().StringVar(&writeTo, "output", "", "output type")

	getTeamDataCmd.MarkFlagRequired("source")
	getTeamDataCmd.MarkFlagRequired("writeTo")

}

func getTeamData(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	t := time.Now().Format(time.UnixDate)

	var response struct {
		Data []schema.Team   `json:"data"`
		Meta schema.PageData `json:"meta"`
	}

	pageIndex := 0

	getUrl := httpHelpers.BaseUrl + "teams" + fmt.Sprint(pageIndex)

	for pageIndex := 0; pageIndex < response.Meta.TotalPages; pageIndex++ {
		// get some data
		resp, err := httpHelpers.MakeHttpRequest("GET", getUrl, nil, "")
		if err != nil {
			log.Printf("error: %s", err)
		}
		defer resp.Body.Close()

		r, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error: %s", err)
		}

		err = json.Unmarshal(r, &response)
		if err != nil {
			log.Printf("error: %s", err)
		}

		if writeTo == "csv" {
			// write a csv
			fileName := ("player_data" + fmt.Sprint(pageIndex) + t + ".csv")
			err = writeTeamCsv(fileName, response.Data)
			if err != nil {
				log.Printf("Error writing csv %s: %s", fileName, err)
			}
		}

		if writeTo == "JSON" {
			// write a JSON file
			fileName := ("player_data" + fmt.Sprint(pageIndex) + t + ".json")
			err = writeTeamJson(fileName, response.Data)
			if err != nil {
				log.Printf("Error writing JSON file %s: %s", fileName, err)
			}
		}

	}
}

func writeTeamCsv(fileName string, data []schema.Team) error {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to create file %s: %s", fileName, err)
	}

	w := csv.NewWriter(file)

	enc := csvutil.NewEncoder(w)
	defer w.Flush()

	for _, row := range data {
		if err := enc.Encode(row); err != nil {
			return fmt.Errorf("error: %s", err)
		}
	}

	if err := w.Error(); err != nil {
		return fmt.Errorf("error: %s", err)
	}
	return nil
}

func writeTeamJson(fileName string, schema []schema.Team) error {
	file, err := os.OpenFile(fileName, os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create file %s: %s", fileName, err)
	}

	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	if err := enc.Encode(schema); err != nil {
		log.Printf("error: %s", err)
	}

	return nil
}
