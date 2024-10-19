package queries

import (
	"context"

	db "gitlab.com/metronero/backend/internal/platform/database"
)

func InstanceVersion(ctx context.Context) (string, error) {
	var version string
	row, err := db.QueryRow(ctx,
		"SELECT version FROM instance")
	if err != nil {
		return version, err
	}
	if err := row.Scan(&version); err != nil {
		return version, err
	}
	return version, nil
}
