package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestPostgresDB(t *testing.T, user, password, host, port, db, ssl string) (*sql.DB, func(...string)) {
	t.Helper()
	databaseURL := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, db)
	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := database.Ping(); err != nil {
		t.Fatal(err)
	}

	return database, func(tables ...string) {
		defer database.Close()

		if len(tables) > 0 {
			if _, err := database.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}
	}
}
