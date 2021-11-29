package jsonHelpers

import (
	"log"
	"os"

	schema "github.com/tacohole/boardman/internal"
)

func WriteJson(fileName string, schema []schema.DbSchema) error {
	_, err := os.Create(fileName)
	if err != nil {
		log.Printf("Failed to create file %s: %s", fileName, err)
	}
	return nil
}
