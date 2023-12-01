package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"architecture/clientApp/socket_client"
	"architecture/logger"
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
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Read err: %s", err)
			continue
		}

		message = strings.TrimSuffix(message, "\n")
		sendMessage := func() error {
			return client.SendMessage(message)
		}

		err = RepeatUtilSuccess(sendMessage, 1*time.Second, 10)
		if err != nil {
			logger.Error("Failed send message: %s due to error: %s", message, err)
			continue
		}

		logger.Info("Send message: %s", message)
	}
}

func RepeatUtilSuccess(fn func() error, timeBetweenAttempts time.Duration, maxRepeatCount int) (err error) {
	repeatCount := 1
	err = fn()
	for err != nil && repeatCount < maxRepeatCount {
		time.Sleep(timeBetweenAttempts)
		err = fn()
		repeatCount++
	}

	return err
}
