package socket_server

import (
	"fmt"
	"net"
	"strings"

	"architecture/logger"
)

const BufferSize = 10000

type SocketServer interface {
	Start()
}

type SocketServerImpl struct {
	port int
}

func NewSocketServer(port int) *SocketServerImpl {
	return &SocketServerImpl{
		port: port,
	}
}

func (s SocketServerImpl) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		logger.Fatal(err.Error())
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			continue //TODO add loger
		}

		go s.handlerClient(conn)
	}

}

func (s SocketServerImpl) handlerClient(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error("Close connection error: %s", err)
		}
	}(conn)

	logger.Info("Handle new client")
	buf := make([]byte, BufferSize)

	_, err := conn.Write([]byte("connection established\n"))
	if err != nil {
		logger.Error("Write to connection error: %s", err)
	}

	for {
		readLen, err := conn.Read(buf)
		if err != nil {
			logger.Error("Read from connection error: %s", err)
			break
		}

		messages := strings.Split(string(buf[:readLen]), "\n")

		for _, message := range messages {
			if message == "" {
				continue
			}
			logger.Info("Received a message:", message)
		}
	}

	_, err = conn.Write([]byte("connection closed\n"))
	if err != nil {
		logger.Error("Write to connection error: %s", err)
	}
}
