package main

import (
	"gitlab.com/moneropay/metronero/metronero-backend/utils/config"
	"gitlab.com/moneropay/metronero/metronero-backend/platform/database"
	"gitlab.com/moneropay/metronero/metronero-backend/utils/server"
)

func main() {
	config.Load()
	database.Migrate()
	database.Init()
	server.StartWithGracefulShutdown()
}
