package csvHelpers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jszwec/csvutil"
	schema "github.com/tacohole/boardman/internal"
)

func ReadCsv(filePath string) ([]struct{}, error) {
	var sourceArray []struct{}
	var source struct{}

	headerSlice, err := csvutil.Header(source, "csv")
	if err != nil {
		return nil, err
	}

	headers := strings.Join(headerSlice, ",")

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("file at %s cannot be opened: %s", filePath, err)
	}

	csvReader := csv.NewReader(file)

	dec, err := csvutil.NewDecoder(csvReader, headers)
	if err != nil {
		return nil, err
	}

	for {
		if err := dec.Decode(&source); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("error decoding into csv: %s", err)
		}
		sourceArray = append(sourceArray, source)
	}

	return sourceArray, nil
}

func WriteTeamCsv(fileName string, data []schema.Data) error {
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

func WriteTeamJson(fileName string, schema []schema.Data) error {
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

func WritePlayerCsv(fileName string, data []schema.Player) error {
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

func WritePlayerJson(fileName string, schema []schema.Player) error {
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
