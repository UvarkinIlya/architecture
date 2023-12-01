package srv

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"architecture/logger"
	"architecture/modellibrary"
	"architecture/serverApp/message_manager"

	"architecture/serverApp/storage"
	"architecture/serverApp/sync_server"
	syncer2 "architecture/serverApp/syncer"
)

const (
	imagesDir    = "public/img"
	imageFeature = "[img]"
)

type Message struct {
	IsImg      bool
	TextOrPath string
}

type Server interface {
	Start()
}

type ServerImpl struct {
	messageManager      message_manager.MessageManager
	syncer              syncer2.Syncer
	syncServer          *sync_server.SyncServer
	db                  storage.Storage
	port                int
	distributedLockPort int
	isSynced            bool
}

func NewServer(messageManager message_manager.MessageManager, syncer syncer2.Syncer, syncServer *sync_server.SyncServer, db storage.Storage, port, distributedLockPort int) *ServerImpl {
	return &ServerImpl{
		messageManager:      messageManager,
		syncer:              syncer,
		syncServer:          syncServer,
		db:                  db,
		port:                port,
		distributedLockPort: distributedLockPort,
	}
}

func (s *ServerImpl) Start() {
	go s.sync()

	isLocked := s.tryDistributedLock()
	for !isLocked {
		time.Sleep(100 * time.Millisecond)
		isLocked = s.tryDistributedLock()
	}

	logger.Info("Get distributed lock")
	s.start()
}

func (s *ServerImpl) tryDistributedLock() (isLocked bool) {
	_, err := net.Listen("tcp", fmt.Sprintf(":%d", s.distributedLockPort))
	return err == nil
}

func (s *ServerImpl) sync() {
	now := time.Now()

	message, err := s.syncer.GetMessageSince(now)
	for err != nil {
		time.Sleep(500 * time.Millisecond)
		message, err = s.syncer.GetMessageSince(now)
	}

	err = s.db.SaveMessagesAndSort(message...)
	if err != nil {
		logger.Error("Failed save message due to error:%s", err)
		return
	}

	logger.Info("Message synced")
}

func (s *ServerImpl) start() {
	go s.syncServer.Start()

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	http.HandleFunc("/img/upload", s.uploadImage)
	http.HandleFunc("/ping", s.ping)
	http.HandleFunc("/watchdog/start", s.watchdogStart)
	http.HandleFunc("/messages", s.showMessages)
	http.HandleFunc("/message", s.receiveMessage)

	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil).Error())
}

func (s *ServerImpl) imageHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		s.getImage(writer, request)
	case http.MethodPost:
		s.uploadImage(writer, request)
	}
}

func (s *ServerImpl) getImage(writer http.ResponseWriter, request *http.Request) {
	path := pathWithoutPrefix(request.URL.String(), imagesDir)

	pathImage := struct {
		Path string
	}{
		Path: path,
	}

	t := template.Must(template.ParseFiles("serverApp/templates/image_page.html"))
	err := t.Execute(writer, pathImage)
	if err != nil {
		return
	}
}

func (s *ServerImpl) uploadImage(writer http.ResponseWriter, request *http.Request) {
	fileMultipart, handel, err := request.FormFile("img")
	if err != nil {
		return
	}

	defer fileMultipart.Close()

	imgDir := "public/img"
	file, err := os.Create(fmt.Sprintf("%s/%s", imgDir, handel.Filename))
	if err != nil {
		return
	}

	fileBytes, err := io.ReadAll(fileMultipart)
	if err != nil {
		fmt.Println(err)
	}

	_, err = file.Write(fileBytes)
	if err != nil {
		return
	}
}

func (s *ServerImpl) ping(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "pong")
}

func (s *ServerImpl) watchdogStart(writer http.ResponseWriter, request *http.Request) {
	watchdog := modellibrary.WatchdogStartRequest{}

	err := json.NewDecoder(request.Body).Decode(&watchdog)
	if err != nil {
		return
	}

	file, err := os.Create(watchdog.FileName)
	if err != nil {
		return
	}

	go startWatchdog(file, watchdog.IntervalSeconds)

	writer.WriteHeader(http.StatusOK)
}

func (s *ServerImpl) showMessages(writer http.ResponseWriter, request *http.Request) {
	t := template.Must(template.ParseFiles("templates/messages_page.html"))
	err := t.Execute(writer, s.getMessages())
	if err != nil {
		logger.Error("Template error: %s", err)
		return
	}
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

func (s *ServerImpl) getMessages() []Message {
	messages, err := s.db.GetMessages()
	if err != nil {
		logger.Error("Failed get message due to error:%s", err)
		return []Message{}
	}

	messagesTmpl := make([]Message, 0, len(messages))
	for _, message := range messages {
		messagesTmpl = append(messagesTmpl, Message{
			IsImg:      message.IsImg,
			TextOrPath: message.Text,
		})
	}

	return messagesTmpl
}

func (s *ServerImpl) receiveMessage(writer http.ResponseWriter, request *http.Request) {
	var message string

	err := json.NewDecoder(request.Body).Decode(&message)
	if err != nil {
		logger.Error("Failed receive message due to error:%s", err)
		return
	}

	logger.Info("Received message: %s", message)

	err = s.messageManager.AddMessages(message)
	if err != nil {
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func pathWithoutPrefix(utlPath string, prefix string) (path string) {
	return strings.TrimPrefix(utlPath, prefix)
}
