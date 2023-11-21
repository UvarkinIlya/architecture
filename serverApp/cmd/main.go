package main

import (
	"fmt"

	"github.com/spf13/pflag"

	"architecture/logger"

	"architecture/serverApp/config"
	"architecture/serverApp/socket_server"
	"architecture/serverApp/srv"
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

	socketServer := socket_server.NewSocketServer(cfg.TCPSocket.Port)

	server := srv.NewServer(socketServer, cfg.HTTP.Port, cfg.DistributedLock.Port)
	server.Start()
}
