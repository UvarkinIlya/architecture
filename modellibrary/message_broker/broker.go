package message_broker

import (
	"encoding/json"
	"os"

	"github.com/nats-io/nats.go"

	"architecture/logger"
	"architecture/modellibrary"
)

const maxMessage = 10
const SubjMessages = "messages"

type Broker interface {
	PublishMessage(message modellibrary.Message) error
	Subscribe() (messageCh chan modellibrary.Message, err error)
}

type MessageBroker struct {
	subj       string
	connection *nats.Conn
}

func NewMessageBroker(subj string) (messageBroker *MessageBroker, err error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	connection, err := nats.Connect(url)
	if err != nil {
		return &MessageBroker{}, err
	}

	return &MessageBroker{
		subj:       subj,
		connection: connection,
	}, nil
}

func (m *MessageBroker) PublishMessage(message modellibrary.Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = m.connection.Publish(m.subj, messageBytes)
	if err != nil {
		return err
	}

	return nil
}

func (m *MessageBroker) Subscribe() (messageCh chan modellibrary.Message, err error) {
	messageCh = make(chan modellibrary.Message, maxMessage)

	_, err = m.connection.Subscribe(m.subj, func(msg *nats.Msg) {
		var message modellibrary.Message
		err := json.Unmarshal(msg.Data, &message)
		if err != nil {
			logger.Error("failed unmarshal message %+v due to error %s", msg, err)
			return
		}

		messageCh <- message
	})

	if err != nil {
		return messageCh, err
	}

	return messageCh, err
}
