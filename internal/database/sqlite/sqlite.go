package sqlite

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func ConnnectSQLite(schemaPath, databasePath string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping was unsuccessful: %w", err)
	}

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}

	db.MustExec(string(schemaBytes))

	return db, nil
}
