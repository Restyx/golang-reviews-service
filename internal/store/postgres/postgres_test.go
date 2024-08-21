package postgres_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var (
	pgUser string
	pgPass string
	pgHost string
	pgPort string
	pgDB   string
	pgSSL  string
)

func TestMain(m *testing.M) {
	godotenv.Load()

	if pgUser = os.Getenv("POSTGRES_USER"); pgUser == "" {
		pgUser = "postgres"
	}
	if pgPass = os.Getenv("POSTGRES_PASS"); pgPass == "" {
		pgPass = "postgres"
	}
	if pgHost = os.Getenv("POSTGRES_HOST"); pgHost == "" {
		pgHost = "localhost"
	}
	if pgPort = os.Getenv("POSTGRES_PORT"); pgPort == "" {
		pgPort = "5432"
	}
	if pgDB = os.Getenv("POSTGRES_DB"); pgDB == "" {
		pgDB = "reviews_testing"
	}
	if pgSSL = os.Getenv("POSTGRES_SSL"); pgDB == "" {
		pgSSL = "disable"
	}

	os.Exit(m.Run())
}
