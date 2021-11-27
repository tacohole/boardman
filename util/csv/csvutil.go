package csvutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/jszwec/csvutil"
)

func ReadCsv(filePath string, schema struct{}) ([]struct{}, error) {
	var objArray []struct{}

	schemaKeys, err := SchemaHeadersToString(schema)
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
		var o struct{}
		if err := dec.Decode(&o); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		objArray = append(objArray, o)
	}

	return objArray, nil
}

func WriteCsv(filePath string, schema string) error {

	return nil
}

func SchemaHeadersToString(schema struct{}) (string, error) {

	headers := string.(reflect.VisibleFields())



	return headers, nil
}
