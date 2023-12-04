package storage

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"architecture/serverApp/common"
)

type MessageStorage interface {
	GetMessages() (messages []common.Message, err error)
	GetMessagesSince(since time.Time) (messages []common.Message, err error)
	SaveMessages(messages ...common.Message) (err error)
	SaveMessagesAndSort(messages ...common.Message) (err error)
}

func (s *StorageImpl) GetMessages() (messages []common.Message, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.getMessages()
}

func (s *StorageImpl) getMessages() (messages []common.Message, err error) {
	_, err = os.OpenFile(s.MessagesFilepath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return []common.Message{}, err
	}

	messagesBin, err := os.ReadFile(s.MessagesFilepath)
	if err != nil {
		return []common.Message{}, err
	}

	if len(messagesBin) == 0 {
		return []common.Message{}, nil
	}

	messages = make([]common.Message, 0)
	err = json.Unmarshal(messagesBin, &messages)
	if err != nil {
		return []common.Message{}, err
	}

	return messages, nil
}

func (s *StorageImpl) GetMessagesSince(since time.Time) (messages []common.Message, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

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

	bin, err := json.Marshal(append(messagesOld, messages...))
	if err != nil {
		return err
	}

	err = os.WriteFile(s.MessagesFilepath, bin, 0644)
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

	messagesBin, err := json.Marshal(removeDuplicate(messages))
	if err != nil {
		return err
	}

	return os.WriteFile(s.MessagesFilepath, messagesBin, 0644)
}

func removeDuplicate(messages []common.Message) (messageWithDuplicate []common.Message) {
	if len(messages) == 0 {
		return messages
	}

	messageWithDuplicate = make([]common.Message, 0, len(messages))

	last := messages[0]
	messageWithDuplicate = append(messageWithDuplicate, messages[0])
	for i := 1; i < len(messages); i++ {
		if last == messages[i] {
			continue
		}

		messageWithDuplicate = append(messageWithDuplicate, messages[i])
		last = messages[i]
	}

	return messageWithDuplicate
}
