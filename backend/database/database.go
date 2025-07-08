// database/database.go
// Handles the database connection and initial table creation.
package database

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
)

func ConnectDB(connStr string) *sql.DB {
	if connStr == "" {
		log.Fatal().Msg("Database URL is not configured")
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to open database connection")
	}

	// --- FIX: Add retry logic ---
	// This loop will attempt to connect to the database up to 5 times
	// over 10 seconds before failing. This gives the database container
	// time to start up properly.
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			log.Info().Msg("Successfully connected to the PostgreSQL database!")
			return db
		}
		log.Warn().Err(err).Int("attempt", i+1).Msg("Failed to connect to database. Retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal().Err(err).Msg("Unable to connect to the database after several retries")
	return nil // This line will not be reached due to log.Fatal
}

func CreateTables(db *sql.DB) {
	// User table
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY, email TEXT NOT NULL UNIQUE, password_hash TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL);`); err != nil {
		log.Fatal().Err(err).Msg("Unable to create users table")
	}
	// Categories table
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS categories (id UUID PRIMARY KEY, user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, name TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL, UNIQUE(user_id, name));`); err != nil {
		log.Fatal().Err(err).Msg("Unable to create categories table")
	}
	// Transactions table
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS transactions (id UUID PRIMARY KEY, user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT, description TEXT, amount NUMERIC(10, 2) NOT NULL, date TIMESTAMPTZ NOT NULL, type TEXT NOT NULL);`); err != nil {
		log.Fatal().Err(err).Msg("Unable to create transactions table")
	}
	// Budgets table
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS budgets (id UUID PRIMARY KEY, user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE, amount NUMERIC(10, 2) NOT NULL, month INT NOT NULL, year INT NOT NULL, created_at TIMESTAMPTZ NOT NULL, UNIQUE(user_id, category_id, month, year));`); err != nil {
		log.Fatal().Err(err).Msg("Unable to create budgets table")
	}
	// Refresh tokens table
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS refresh_tokens (id UUID PRIMARY KEY, user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, token_hash TEXT NOT NULL UNIQUE, expires_at TIMESTAMPTZ NOT NULL);`); err != nil {
		log.Fatal().Err(err).Msg("Unable to create refresh_tokens table")
	}

	log.Info().Msg("All tables are ready.")
}
