package sync_and_watchdog_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"architecture/logger"
	"architecture/modellibrary"

	"architecture/serverApp/storage"
)

type SyncAndWatchdogServer struct {
	db   storage.Storage
	port int
}

func NewSyncAndWatchdogServer(db storage.Storage, port int) *SyncAndWatchdogServer {
	return &SyncAndWatchdogServer{
		db:   db,
		port: port,
	}
}

func (s *SyncAndWatchdogServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/messages", s.SendMessages).Methods(http.MethodGet)
	r.HandleFunc("/messages", s.ReceiveMessages).Methods(http.MethodPost)
	r.HandleFunc("/watchdog/start", s.watchdogStart).Methods(http.MethodPost)

	logger.Info("Internal server started!")
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), r).Error())
}

func (s *SyncAndWatchdogServer) SendMessages(writer http.ResponseWriter, request *http.Request) {
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

func (s *SyncAndWatchdogServer) ReceiveMessages(writer http.ResponseWriter, request *http.Request) {
	messages := make([]modellibrary.Message, 0)
	err := json.NewDecoder(request.Body).Decode(&messages)
	if err != nil {
		logger.Info("Failed sync messages:%+v due to err:%s", messages, err)
		return
	}

	err = s.db.SaveMessages(messages...)
	if err != nil {
		return
	}

	logger.Info("Sync from neighbour messages:%+v", messages)

	writer.WriteHeader(http.StatusOK)
}

func (s *SyncAndWatchdogServer) watchdogStart(writer http.ResponseWriter, request *http.Request) {
	watchdog := modellibrary.WatchdogStartRequest{}

	err := json.NewDecoder(request.Body).Decode(&watchdog)
	if err != nil {
		logger.Error("Failed watchdog start due to err: %s", err)
		return
	}

	file, err := os.Create(watchdog.FileName)
	if err != nil {
		logger.Error("Failed watchdog start due to err: %s", err)
		return
	}

	go startWatchdog(file, watchdog.IntervalSeconds)

	logger.Info("Watchdog stated!")
	writer.WriteHeader(http.StatusOK)
}

func startWatchdog(file *os.File, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		err := file.Truncate(0)
		if err != nil {
			logger.Error("Failed watchdog due to error: %s", err.Error())
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			return
		}

		_, err = file.Write([]byte(time.Now().Format(time.RFC3339)))
		if err != nil {
			logger.Error("Failed watchdog due to error: %s", err.Error())
		}
	}
}
