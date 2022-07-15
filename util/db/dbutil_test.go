package dbutil_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	dbutil "github.com/tacohole/boardman/util/db"
)

const (
	TEST_DB_URL     = "postgresql://127.0.0.1:5432/"
	TEST_DB_NAME    = "nba_data_test"
	TEST_DB_TIMEOUT = "10s"
)

// test DbConn
func TestDbConn(t *testing.T) {
	before()

	testDb, err := dbutil.DbConn("")
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	proxy := &sqlx.DB{}

	if reflect.TypeOf(testDb) != reflect.TypeOf(proxy) {
		after(false)
		t.Fail()
	}
	after(false)
}

func TestDbConnWithName(t *testing.T) {
	before()

	testDb, err := dbutil.DbConn(TEST_DB_NAME)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}
	proxy := &sqlx.DB{}

	if reflect.TypeOf(testDb) != reflect.TypeOf(proxy) {
		log.Printf("%v", reflect.TypeOf(testDb))
		t.Fail()
		after(true)
	}
	after(true)
}

func TestGenerateTimeout(t *testing.T) {
	before()

	timeout, err := dbutil.GenerateTimeout()
	if err != nil {
		log.Printf("%v", err)
	}

	proxy := time.Duration(10 * time.Second)

	if reflect.TypeOf(timeout) != reflect.TypeOf(&proxy) {
		log.Printf("%v", reflect.TypeOf(timeout))
		after(false)
		t.Fail()
	}

	after(false)
}

// test db create
func TestPrepareValidSchema(t *testing.T) {
	before()

	schema := fmt.Sprintf("CREATE DATABASE %s", TEST_DB_NAME)

	err := dbutil.PrepareSchema(schema, "")
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		after(true)
	}
	after(true)
}

func before() {
	os.Setenv("DATABASE_URL", TEST_DB_URL)
	os.Setenv("DB_TIMEOUT", TEST_DB_TIMEOUT)
}

func after(clearDb bool) {
	if clearDb {
		stmt := fmt.Sprintf("DROP DATABASE %s;", TEST_DB_NAME)
		err := dbutil.PrepareSchema(stmt, "")
		if err != nil {
			log.Printf("error on test teardown: %v", err)
		}
	}

	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("DB_TIMEOUT")
}
