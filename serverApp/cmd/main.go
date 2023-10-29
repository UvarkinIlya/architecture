package main

import (
	"architecture/logger"

	"architecture/serverApp/socket_server"
	"architecture/serverApp/srv"
)

const logFile = "server.log"

func main() {
	logger.ConfigurateLogger(logFile)

	socketServer := socket_server.NewSocketServer(7070)

	server := srv.NewServer(socketServer, 8080)
	server.Start()
}
