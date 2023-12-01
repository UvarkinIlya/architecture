package message_manager

import (
	"architecture/logger"
	"architecture/serverApp/common"
	"architecture/serverApp/storage"
)

const (
	BufferSize = 10000
)

type MessageManager interface {
	AddMessages(...string) (err error)
	GetMessages() (messages []common.Message, err error)
}

type MessageManagerImpl struct {
	db storage.Storage
}

func NewMessageManager(db storage.Storage) *MessageManagerImpl {
	return &MessageManagerImpl{
		db: db,
	}
}

func (s *MessageManagerImpl) AddMessages(messages ...string) (err error) {
	messageDB := make([]common.Message, 0, len(messages))
	for _, msg := range messages {
		messageDB = append(messageDB, common.NewMessage(msg))
	}

	err = s.db.SaveMessages(messageDB...)
	if err != nil {
		logger.Error("Failed save message due to error: %s", err)
	}

	return nil
}

func (s *MessageManagerImpl) GetMessages() (messages []common.Message, err error) {
	messages, err = s.db.GetMessages()
	if err != nil {
		logger.Error("Failed get messages due to error: %s", err)
		return nil, err
	}

	return messages, nil
}
