package queries

import (
	"context"

	db "gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/pkg/models"
)

func InstanceInfo(ctx context.Context) (*models.InstanceInfo, error) {
	var (
		info models.InstanceInfo
		err  error
	)
	info.Details, err = InstanceDetails(ctx)
	if err != nil {
		return nil, err
	}
	info.Stats, err = InstanceStats(ctx)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func InstanceDetails(ctx context.Context) (models.Instance, error) {
	var details models.Instance
	row, err := db.QueryRow(ctx,
		"SELECT version FROM instance")
	if err != nil {
		return details, err
	}
	if err := row.Scan(&details.Version); err != nil {
		return details, err
	}
	return details, nil
}

func InstanceStats(ctx context.Context) (models.InstanceStats, error) {
	var stats models.InstanceStats
	row, err := db.QueryRow(ctx,
		"SELECT wallet_balance,total_sales FROM instance_stats")
	if err != nil {
		return stats, err
	}
	if err := row.Scan(&stats.WalletBalance, &stats.TotalProfits); err != nil {
		return stats, err
	}
	return stats, nil
}

func UpdateInstance(ctx context.Context, conf *models.Instance) error {
	// TODO: deprecate
	return nil
}
