package storage

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"architecture/modellibrary"
)

type MessageStorage interface {
	GetMessages() (messages []modellibrary.Message, err error)
	GetMessagesSince(since time.Time) (messages []modellibrary.Message, err error)
	SaveMessages(messages ...modellibrary.Message) (err error)
	SaveMessagesAndSort(messages ...modellibrary.Message) (err error)
}

func (s *StorageImpl) GetMessages() (messages []modellibrary.Message, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.getMessages()
}

func (s *StorageImpl) getMessages() (messages []modellibrary.Message, err error) {
	_, err = os.OpenFile(s.MessagesFilepath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return []modellibrary.Message{}, err
	}

	messagesBin, err := os.ReadFile(s.MessagesFilepath)
	if err != nil {
		return []modellibrary.Message{}, err
	}

	if len(messagesBin) == 0 {
		return []modellibrary.Message{}, nil
	}

	messages = make([]modellibrary.Message, 0)
	err = json.Unmarshal(messagesBin, &messages)
	if err != nil {
		return []modellibrary.Message{}, err
	}

	return messages, nil
}

func (s *StorageImpl) GetMessagesSince(since time.Time) (messages []modellibrary.Message, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	messages, err = s.getMessages()
	if err != nil {
		return []modellibrary.Message{}, err
	}

	for i := 0; i < len(messages); i++ {
		if messages[i].Time.Before(since) {
			continue
		}

		return messages[i:], nil
	}

	return []modellibrary.Message{}, nil
}

func (s *StorageImpl) SaveMessages(messages ...modellibrary.Message) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.saveMessages(messages...)
}

func (s *StorageImpl) saveMessages(messages ...modellibrary.Message) (err error) {
	messagesOld, err := s.getMessages()
	if err != nil {
		messagesOld = make([]modellibrary.Message, 0)
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

func (s *StorageImpl) SaveMessagesAndSort(newMessages ...modellibrary.Message) (err error) {
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

func removeDuplicate(messages []modellibrary.Message) (messageWithDuplicate []modellibrary.Message) {
	if len(messages) == 0 {
		return messages
	}

	messageWithDuplicate = make([]modellibrary.Message, 0, len(messages))

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
