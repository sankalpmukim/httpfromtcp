package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sankalpmukim/httpfromtcp/internal/request"
	"github.com/sankalpmukim/httpfromtcp/internal/response"
	"github.com/sankalpmukim/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server.ShuttingDown.Store(false)
	srv, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandleError {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandleError{StatusCode: response.BadRequest, Message: "Your problem is not my problem\n"}
		case "/myproblem":
			return &server.HandleError{StatusCode: response.InternalServerError, Message: "Woopsie, my bad\n"}
		default:
			w.Write([]byte("All good, frfr\n"))
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	server.ShuttingDown.Store(true)

	log.Println("Server gracefully stopped")
}
