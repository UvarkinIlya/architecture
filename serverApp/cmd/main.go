package main

import (
	"architecture/serverApp/socket_server"
	"architecture/serverApp/srv"
)

func main() {
	socketServer := socket_server.NewSocketServer(7070)
	go socketServer.Start()

	server := srv.NewServer(8080)
	server.Start()
}
