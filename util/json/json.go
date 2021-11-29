package jsonHelpers

import (
	"encoding/json"
	"log"
	"os"

	schema "github.com/tacohole/boardman/internal"
)

func WriteJson(fileName string, schema []schema.DbSchema) error {
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
