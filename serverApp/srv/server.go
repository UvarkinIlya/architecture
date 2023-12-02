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
	"architecture/serverApp/message_manager"
	"architecture/serverApp/sync_and_watchdog_server"
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
	messageManager        message_manager.MessageManager
	SyncAndWatchdogServer *sync_and_watchdog_server.SyncAndWatchdogServer
	port                  int
	distributedLockPort   int
	isSynced              bool
}

func NewServer(messageManager message_manager.MessageManager, SyncAndWatchdogServer *sync_and_watchdog_server.SyncAndWatchdogServer, port, distributedLockPort int) *ServerImpl {
	return &ServerImpl{
		messageManager:        messageManager,
		SyncAndWatchdogServer: SyncAndWatchdogServer,
		port:                  port,
		distributedLockPort:   distributedLockPort,
	}
}

func (s *ServerImpl) Start() {
	go s.SyncAndWatchdogServer.Start()
	go s.sync()

	isLocked := s.tryDistributedLock()
	for !isLocked {
		time.Sleep(500 * time.Millisecond)
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
	s.messageManager.SyncMessages()
	logger.Info("Messages synced")
}

func (s *ServerImpl) start() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	http.HandleFunc("/img/upload", s.uploadImage)
	http.HandleFunc("/ping", s.ping)
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

func (s *ServerImpl) showMessages(writer http.ResponseWriter, request *http.Request) {
	t := template.Must(template.ParseFiles("templates/messages_page.html"))
	err := t.Execute(writer, s.getMessages())
	if err != nil {
		logger.Error("Template error: %s", err)
		return
	}
}

func (s *ServerImpl) getMessages() []Message {
	messages, err := s.messageManager.GetMessages()
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
