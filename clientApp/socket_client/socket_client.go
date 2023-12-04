package socket_client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"architecture/logger"
	"architecture/modellibrary"
)

const (
	permissionsSendMessage = "sendMessage"
	permissionsSendImg     = "sendImg"
	PermissionsDenied      = "permissions denied"
)

var (
	ErrPermissionsDenied = errors.New(PermissionsDenied)
)

type Client interface {
	Auth(login string) error
	SendMessage(message string) error
}

type ClientImpl struct {
	serverAddress string
	client        *http.Client
	permissions   map[string]struct{}
}

func NewClient(serverPort string) *ClientImpl {
	return &ClientImpl{
		serverAddress: fmt.Sprintf("http://localhost:%s", serverPort),
		client:        &http.Client{}}
}

func (c *ClientImpl) SendMessage(message string) (err error) {
	//TODO rename
	if !c.checkPermissions(message) {
		return ErrPermissionsDenied
	}

	if isFile(message) {
		if c.sendFile(message) != nil {
			return err
		}

		message = fmt.Sprintf("[img]%s", filepath.Base(message))
	}

	return c.sendMessage(message)
}

func (c *ClientImpl) checkPermissions(message string) (ok bool) {
	if isFile(message) {
		_, ok = c.permissions[permissionsSendImg]
		return ok
	}

	_, ok = c.permissions[permissionsSendMessage]
	return ok
}

func (c *ClientImpl) Auth(user modellibrary.User) (isAuthorised bool, err error) {
	userWithActions, statusCode, err := c.login(user)
	if err != nil {
		return false, err
	}

	if statusCode == http.StatusOK {
		c.permissions = userWithActions.Actions
	}

	return statusCode == http.StatusOK, err
}

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

func (c *ClientImpl) sendMessage(message string) (err error) {
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

func (c *ClientImpl) login(user modellibrary.User) (userWithActions modellibrary.UserWithActions, statusCode int, err error) {
	bin, err := json.Marshal(user)
	if err != nil {
		return modellibrary.UserWithActions{}, http.StatusInternalServerError, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/auth", c.serverAddress), "application/json", bytes.NewBuffer(bin))
	if err != nil {
		return modellibrary.UserWithActions{}, http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	token := string(body)
	userWithActionsBin, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return modellibrary.UserWithActions{}, 0, err
	}

	err = json.Unmarshal(userWithActionsBin, &userWithActions)
	if err != nil {
		return modellibrary.UserWithActions{}, 0, err
	}

	return userWithActions, resp.StatusCode, nil
}
