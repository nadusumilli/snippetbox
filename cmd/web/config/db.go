package config

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	DSN *string
	*sql.DB
}

func NewDatabaseConnection(dsn *string, logger *slog.Logger) *Database {
	// Initialize the databse.
	db := &Database{
		DSN: dsn,
	}

	// Initializing the database connection and module.
	db.Init(logger)

	return db
}

// SetupDBConn initializes a new database connection using the provided DSN.
func (db *Database) Init(logger *slog.Logger) error {
	var err error
	db.DB, err = sql.Open("postgres", *db.DSN)
	if err != nil {
		logger.Error(err.Error())
		logger.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}

	err = db.DB.Ping()
	if err != nil {
		logger.Error(err.Error())
		db.DB.Close()
		logger.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}

	return nil
}
