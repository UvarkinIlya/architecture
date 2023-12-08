package main

import (
	"fmt"

	"github.com/spf13/pflag"

	"architecture/logger"
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

	viewerApp := viewer.NewViewer(cfg.Server.Port)
	err = viewerApp.Start()
	if err != nil {
		logger.Fatal("Failed start due to error: %s", err)
	}
}
