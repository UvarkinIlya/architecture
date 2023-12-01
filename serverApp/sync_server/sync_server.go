package sync_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"architecture/logger"

	"architecture/serverApp/common"
	"architecture/serverApp/storage"
)

type SyncServer struct {
	db   storage.Storage
	port int
}

func NewSyncServer(db storage.Storage, port int) *SyncServer {
	return &SyncServer{
		db:   db,
		port: port,
	}
}

func (s *SyncServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/messages", s.SendMessages).Methods(http.MethodGet)
	r.HandleFunc("/messages", s.ReceiveMessages).Methods(http.MethodPost)

	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), r).Error())
}

func (s *SyncServer) SendMessages(writer http.ResponseWriter, request *http.Request) {
	unix, err := strconv.Atoi(request.URL.Query().Get("since"))
	if err != nil {
		logger.Error("Send message failed parse due to err:%s", err)
		return
	}
	since := time.Unix(int64(unix), 0)

	messages, err := s.db.GetMessagesSince(since)
	if err != nil {
		//logger.Error("Send message failed get message from db due to err:%s", err)
		return
	}

	bin, err := json.Marshal(messages)
	if err != nil {
		logger.Error("Send message failed marshal message due to err:%s", err)
		return
	}

	_, err = writer.Write(bin)
	if err != nil {
		logger.Error("Send message failed write message due to err:%s", err)
		return
	}
}

func (s *SyncServer) ReceiveMessages(writer http.ResponseWriter, request *http.Request) {
	messages := make([]common.Message, 0)
	err := json.NewDecoder(request.Body).Decode(&messages)
	if err != nil {
		return
	}

	err = s.db.SaveMessages(messages...)
	if err != nil {
		return
	}

	writer.WriteHeader(http.StatusOK)
}
