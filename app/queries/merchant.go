package queries

import (
	"time"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func GetMerchantInfo(id string) (models.MerchantDashboardInfo, error) {
	// As of now do nothing with the id as we will be returning dummy data
	_ = id

	t1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	t2, _ := time.Parse(time.RFC3339, "2023-03-11T15:04:05Z")
	t3, _ := time.Parse(time.RFC3339, "2023-01-28T15:04:05Z")

	return models.MerchantDashboardInfo{
		Stats: models.MerchantStats{
			AccountId:  "dummy",
			Balance:    13123862123,
			TotalSales: 9387397238432,
		},
		Recent: []models.Payment{
			models.Payment{
				MerchantName: "Siren",
				Amount:       12834232321,
				Fee:          2728323,
				LastUpdate:   t1,
				Status:       "Finished",
			},
			models.Payment{
				MerchantName: "Siren",
				Amount:       8263721312312,
				Fee:          3432642,
				LastUpdate:   t2,
				Status:       "Finished",
			},
			models.Payment{
				MerchantName: "Siren",
				Amount:       38894233298788,
				Fee:          2337432423,
				LastUpdate:   t3,
				Status:       "Finished",
			},
		},
	}, nil
}
