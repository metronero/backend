package queries

import (
	"time"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func GetAdminInfo() (models.AdminDashboardInfo, error) {
	t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	return models.AdminDashboardInfo{
		Stats: models.InstanceStats{
			WalletBalance: 23250000000000,
			TotalProfits: 67887000000000,
			TotalMerchants: 6,
		},
		Recent: []models.Withdrawal{
			models.Withdrawal{
				MerchantName: "Siren",
        			Amount: 12832321,
        			Date: t,
			},
			models.Withdrawal{
				MerchantName: "Koutsie",
        			Amount: 8263721312312,
        			Date: t,
			},
			models.Withdrawal{
				MerchantName: "Moin",
        			Amount: 38894233298788,
        			Date: t,
			},
		},
	}, nil
}
