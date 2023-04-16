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
		"SELECT version,default_commission,custodial_mode,registrations_allowed,withdrawal_times FROM instance")
	if err != nil {
		return details, err
	}
	if err := row.Scan(&details.Version, &details.DefaultCommission, &details.CustodialMode,
		&details.RegistrationsAllowed, &details.WithdrawalTimes); err != nil {
		return details, err
	}
	return details, nil
}

func InstanceStats(ctx context.Context) (models.InstanceStats, error) {
	var stats models.InstanceStats
	row, err := db.QueryRow(ctx,
		"SELECT wallet_balance,total_profits,total_merchants FROM instance_stats")
	if err != nil {
		return stats, err
	}
	if err := row.Scan(&stats.WalletBalance, &stats.TotalProfits, &stats.TotalMerchants); err != nil {
		return stats, err
	}
	return stats, nil
}

func UpdateInstance(ctx context.Context, conf *models.Instance) error {
	return db.Exec(ctx,
		"UPDATE instance SET custodial_mode=$1,default_commission=$2,registrations_allowed=$3,withdrawal_times=$4",
		conf.CustodialMode, conf.DefaultCommission, conf.RegistrationsAllowed, conf.WithdrawalTimes)
}
