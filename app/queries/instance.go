package queries

import (
	"context"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
	db "gitlab.com/moneropay/metronero/metronero-backend/platform/database"
)

func InstanceInfo(ctx context.Context) (*models.InstanceInfo, error) {
	var (
		info models.InstanceInfo
		err error
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
	if err := row.Scan(&stats.WalletBalance, &stats.TotalProfits, &stats.TotalMerchants);
	    err != nil {
		return stats, err
	}
	return stats, nil
}
