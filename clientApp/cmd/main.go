package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"architecture/clientApp/socket_client"
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	client := socket_client.NewClient(os.Args[1])
	isEstablishConnection := RepeatUtilSuccess(client.EstablishConnection, 1*time.Second, 10)
	if !isEstablishConnection {
		log.Panic("Failed connection not establish")
	}
	defer func() {
		err := client.CloseConnection()
		if err != nil {
			log.Println("Close connection err:", err)
			return
		}
	}()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Read err:", err)
			return
		}

		err = client.SendMessage(message)
		log.Println("Send message:", message)
		if err != nil {
			log.Println("Send message err:", err)
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
