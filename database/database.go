package database

import (
	"database/sql"
	"log"
	"net"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// If connection string is a postgres URL, resolve hostname to an IPv4 address
	if strings.HasPrefix(connectionString, "postgres://") || strings.HasPrefix(connectionString, "postgresql://") {
		if u, err := url.Parse(connectionString); err == nil {
			host := u.Hostname()
			port := u.Port()
			if port == "" {
				port = "5432"
			}
			if ips, err := net.LookupIP(host); err == nil {
				for _, ip := range ips {
					if ip4 := ip.To4(); ip4 != nil {
						u.Host = net.JoinHostPort(ip4.String(), port)
						log.Printf("Resolved DB host %s -> %s", host, ip4.String())
						break
					}
				}
			} else {
				log.Printf("Could not resolve DB host %s: %v", host, err)
			}
			
			// Add SSL mode if not present (Supabase requires SSL)
			query := u.Query()
			if !query.Has("sslmode") {
				query.Set("sslmode", "require")
				u.RawQuery = query.Encode()
				log.Println("Added sslmode=require to connection string")
			}
			
			connectionString = u.String()
		}
	}

	// Open database
	db, err := sql.Open("postgres", connectionString)
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
