package dbutil

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/tacohole/boardman/util/config"
)

func dbConn() {
	psqlconn := fmt.Sprintf(config.DbUrlEnvironmentName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatalf("%s", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf("Connected to PostgreSQL")
}
