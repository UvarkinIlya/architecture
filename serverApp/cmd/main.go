package main

import (
	"architecture/serverApp/srv"
)

func main() {

	server := srv.NewServer(8080)

	server.Start()
}
