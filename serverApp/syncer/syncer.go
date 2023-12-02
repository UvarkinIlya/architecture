package syncer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"architecture/logger"

	"architecture/serverApp/common"
)

const messagesUrl = "http://127.0.0.1:%d/messages"

type Syncer interface {
	GetMessageSince(since time.Time) (message []common.Message, err error)
	SendMessage(messages []common.Message) (err error)
}

type SyncerIpml struct {
	neighbourSyncAddr int
	newMessageCh      chan common.Message
}

func NewSyncerImpl(neighbourSyncAddr int) *SyncerIpml {
	syncer := &SyncerIpml{
		neighbourSyncAddr: neighbourSyncAddr,
		newMessageCh:      make(chan common.Message),
	}

	go syncer.syncMessages()

	return syncer
}

func (s *SyncerIpml) GetMessageSince(since time.Time) (message []common.Message, err error) {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/messages?since=%d", s.neighbourSyncAddr, since.Unix())) //TODO change to url params
	if err != nil {
		//		logger.Error("Failed get messages due to error:%s", err)
		return nil, err
	}

	message = make([]common.Message, 0)
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		//logger.Error("Failed get messages due to error:%s", err)
		return nil, err
	}

	return message, err
}

func (s *SyncerIpml) SendMessage(messages []common.Message) (err error) {
	for _, message := range messages {
		s.newMessageCh <- message
	}

	return nil
}

func (s *SyncerIpml) syncMessages() {
	messages := make([]common.Message, 0)
	messageCh := make(chan common.Message)

	go func() {
		for {
			messages = append(messages, <-s.newMessageCh)
			messageCh <- messages[0]
			messages = messages[1:]
		}
	}()

	for {
		msg := <-messageCh
		messageBin, err := json.Marshal([]common.Message{msg})
		if err != nil {
			logger.Error("syncMessages failed marshal due to error: %s", err)
			continue
		}

		messageBuffer := bytes.NewBuffer(messageBin)

		err = s.syncMessageBin(messageBuffer)
		if err != nil {
			logger.Error("Failed sync message %+v", msg)
			continue
		}
		//for err != nil {
		//	time.Sleep(500 * time.Millisecond)
		//	err = s.syncMessageBin(messageBuffer)
		//}
		logger.Info("Message %+v synced to neighbour", msg)
	}
}

func (s *SyncerIpml) syncMessageBin(messagesBin *bytes.Buffer) error {
	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/messages", s.neighbourSyncAddr), "application/json", messagesBin)
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Error("Failed send messages due to error:%s", err)
		return err
	}

	return nil
}
