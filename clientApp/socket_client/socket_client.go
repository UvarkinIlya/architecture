package socket_client

import (
	"io"
	"net"
	"os"

	"architecture/logger"
)

type Client interface {
	EstablishConnection() error
	SendMessage(message string) error
	CloseConnection() error
}

type ClientImpl struct {
	serverAddress string
	conn          net.Conn
}

func NewClient(serverAddress string) *ClientImpl {
	return &ClientImpl{serverAddress: serverAddress}
}

func (c *ClientImpl) EstablishConnection() (err error) {
	conn, err := net.Dial("tcp", c.serverAddress)
	if err != nil {
		return
	}

	c.conn = conn
	go c.readMessages(os.Stdout, conn)
	return
}

func (c *ClientImpl) SendMessage(message string) (err error) {
	_, err = c.conn.Write([]byte(message))
	return err
}

func (c *ClientImpl) CloseConnection() (err error) {
	return c.conn.Close()
}

func (c *ClientImpl) readMessages(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		logger.Error("Copy err: %s", err)
	}
}
