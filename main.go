package main

import (
	"gitlab.com/metronero/backend/platform/database"
	"gitlab.com/metronero/backend/utils/config"
	"gitlab.com/metronero/backend/utils/server"
)

func main() {
	config.Load()
	database.Migrate()
	database.Init()
	server.StartWithGracefulShutdown()
}
