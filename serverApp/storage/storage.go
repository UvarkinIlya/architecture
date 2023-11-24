package storage

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"

	"architecture/serverApp/common"
)

type Storage interface {
	GetMessages() (messages []common.Message, err error)
	GetMessagesSince(since time.Time) (messages []common.Message, err error)
	SaveMessages(messages ...common.Message) (err error)
	SaveMessagesAndSort(messages ...common.Message) (err error)
}

type StorageImpl struct {
	MessagesFilepath string
	lock             sync.RWMutex
}

func NewStorageImpl(messageFilepath string) *StorageImpl {
	return &StorageImpl{
		MessagesFilepath: messageFilepath,
	}
}

func (s *StorageImpl) GetMessages() (messages []common.Message, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.getMessages()
}

func (s *StorageImpl) getMessages() (messages []common.Message, err error) {
	messagesBin, err := os.ReadFile(s.MessagesFilepath)
	if err != nil {
		return []common.Message{}, err
	}

	messages = make([]common.Message, 0)
	err = json.Unmarshal(messagesBin, &messages)
	if err != nil {
		return []common.Message{}, err
	}

	return messages, nil
}

func (s *StorageImpl) GetMessagesSince(since time.Time) (messages []common.Message, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	messages, err = s.getMessages()
	if err != nil {
		return []common.Message{}, err
	}

	for i := 0; i < len(messages); i++ {
		if messages[i].Time.Before(since) {
			continue
		}

		return messages[i:], nil
	}

	return []common.Message{}, nil
}

func (s *StorageImpl) SaveMessages(messages ...common.Message) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.saveMessages(messages...)
}

func (s *StorageImpl) saveMessages(messages ...common.Message) (err error) {
	messagesOld, err := s.getMessages()
	if err != nil {
		messagesOld = make([]common.Message, 0)
	}

	file, err := os.OpenFile(s.MessagesFilepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	bin, err := json.Marshal(append(messagesOld, messages...))
	if err != nil {
		return err
	}

	_, err = file.Write(bin)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageImpl) SaveMessagesAndSort(newMessages ...common.Message) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	messages, err := s.getMessages()
	if err != nil {
		return err
	}

	messages = append(messages, newMessages...)

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Time.Before(messages[j].Time)
	})

	messagesBin, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	return os.WriteFile(s.MessagesFilepath, messagesBin, 0644)
}
