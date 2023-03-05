package daemon

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var pdb *pgxpool.Pool

func pdbConnect() {
	var err error
	if pdb, err = pgxpool.Connect(context.Background(), Config.postgresCS); err != nil {
		log.Fatal().Err(err).Msg("Startup failure")
	}
}

func pdbMigrate() {
	m, err := migrate.New("file://db", Config.postgresCS)
	if err != nil {
		log.Fatal().Err(err).Msg("Startup failure")
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
		log.Fatal().Err(err).Msg("Startup failure")
	}
}

func pdbQueryRow(ctx context.Context, query string, args ...interface{}) (pgx.Row, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	c := make(chan pgx.Row, 1)
	go func() { c <- pdb.QueryRow(ctx, query, args...) }()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case row := <-c:
		return row, nil
	}
}

func pdbQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	type queryRet struct { rows pgx.Rows; err error }
	c := make(chan queryRet, 1)
	go func() {
		rows, err := pdb.Query(ctx, query, args...)
		c <- queryRet{rows, err}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ret := <-c:
		return ret.rows, ret.err
	}
}

func pdbExec(ctx context.Context, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	c := make(chan error, 1)
	go func() {
		_, err := pdb.Exec(ctx, query, args...)
		c <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}
