package dbutil

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jmoiron/sqlx"
)

func DbConn(dbName string) (*sqlx.DB, error) {
	connString := os.Getenv("DATABASE_URL") + dbName
	if connString == "" {
		log.Fatalf("No database connection string provided")
	}

	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func GenerateTimeout() (*time.Duration, error) {
	timeoutStr := os.Getenv("DB_TIMEOUT")

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, err
	}

	return &timeout, nil
}

func PrepareSchema(schema string, dbName string) error {
	db, err := DbConn(dbName)
	if err != nil {
		return err
	}

	timeout, err := GenerateTimeout()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	db.MustExecContext(ctx, schema)

	return nil
}

// github.com/jackc/pgx/v4/pgxpool - look into connection pooling later

// exec the schema or fail; multi-statement Exec behavior varies between
// database drivers;  pq will exec them all, sqlite3 won't, ymmv

// tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
// tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
// tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)", "United States", "New York", "1")
// tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Hong Kong", "852")
// tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
// // Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person

// // Query the database, storing results in a []Person (wrapped in []interface{})
// people := []Person{}
// db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
// jason, john := people[0], people[1]

// fmt.Printf("%#v\n%#v", jason, john)
// // Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
// // Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}

// // You can also get a single result, a la QueryRow
// jason = Person{}
// err = db.Get(&jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
// fmt.Printf("%#v\n", jason)
// // Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}

// // if you have null fields and use SELECT *, you must use sql.Null* in your struct
// places := []Place{}
// err = db.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
// if err != nil {
//     fmt.Println(err)
//     return
// }
// usa, singsing, honkers := places[0], places[1], places[2]

// fmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
// // Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
// // Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
// // Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}

// // Loop through rows using only one struct
// place := Place{}
// rows, err := db.Queryx("SELECT * FROM place")
// for rows.Next() {
//     err := rows.StructScan(&place)
//     if err != nil {
//         log.Fatalln(err)
//     }
//     fmt.Printf("%#v\n", place)
// }
// // Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
// // Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
// // Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}

// // Named queries, using `:name` as the bindvar.  Automatic bindvar support
// // which takes into account the dbtype based on the driverName on sqlx.Open/Connect
// _, err = db.NamedExec(`INSERT INTO person (first_name,last_name,email) VALUES (:first,:last,:email)`,
//     map[string]interface{}{
//         "first": "Bin",
//         "last": "Smuth",
//         "email": "bensmith@allblacks.nz",
// })

// // Selects Mr. Smith from the database
// rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`, map[string]interface{}{"fn": "Bin"})

// // Named queries can also use structs.  Their bind names follow the same rules
// // as the name -> db mapping, so struct fields are lowercased and the `db` tag
// // is taken into consideration.
// rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:first_name`, jason)

// // batch insert

// // batch insert with structs
// personStructs := []Person{
//     {FirstName: "Ardie", LastName: "Savea", Email: "asavea@ab.co.nz"},
//     {FirstName: "Sonny Bill", LastName: "Williams", Email: "sbw@ab.co.nz"},
//     {FirstName: "Ngani", LastName: "Laumape", Email: "nlaumape@ab.co.nz"},
// }

// _, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
//     VALUES (:first_name, :last_name, :email)`, personStructs)

// // batch insert with maps
// personMaps := []map[string]interface{}{
//     {"first_name": "Ardie", "last_name": "Savea", "email": "asavea@ab.co.nz"},
//     {"first_name": "Sonny Bill", "last_name": "Williams", "email": "sbw@ab.co.nz"},
//     {"first_name": "Ngani", "last_name": "Laumape", "email": "nlaumape@ab.co.nz"},
// }

// _, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
//     VALUES (:first_name, :last_name, :email)`, personMaps)
