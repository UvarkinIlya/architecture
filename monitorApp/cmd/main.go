package main

import (
	"architecture/logger"

	"architecture/monitorApp/app_manager"

	"architecture/modellibrary"
)

const logFile = "monitor.log"

func main() {
	logger.ConfigurateLogger(logFile)

	watchdogReq := modellibrary.WatchdogStartRequest{
		FileName:        "server_life",
		IntervalSeconds: 5,
	}
	watchdogChecker := app_manager.NewWatchdogChecker(watchdogReq, "http://localhost:8080/watchdog/start", 5, 10)

	manager := app_manager.NewManager(watchdogChecker)
	err := manager.Start()
	if err != nil {
		logger.Fatal("Failed start err: %s", err)
	}
}
