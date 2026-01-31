package database

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Add sslmode if not present (Supabase requires SSL)
	if strings.HasPrefix(connectionString, "postgres://") || strings.HasPrefix(connectionString, "postgresql://") {
		if !strings.Contains(connectionString, "sslmode=") {
			separator := "?"
			if strings.Contains(connectionString, "?") {
				separator = "&"
			}
			connectionString = connectionString + separator + "sslmode=require"
			log.Println("Added sslmode=require to connection string")
		}
	}

	// Open database with pgx driver
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (optional tapi recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
