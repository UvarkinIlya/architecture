package main

import (
	"fmt"

	"github.com/spf13/pflag"

	"architecture/logger"
	"architecture/modellibrary"
	"architecture/monitorApp/app_manager"
	"architecture/monitorApp/config"
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

	watchdogReq := modellibrary.WatchdogStartRequest{
		FileName:        cfg.Watchdog.Filename,
		IntervalSeconds: cfg.Watchdog.Interval,
	}

	watchdogChecker := app_manager.NewWatchdogChecker(watchdogReq, cfg.Watchdog.StartURL, cfg.Watchdog.Interval, cfg.Watchdog.MaxWait)

	manager := app_manager.NewManager(watchdogChecker, cfg.Server.ConfigPath, cfg.Server.BinPath)
	err = manager.Start()
	if err != nil {
		logger.Fatal("Failed start err: %s", err)
	}
}
