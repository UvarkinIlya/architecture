package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"architecture/logger"

	"architecture/clientApp/socket_client"
)

const logFile = "client.log"

func main() {
	logger.ConfigurateLogger(logFile)

	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	client := socket_client.NewClient(os.Args[1])
	isEstablishConnection := RepeatUtilSuccess(client.EstablishConnection, 1*time.Second, 10)
	if !isEstablishConnection {
		logger.Fatal("Failed connection not establish")
	}
	defer func() {
		err := client.CloseConnection()
		if err != nil {
			logger.Error("Close connection err: %s", err)
			return
		}
	}()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Read err: %s", err)
			return
		}

		err = client.SendMessage(message)
		logger.Info("Send message: %s", message)
		if err != nil {
			logger.Error("Send message err: %s", err)
			return
		}
	}
}

func RepeatUtilSuccess(fn func() error, timeBetweenAttempts time.Duration, maxRepeatCount int) (isEstablish bool) {
	repeatCount := 1
	err := fn()
	for err != nil && repeatCount < maxRepeatCount {
		time.Sleep(timeBetweenAttempts)
		err = fn()
		repeatCount++
	}

	return err == nil
}
