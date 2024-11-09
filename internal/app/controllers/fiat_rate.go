package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/monero-atm/pricefetcher"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

func GetFiatRate(w http.ResponseWriter, r *http.Request) {
	// TODO: add more pairs by comparing to ECB.
	fiat := chi.URLParam(r, "fiat")
	fetcher := pricefetcher.New(nil)
	var (
		resp models.FiatRate
		err  error
	)
	resp.Price, resp.Source, err = fetcher.FetchXMRPrice(fiat)
	if err != nil {
		helpers.WriteError(w, apierror.ErrFetchRate, err)
	}
	json.NewEncoder(w).Encode(resp)
}
