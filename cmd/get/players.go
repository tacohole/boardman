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
	schema "github.com/tacohole/boardman/internal"
	httpHelpers "github.com/tacohole/boardman/util/http"

	"github.com/spf13/cobra"
)

var getPlayersCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getPlayerData,
}

func init() {
	getPlayersCmd.Flags().StringVar(&writeTo, "output", "", "output type, options are JSON or csv")

	getCmd.AddCommand(getPlayersCmd)

}

func getPlayerData(cmd *cobra.Command, args []string) {
	loadDefaultVariables()

	var data []schema.Player

	t := time.Now().Format(time.UnixDate)

	var response schema.Page

	pageIndex := 0

	getUrl := httpHelpers.BaseUrl + httpHelpers.Players + fmt.Sprint(pageIndex)

	for pageIndex := 0; pageIndex < response.PageData.TotalPages; pageIndex++ {
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

		err = json.Unmarshal(r, &data)
		if err != nil {
			log.Printf("error: %s", err)
		}

		if writeTo == "csv" {
			// write a csv
			fileName := ("player_data" + fmt.Sprint(pageIndex) + t + ".csv")
			err = writePlayerCsv(fileName, data)
			if err != nil {
				log.Printf("Error writing csv %s: %s", fileName, err)
			}
		}

		if writeTo == "JSON" {
			// write a JSON file
			fileName := ("player_data" + fmt.Sprint(pageIndex) + t + ".json")
			err = writePlayerJson(fileName, data)
			if err != nil {
				log.Printf("Error writing JSON file %s: %s", fileName, err)
			}
		}

	}
}

func writePlayerCsv(fileName string, data []schema.Player) error {
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

func writePlayerJson(fileName string, schema []schema.Player) error {
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
