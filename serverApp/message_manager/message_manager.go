package message_manager

import (
	"time"

	"architecture/logger"
	"architecture/serverApp/common"
	"architecture/serverApp/storage"
	syncer2 "architecture/serverApp/syncer"
)

const (
	BufferSize = 10000
)

type MessageManager interface {
	AddMessages(...string) (err error)
	GetMessages() (messages []common.Message, err error)
	SyncMessages()
}

type MessageManagerImpl struct {
	db     storage.MessageStorage
	syncer syncer2.Syncer
}

func NewMessageManager(db storage.Storage, syncer syncer2.Syncer) *MessageManagerImpl {
	return &MessageManagerImpl{
		db:     db,
		syncer: syncer,
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
		return err
	}

	err = s.syncer.SendMessage(messageDB)
	if err != nil {
		logger.Error("Failed sync message due to error: %s", err)
		return err
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
