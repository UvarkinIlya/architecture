package message_manager

import (
	"time"

	"architecture/logger"
	"architecture/modellibrary"
	"architecture/modellibrary/message_broker"
	"architecture/serverApp/storage"
	syncer2 "architecture/serverApp/syncer"
)

const (
	BufferSize = 10000
)

type MessageManager interface {
	AddMessages(...string) (err error)
	GetMessages() (messages []modellibrary.Message, err error)
	SyncMessages()
}

type MessageManagerImpl struct {
	db            storage.MessageStorage
	syncer        syncer2.Syncer
	messageBroker message_broker.Broker
}

func NewMessageManager(db storage.Storage, syncer syncer2.Syncer, messageBroker message_broker.Broker) *MessageManagerImpl {
	return &MessageManagerImpl{
		db:            db,
		syncer:        syncer,
		messageBroker: messageBroker,
	}
}

func (s *MessageManagerImpl) AddMessages(messages ...string) (err error) {
	messageDB := make([]modellibrary.Message, 0, len(messages))
	for _, msg := range messages {
		messageDB = append(messageDB, modellibrary.NewMessage(msg))
	}

	err = s.db.SaveMessages(messageDB...)
	if err != nil {
		logger.Error("Failed save message due to error: %s", err)
		return err
	}

	err = s.syncer.SendMessage(messageDB)
	if err != nil {
		logger.Error("Failed sync message due to error: %s", err)
	}

	s.publishMessages(messageDB)

	return nil
}

func (s *MessageManagerImpl) GetMessages() (messages []modellibrary.Message, err error) {
	messages, err = s.db.GetMessages()
	if err != nil {
		logger.Error("Failed get messages due to error: %s", err)
		return nil, err
	}

	return messages, nil
}

func (s *MessageManagerImpl) SyncMessages() {
	lastMessageTime, _ := s.lastMessageTime()

	message, err := s.syncer.GetMessageSince(lastMessageTime)
	for err != nil {
		time.Sleep(500 * time.Millisecond)
		message, err = s.syncer.GetMessageSince(lastMessageTime)
	}

	err = s.db.SaveMessagesAndSort(message...)
	if err != nil {
		logger.Error("Failed save message due to error:%s", err)
		return
	}
}

func (s *MessageManagerImpl) lastMessageTime() (message time.Time, err error) {
	messages, err := s.db.GetMessages()
	if err != nil {
		return time.Unix(0, 0), err
	}

	if len(messages) == 0 {
		return time.Unix(0, 0), nil
	}

	return messages[len(messages)-1].Time, nil
}

func (s *MessageManagerImpl) publishMessages(messages []modellibrary.Message) {
	for _, message := range messages {
		err := s.messageBroker.PublishMessage(message)
		if err != nil {
			logger.Error("Failed publish message %s due to error", message, err)
		}
	}
}
