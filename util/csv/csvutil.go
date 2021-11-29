package csvHelpers

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	schema "github.com/tacohole/boardman/internal"

	"github.com/jszwec/csvutil"
)

func ReadCsv(filePath string) ([]schema.Source, error) {
	var sourceArray []schema.Source
	var source schema.Source

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
			log.Printf("error decoding %s into csv: %s", source.Name, err)
		}
		sourceArray = append(sourceArray, source)
	}

	return sourceArray, nil
}

func WriteCsv(fileName string, data []schema.DbSchema) error {
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
