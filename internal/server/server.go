package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/sankalpmukim/httpfromtcp/internal/request"
	"github.com/sankalpmukim/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	netListener, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		return nil, err
	}
	serverInstance := Server{listener: netListener, handler: handler}
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
			fmt.Printf("Error accepting connections: %v", err)
			continue
		}
		fmt.Printf("A new connection has been accepted. %v\n", connection)

		go s.handle(connection)
	}
}

func (s *Server) handle(conn net.Conn) {
	request, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("Error in RequestFromReader", err)
	}

	responseWriter := response.NewResponseWriter(conn)
	s.handler(&responseWriter, request)

	conn.Close()
}

func HandleWritingError(w io.Writer, err HandleError) error {
	response.WriteStatusLine(w, err.StatusCode)
	outgoingMessage := fmt.Sprintf("An error occurred: %s", err.Message)
	headers := response.GetDefaultHeaders(len(outgoingMessage))
	response.WriteHeaders(w, headers)
	_, erro := w.Write([]byte(outgoingMessage))
	return erro
}

type HandleError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)
