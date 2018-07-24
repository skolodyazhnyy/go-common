package migration

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", "file::memory:")
	if err != nil {
		t.Fatalf("An error occurred while connection to database: %v", err)
	}

	migrations := []string{
		`CREATE TABLE test (value INT)`,
		`INSERT INTO test VALUES(1)`,
	}

	value := 0

	// apply few migrations
	if err = Apply(db, "foo", migrations); err != nil {
		t.Fatalf("An error occurred while applying migrations (1): %v", err)
	}

	// check current version
	if ver, _ := Version(db, "foo"); ver != 2 {
		t.Errorf("Expected version for \"foo\" is 2, but got %v instead", ver)
	}

	// check if migrations applied
	if err := db.QueryRow("SELECT MAX(value) FROM test").Scan(&value); err != nil {
		t.Fatalf("An error occurred while fetching value from test table: %v", err)
	}

	if value != 1 {
		t.Errorf("Expected value in test table is 1, but got %v instead", value)
	}

	// apply one more migrations
	migrations = append(migrations, `INSERT INTO test VALUES(2)`)
	if err = Apply(db, "foo", migrations); err != nil {
		t.Fatalf("An error occurred while applying migrations (2): %v", err)
	}

	// check current version
	if ver, _ := Version(db, "foo"); ver != 3 {
		t.Errorf("Expected version for \"foo\" is 3, but got %v instead", ver)
	}

	// check if migrations applied
	if err := db.QueryRow("SELECT MAX(value) FROM test").Scan(&value); err != nil {
		t.Fatalf("An error occurred while fetching value from test table: %v", err)
	}

	if value != 2 {
		t.Errorf("Expected value in test table is 2, but got %v instead", value)
	}
}
