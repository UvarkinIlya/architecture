package main

import (
	"fmt"

	"github.com/spf13/pflag"

	"architecture/logger"
	"architecture/modellibrary/message_broker"
	"architecture/viewerApp/config"
	"architecture/viewerApp/viewer"
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

	messageBroker, err := message_broker.NewMessageBroker(message_broker.SubjMessages)
	if err != nil {
		println("Failed start message broker failed due to error: ", err)
		logger.Fatal("Failed start message broker failed due to error: %s", err)
	}

	viewerApp := viewer.NewViewer(messageBroker, cfg.Server.Port)
	err = viewerApp.Start()
	if err != nil {
		println("Failed start Viewer due to error: ", err.Error())
		logger.Fatal("Failed start Viewer due to error: %s", err)
	}
}
