package csvHelpers

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jszwec/csvutil"
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
