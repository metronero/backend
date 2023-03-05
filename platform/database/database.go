package database

import (
	"context"
	"time"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gitlab.com/moneropay/metronero/metronero-backend/utils/config"
)

var Db *pgxpool.Pool

func Init() {
	var err error
	if Db, err = pgxpool.Connect(context.Background(), config.PostgresUri); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}

func Migrate() {
	m, err := migrate.New("file://platform/migrations", config.PostgresUri)
	if err != nil {
		log.Fatal("Failed to begin migration: ", err)
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
		log.Fatal("Failed to migrate: ", err)
	}
}

func QueryRow(ctx context.Context, query string, args ...interface{}) (pgx.Row, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	c := make(chan pgx.Row, 1)
	go func() { c <- Db.QueryRow(ctx, query, args...) }()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case row := <-c:
		return row, nil
	}
}

func Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	type queryRet struct { rows pgx.Rows; err error }
	c := make(chan queryRet, 1)
	go func() {
		rows, err := Db.Query(ctx, query, args...)
		c <- queryRet{rows, err}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ret := <-c:
		return ret.rows, ret.err
	}
}

func Exec(ctx context.Context, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	c := make(chan error, 1)
	go func() {
		_, err := Db.Exec(ctx, query, args...)
		c <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}
