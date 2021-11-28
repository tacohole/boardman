package csvutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"boardman/internal/schema"

	"github.com/jszwec/csvutil"
)

func ReadCsv(filePath string, schema struct{}) ([]struct{}, error) {
	var objArray []struct{}

	schemaKeys, err := StructKeysToString(schema)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("file at %s cannot be opened: %s", filePath, err)
	}

	csvReader := csv.NewReader(file)

	dec, err := csvutil.NewDecoder(csvReader, schemaKeys)
	if err != nil {
		return nil, err
	}

	for {
		if err := dec.Decode(&schema); err == io.EOF {
			break
		} else if err != nil {
			log.Printf("error decoding into csv: %s", err)
		}
		objArray = append(objArray, schema)
	}

	return objArray, nil
}

func WriteCsv(filePath string, schema string) error {

	return nil
}

func StructKeysToString(data schema.Data) (string, error) {
	var headerSlice []string
	structFields := reflect.VisibleFields(data)

	for _, field := range structFields {
		header := field.Name

		headerSlice = append(headerSlice, header)
	}

	headers := strings.Join(headerSlice, ",")

	return headers, nil
}
