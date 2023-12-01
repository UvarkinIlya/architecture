package main

import (
	"fmt"

	"github.com/spf13/pflag"

	"architecture/logger"
	"architecture/serverApp/message_manager"

	"architecture/serverApp/config"
	"architecture/serverApp/srv"
	"architecture/serverApp/storage"
	"architecture/serverApp/sync_server"
	syncer2 "architecture/serverApp/syncer"
)

const logFile = "server.log"

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

	db := storage.NewStorageImpl(cfg.Storage.MessageFilePath)
	syncer := syncer2.NewSyncerIpml(cfg.Syncer.Port) //TODO get addr from config
	syncServer := sync_server.NewSyncServer(db, cfg.Syncer.Port)

	messageManager := message_manager.NewMessageManager(db)

	server := srv.NewServer(messageManager, syncer, syncServer, db, cfg.HTTP.Port, cfg.DistributedLock.Port)
	server.Start()
}
