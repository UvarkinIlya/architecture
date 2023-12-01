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
	SendMessage(message []common.Message) (err error)
}

type SyncerIpml struct {
	neighbourSyncAddr int
}

func NewSyncerIpml(neighbourSyncAddr int) *SyncerIpml {
	return &SyncerIpml{neighbourSyncAddr: neighbourSyncAddr}
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

func (s *SyncerIpml) SendMessage(message []common.Message) (err error) {
	messagesBin, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/messages", s.neighbourSyncAddr), "application/json", bytes.NewBuffer(messagesBin))
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Error("Failed send messages due to error:%s", err)
		return err
	}

	return nil
}
