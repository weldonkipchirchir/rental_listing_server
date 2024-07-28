package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
)

var testQueries *db.Queries
var testDB *sql.DB

// TestMain sets up the test database and initializes the test environment.
func TestMain(m *testing.M) {
	var err error

	// Initialize a test database connection
	testDB, err = sql.Open("postgres", "user=root dbname=rental_listing_test password=secret sslmode=disable")
	if err != nil {
		log.Fatalf("cannot connect to test database: %v", err)
	}

	// Initialize the SQLC queries with the test database connection
	testQueries = db.New(testDB)

	// Run tests
	code := m.Run()

	// Clean up resources (close database connection, etc.)
	err = testDB.Close()
	if err != nil {
		log.Fatalf("error closing test database connection: %v", err)
	}

	// Exit with the status code from tests
	os.Exit(code)
}
