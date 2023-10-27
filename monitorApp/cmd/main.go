package main

import (
	"log"

	"architecture/monitorApp/app_manager"

	"architecture/modellibrary"
)

func main() {
	watchdogReq := modellibrary.WatchdogStartRequest{
		FileName:        "server_life",
		IntervalSeconds: 5,
	}
	watchdogChecker := app_manager.NewWatchdogChecker(watchdogReq, "http://localhost:8080/watchdog/start", 5, 10)

	manager := app_manager.NewManager(watchdogChecker)
	err := manager.Start()
	if err != nil {
		log.Fatal("Failed start err:", err)
	}
}
