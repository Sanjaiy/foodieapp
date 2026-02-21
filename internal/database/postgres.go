package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Connect(ctx context.Context, databaseURL string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	maxRetries := 5
	for i := range maxRetries {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			log.Printf("Failed to open database: %v (attempt %d/%d)", err, i+1, maxRetries)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.PingContext(ctx)
		if err == nil {
			break
		}

		log.Printf("Failed to ping database: %v (attempt %d/%d)", err, i+1, maxRetries)
		db.Close()

		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to database after %d attempts: %w", maxRetries, err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	log.Println("DB connection established!!!")
	return db, nil
}
