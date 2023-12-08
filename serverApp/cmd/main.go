package main

import (
	"fmt"

	"github.com/spf13/pflag"

	"architecture/logger"
	"architecture/modellibrary/message_broker"
	"architecture/serverApp/auth"
	"architecture/serverApp/message_manager"

	"architecture/serverApp/config"
	"architecture/serverApp/srv"
	"architecture/serverApp/storage"
	"architecture/serverApp/sync_and_watchdog_server"
	syncer2 "architecture/serverApp/syncer"
)

func main() {
	configPath := pflag.StringP("config", "c", "", "Config file path")
	showHelp := pflag.BoolP("help", "h", false,
		"Show help message")

	pflag.Parse()
	if *showHelp {
		pflag.Usage()
		return
	}

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed parse config due to error %s", err))
	}

	logger.ConfigurateLogger(cfg.Logger.Filename)

	db := storage.NewStorageImpl(cfg.Storage.MessageFilePath, cfg.Storage.UsersFilePath)
	syncer := syncer2.NewSyncerImpl(cfg.Neighbour.Syncer.Port) //TODO get addr from config
	SyncAndWatchdogServer := sync_and_watchdog_server.NewSyncAndWatchdogServer(db, cfg.Syncer.Port)

	messageBroker, err := message_broker.NewMessageBroker(message_broker.SubjMessages)
	if err != nil {
		logger.Fatal("Failed start message broker failed due to error: %s", err)
	}

	messageManager := message_manager.NewMessageManager(db, syncer, messageBroker)
	auther := auth.NewAutherImpl(db)

	server := srv.NewServer(auther, messageManager, SyncAndWatchdogServer, cfg.HTTP.Port, cfg.DistributedLock.Port)
	server.Start()
}
