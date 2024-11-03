package daemon

import "time"

func Start() {
	go expireInvoices(5 * time.Minute)
}
