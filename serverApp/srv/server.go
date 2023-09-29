package srv

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

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

func pathWithoutPrefix(utlPath string, prefix string) (path string) {
	return strings.TrimPrefix(utlPath, prefix)
}
