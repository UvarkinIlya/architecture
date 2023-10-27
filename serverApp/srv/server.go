package srv

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"architecture/modellibrary"
	"architecture/serverApp/image_manager"
)

const imageURLPrefix = "/images/"

type Server interface {
	Start()
}

type ServerImpl struct {
	port         int
	imageManager image_manager.ImageManager
}

func NewServer(port int) *ServerImpl {
	return &ServerImpl{
		port: port,
		//imageManager: imageManager,
	}
}

func (server *ServerImpl) Start() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("serverApp/public/"))))

	http.HandleFunc(imageURLPrefix, server.imageHandler)
	http.HandleFunc("/ping", server.ping)
	http.HandleFunc("/watchdog/start", server.watchdogStart)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil))
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

func startWatchdog(file *os.File, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		err := file.Truncate(0)
		if err != nil {
			log.Println("Failed watchdog due to error: ", err.Error())
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			return
		}

		_, err = file.Write([]byte(time.Now().Format(time.RFC3339)))
		if err != nil {
			log.Println("Failed watchdog due to error: ", err.Error())
		}
	}
}

func pathWithoutPrefix(utlPath string, prefix string) (path string) {
	return strings.TrimPrefix(utlPath, prefix)
}
