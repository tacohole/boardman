package csvutil

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/tacohole/boardman/internal/schema"

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
			log.Printf("error decoding %s into csv: %s", schema, err)
		}
		objArray = append(objArray, schema)
	}

	return objArray, nil
}

func WriteCsv(filePath string, data []struct{}) error {
	// check filename

	var buf bytes.Buffer

	w := csv.NewWriter(&buf)

	enc := csvutil.NewEncoder(w)
	defer w.Flush()

	for _, row := range data {
		if err := enc.Encode(row); err != nil {
			return fmt.Errorf("error: %s", err)
		}
	}

	// write to filename

	if err := w.Error(); err != nil {
		return fmt.Errorf("error: %s", err)
	}
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
