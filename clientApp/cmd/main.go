package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"architecture/modellibrary"

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

	auth(client, reader)
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

func auth(client *socket_client.ClientImpl, reader *bufio.Reader) {
	isAuthorised := false

	for !isAuthorised {
		fmt.Println("login:")
		login, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Read err: %s", err)
		}

		fmt.Println("password:")
		password, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Read err: %s", err)
		}

		user := modellibrary.User{
			Login:    strings.TrimSuffix(login, "\n"),
			Password: strings.TrimSuffix(password, "\n"),
		}

		isAuthorised, err = client.Auth(user)
		if err != nil {
			logger.Error("Auth err: %s", err)
		}

		if isAuthorised {
			logger.Info("Authorised for %s", login)
			fmt.Println("Authorised, pleas write messages:")
			break
		}

		logger.Error("Failed Authorised for %s", login)
		fmt.Println("Failed Authorised")
	}
}

func RepeatUtilSuccess(fn func() error, timeBetweenAttempts time.Duration, maxRepeatCount int) (err error) {
	repeatCount := 1
	err = fn()
	if errors.Is(err, socket_client.ErrPermissionsDenied) {
		fmt.Println(socket_client.PermissionsDenied)
		return err
	}

	for err != nil && repeatCount < maxRepeatCount {
		time.Sleep(timeBetweenAttempts)
		err = fn()
		repeatCount++
	}

	return err
}
