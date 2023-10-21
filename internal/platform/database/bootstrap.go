package database

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/google/uuid"

	"gitlab.com/metronero/backend/internal/utils/config"
	"gitlab.com/metronero/backend/internal/utils/auth"
)

func applyMigrations() {
	m, err := migrate.New("file://internal/platform/migrations", config.PostgresUri)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
		log.Fatal().Err(err).Msg("Failed to migrate database schema")
	}
}

func bootstrap() {
	defer log.Info().Msg("Instance bootstrapped, starting...")
	ctx := context.Background()
	row, err := QueryRow(ctx, "SELECT version FROM instance")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to bootstrap")
	}
	var version string
	err = row.Scan(&version)
	if err == nil {
		return
	}
	if err != pgx.ErrNoRows {
		log.Fatal().Err(err).Msg("Failed to bootstrap")
	}
	log.Info().Msg("Bootstrapping instance for the first time")
	tx, err := Db.Begin(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start database transaction")
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, "INSERT INTO instance (version) VALUES ($1)",
		config.Version); err != nil {
		log.Fatal().Err(err).Msg("Failed to bootstrap instance table")
	}

	if _, err := tx.Exec(ctx, "INSERT INTO instance_stats DEFAULT VALUES"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bootstrap instance_stats table")
	}
	// Create admin user and generate a password.
	password, err := auth.GeneratePassword()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate password for the admin user.")
	}
	hashBytes, err := auth.HashPassword(password)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to hash password for the admin user.")
	}
	if _, err := tx.Exec(ctx,
		"INSERT INTO accounts (account_id, username, password, role) VALUES ($1, 'admin', $2, 'admin')",
		uuid.New().String(), hashBytes);
		err != nil {
		log.Fatal().Err(err).Msg("Failed to bootstrap")
	}
	if err := tx.Commit(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to bootstrap")
	}
	log.Info().Str("username", "admin").Str("password", password).Msg("Created admin user")
}
