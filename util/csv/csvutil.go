package csvutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jszwec/csvutil"
)

type Schema []string

func readCsv(filePath string, s Schema) ([]struct{}, error) {
	var objArray []struct{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("file at %s cannot be opened: %s", filePath, err)
	}

	r, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(r)

	dec, err := csvutil.NewDecoder(csvReader, s)

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
