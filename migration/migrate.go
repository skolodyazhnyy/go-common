package migration

import (
	"database/sql"
)

func Version(db *sql.DB, name string) (int, bool) {
	version := -1
	err := db.QueryRow(`SELECT version FROM migrations WHERE name = ?`, name).Scan(&version)

	return version, err == nil
}

func Apply(db *sql.DB, name string, migrations []string) error {
	version, ok := Version(db, name)
	if !ok {
		if err := initialize(db, name); err != nil {
			return err
		}
	}

	for v, m := range migrations {
		if v < version {
			continue
		}

		if err := execute(db, name, v+1, m); err != nil {
			return err
		}
	}

	return nil
}

func initialize(conn *sql.DB, name string) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		tx.Rollback()
	}()

	// as long as insert query below does not fail, we are fine
	tx.Exec(`CREATE TABLE migrations (name VARCHAR(25) PRIMARY KEY, version INT NOT NULL)`)

	if _, err := tx.Exec("INSERT INTO migrations VALUES(?, ?)", name, 0); err != nil {
		return err
	}

	return tx.Commit()
}

func execute(conn *sql.DB, name string, v int, m string) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		tx.Rollback()
	}()

	_, err = tx.Exec(m)
	if err != nil {
		return err
	}

	if _, err = tx.Exec("UPDATE migrations SET version = ? WHERE name = ?", v, name); err != nil {
		return err
	}

	return tx.Commit()
}
