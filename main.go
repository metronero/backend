package main

import (
	"gitlab.com/metronero/backend/internal/platform/database"
	"gitlab.com/metronero/backend/internal/utils/config"
	"gitlab.com/metronero/backend/internal/utils/daemon"
	"gitlab.com/metronero/backend/internal/utils/server"
)

func main() {
	config.Load()
	database.Init()
	daemon.Start()
	server.StartWithGracefulShutdown()
}
