package server

import (
	"fmt"
	"net"

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

func (s *Server) listen() error {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			return err
		}
		fmt.Printf("A new connection has been accepted. %v\n", connection)

		request, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Println("Error in RequestFromReader")
		}

		utils.PrintRequest(*request)
	}
}
