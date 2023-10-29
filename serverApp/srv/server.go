package srv

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"architecture/logger"
	"architecture/modellibrary"

	"architecture/serverApp/socket_server"
)

const imageURLPrefix = "/images/"

type Server interface {
	Start()
}

type ServerImpl struct {
	socketServer socket_server.SocketServer
	port         int
}

func NewServer(socketServer socket_server.SocketServer, port int) *ServerImpl {
	return &ServerImpl{
		socketServer: socketServer,
		port:         port,
	}
}

func (server *ServerImpl) Start() {
	go server.socketServer.Start()

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("serverApp/public/"))))
	http.HandleFunc("/img/upload", server.uploadImage)
	http.HandleFunc("/ping", server.ping)
	http.HandleFunc("/watchdog/start", server.watchdogStart)
	http.HandleFunc("messages", server.showMessages)

	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil).Error())
}

func (server *ServerImpl) imageHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		server.getImage(writer, request)
	case http.MethodPost:
		server.uploadImage(writer, request)
	}
}

func (server *ServerImpl) getImage(writer http.ResponseWriter, request *http.Request) {
	path := pathWithoutPrefix(request.URL.String(), imageURLPrefix)

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

func (server *ServerImpl) uploadImage(writer http.ResponseWriter, request *http.Request) {
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

func (server *ServerImpl) ping(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "pong")
}

func (server *ServerImpl) watchdogStart(writer http.ResponseWriter, request *http.Request) {
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

func (server *ServerImpl) showMessages(writer http.ResponseWriter, request *http.Request) {
	imagesPath := pathWithoutPrefix(request.URL.String(), imageURLPrefix)

	pathImage := struct {
		Path string
	}{
		Path: imagesPath,
	}

	t := template.Must(template.ParseFiles("serverApp/templates/image_page.html"))
	err := t.Execute(writer, pathImage)
	if err != nil {
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

func pathWithoutPrefix(utlPath string, prefix string) (path string) {
	return strings.TrimPrefix(utlPath, prefix)
}
