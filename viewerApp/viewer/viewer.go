package viewer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"architecture/logger"
	"architecture/modellibrary"
	"architecture/modellibrary/message_broker"
)

type Viewer struct {
	messageBroker message_broker.Broker
	serverPort    int
}

func NewViewer(messageBroker message_broker.Broker, serverPort int) *Viewer {
	return &Viewer{
		messageBroker: messageBroker,
		serverPort:    serverPort,
	}
}

func (v *Viewer) Start() error {
	messages, err := v.getMessages()
	if err != nil {
		return err
	}

	printMessages(messages...)

	messageCh, err := v.messageBroker.Subscribe()
	if err != nil {
		println("Failed subscribe to messages due to err: ", err.Error())
		logger.Fatal("Failed subscribe to messages due to err: ", err)
	}

	for {
		message := <-messageCh
		logger.Info("Receive new message: %s", message.Text)
		printMessages(message)
	}
}

func (v *Viewer) getMessages() (messages []modellibrary.Message, err error) {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/int/messages", v.serverPort))
	if err != nil {
		return nil, err
	}

	messages = make([]modellibrary.Message, 0)
	err = json.NewDecoder(resp.Body).Decode(&messages)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func printMessages(messages ...modellibrary.Message) {
	for _, message := range messages {
		if message.IsImg {
			fmt.Printf("%s%s\n", modellibrary.ImageFeature, message.Text)
			return
		}

		fmt.Printf("%s\n", message.Text)
	}
}
