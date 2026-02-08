package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/sankalpmukim/httpfromtcp/internal/request"
	"github.com/sankalpmukim/httpfromtcp/internal/utils"
)

type Server struct {
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	netListener, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		return nil, err
	}
	serverInstance := Server{listener: netListener}
	go serverInstance.listen()

	return &serverInstance, nil
}

func (s *Server) Close() error {
	fmt.Println("Server Close() called")
	return s.listener.Close()
}

var ShuttingDown atomic.Bool

func (s *Server) listen() error {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			if !ShuttingDown.Load() {
				return err
			}
		}
		fmt.Printf("A new connection has been accepted. %v\n", connection)

		go func() {
			request, err := request.RequestFromReader(connection)
			s.handle(connection)
			if err != nil {
				fmt.Println("Error in RequestFromReader", err)
			}

			if request != nil {
				utils.PrintRequest(*request)
			} else {
				fmt.Println("Request was nil")
			}
		}()
	}
}

func (s *Server) handle(conn net.Conn) {
	conn.Write([]byte(
		"HTTP/1.1 200 OK\n" +
			"Content-Type: text/plain\n" +
			"Content-Length: 13\n" +
			"\n" +
			"Hello World!\n"))
	conn.Close()
}
