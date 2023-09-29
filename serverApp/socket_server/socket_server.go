package socket_server

import (
	"fmt"
	"log"
	"net"
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
		log.Fatal(err)
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
			log.Println("Close connection error: ", err)
		}
	}(conn)

	log.Println("Handle new client")
	buf := make([]byte, BufferSize)

	_, err := conn.Write([]byte("connection established\n"))
	if err != nil {
		log.Println("Write to connection error: ", err)
	}

	for {
		readLen, err := conn.Read(buf)
		if err != nil {
			log.Println("Read from connection error: ", err)
			break
		}

		log.Print("Received a message:", string(buf[:readLen]))
	}

	_, err = conn.Write([]byte("connection closed\n"))
	if err != nil {
		log.Println("Write to connection error: ", err)
	}
}
