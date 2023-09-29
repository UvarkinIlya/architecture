package main

import (
	"log"
	"time"

	"architecture/monitorApp/app_manager"
)

func main() {
	manager := app_manager.NewManager("serverApp", "http://localhost:8080/ping", 5*time.Second)
	err := manager.Start()
	if err != nil {
		log.Fatal("Failed start err:", err)
	}

	time.Sleep(1 * time.Second)
	log.Println("Start check")
	manager.Check()
}
