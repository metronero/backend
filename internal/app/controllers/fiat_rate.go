package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/monero-atm/pricefetcher"
	"gitlab.com/metronero/backend/internal/app/queries"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetFiatRate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := ctx.Value("account_id").(string)
	settings, err := queries.GetMerchantSettings(r.Context(), id)
	if err != nil {
		helpers.WriteError(w, apierror.ErrDatabase, err)
	}
	fetcher := pricefetcher.New(nil)
	var resp models.FiatRate
	resp.Price, resp.Source, err = fetcher.FetchXMRPrice(*settings.FiatCurrency)
	if err != nil {
		helpers.WriteError(w, apierror.ErrFetchRate, err)
	}
	json.NewEncoder(w).Encode(resp)
}
