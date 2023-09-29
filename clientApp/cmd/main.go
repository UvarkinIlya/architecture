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

	client := socket_client.NewClient(os.Args[1])
	RepeatUtilSuccess(client.EstablishConnection, 10*time.Second)
	defer func() {
		err := client.CloseConnection()
		if err != nil {
			log.Println("Close connection err:", err)
			return
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Read err:", err)
			return
		}

		err = client.SendMessage(message)
		if err != nil {
			log.Println("Send message err:", err)
			return
		}
	}
}

func RepeatUtilSuccess(fn func() error, timeBetweenAttempts time.Duration) {
	err := fn()
	time.Sleep(timeBetweenAttempts)
	for err != nil {
		err = fn()
		time.Sleep(timeBetweenAttempts)
	}
}
