package daemon

import (
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/metronero/backend/internal/app/queries"
)

func expireInvoices(d time.Duration) {
	for {
		if err := queries.ExpireIncompleteInvoices(); err != nil {
			log.Err(err).Msg("Failed to mark invoice(s) as expired")
		}
		time.Sleep(d)
	}
}
