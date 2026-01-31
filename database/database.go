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
		separator := "?"
		if strings.Contains(connectionString, "?") {
			separator = "&"
		}

		if !strings.Contains(connectionString, "sslmode=") {
			connectionString = connectionString + separator + "sslmode=require"
			separator = "&"
			log.Println("Added sslmode=require to connection string")
		}

		// Disable prepared statements for Supabase pooler (port 6543)
		if !strings.Contains(connectionString, "default_query_exec_mode=") {
			connectionString = connectionString + separator + "default_query_exec_mode=simple_protocol"
			log.Println("Added default_query_exec_mode=simple_protocol for pooler compatibility")
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
