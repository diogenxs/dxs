package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	// DefaultDBPath is the default path to the sqlite3 database
	defaultDBPath = "~/.dxs.sqlite"

	// Create SQLite table for storing alerts
	createTableSQL = `CREATE TABLE IF NOT EXISTS alerts (
		"fingerprint" TEXT PRIMARY KEY,
		"labels" TEXT,
		"annotations" TEXT,
		"state" TEXT,
		"generatorURL" TEXT,
		"acknowledged" BOOLEAN DEFAULT FALSE
);`
)

// GetDB returns a pointer to a sql.DB object
func GetDB() (*sql.DB, error) {
	dbPath := viper.GetString("dbPath")
	if dbPath == "" {
		home, _ := homedir.Dir()
		dbPath = home + "/.dxs.sqlite"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	return db, nil
}

func MigrateDB() error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	// Create the alerts table if it doesn't exist
	if _, err := db.Exec(createTableSQL); err != nil {
		return err
	}

	return nil
}
