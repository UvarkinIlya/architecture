package socket_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"architecture/logger"
)

type Client interface {
	//EstablishConnection() error
	SendMessage(message string) error
	//CloseConnection() error
}

type ClientImpl struct {
	serverAddress string
	conn          net.Conn
	client        *http.Client
}

func NewClient(serverPort string) *ClientImpl {
	return &ClientImpl{
		serverAddress: fmt.Sprintf("http://localhost:%s", serverPort),
		client:        &http.Client{}}
}

//func (c *ClientImpl) EstablishConnection() (err error) {
//	conn, err := net.Dial("tcp", c.serverAddress)
//	if err != nil {
//		return
//	}
//
//	c.conn = conn
//	go c.readMessages(os.Stdout, conn)
//	return
//}

func (c *ClientImpl) SendMessage(message string) (err error) {
	if isFile(message) {
		if c.sendFile(message) != nil {
			return err
		}

		message = fmt.Sprintf("[img]%s", filepath.Base(message))
	}

	bin, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(fmt.Sprintf("%s/message", c.serverAddress), "application/json", bytes.NewBuffer(bin))
	if err != nil {
		return err
	}

	return nil
}

//func (c *ClientImpl) CloseConnection() (err error) {
//	return c.conn.Close()
//}

func (c *ClientImpl) readMessages(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		logger.Error("Copy err: %s", err)
	}
}

func (c *ClientImpl) sendFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileName := filepath.Base(path)
	formFile, err := writer.CreateFormFile("img", fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/img/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, err = c.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func isFile(message string) bool {
	_, err := os.Stat(message)
	if err == nil {
		return true
	}

	return !errors.Is(err, os.ErrNotExist)
}
