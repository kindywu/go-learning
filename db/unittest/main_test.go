package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// postgres://kindy:kindy@postgresql-postgresql-master-1:5432/demo2

const (
	dsn          = "user=postgres password=dk82Fi2A_2d host=postgresql-postgresql-master-1 port=5432 sslmode=disable"
	test_db_name = "test_db"
	dsn_test_db  = "user=postgres password=dk82Fi2A_2d host=postgresql-postgresql-master-1 port=5432 dbname=test_db sslmode=disable"
)

func createTestDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", test_db_name))
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %v", err)
	}

	return db, nil
}

func dropTestDB(db *sql.DB) error {
	query := `
        SELECT pg_terminate_backend(pid)
        FROM pg_stat_activity
        WHERE pid <> pg_backend_pid() AND datname = $1;
    `
	_, err := db.Exec(query, test_db_name)
	if err != nil {
		return fmt.Errorf("failed to drop test database: %v", err)
	}

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", test_db_name))
	if err != nil {
		return fmt.Errorf("failed to drop test database: %v", err)
	}

	return db.Close()
}

func TestMain(m *testing.M) {
	fmt.Println("create test db")
	db, err := createTestDB()
	if err != nil {
		log.Fatalf("Could not create test DB: %v", err)
	}

	// Run tests
	code := m.Run()

	fmt.Println("drop test db")
	if err := dropTestDB(db); err != nil {
		log.Fatalf("Could not drop test DB: %v", err)
	}

	// Exit with the proper exit code
	os.Exit(code)
}

func TestExample(t *testing.T) {
	db, err := sql.Open("postgres", dsn_test_db)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Your test code here, e.g., creating tables, inserting test data, etc.
	_, err = db.Exec("CREATE TABLE example (id SERIAL PRIMARY KEY, name VARCHAR(255))")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Further testing logic...
}
