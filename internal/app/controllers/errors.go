package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"gitlab.com/metronero/metronero-go/models"
)

func writeError(w http.ResponseWriter, apiErr *models.ApiError, err error) {
	w.WriteHeader(apiErr.Status)
	json.NewEncoder(w).Encode(apiErr)
	log.Error().Err(err).Int("code", apiErr.Code).Msg(apiErr.Msg)
}
