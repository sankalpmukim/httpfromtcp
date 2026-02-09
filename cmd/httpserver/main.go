package main

import (
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
	srv, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			w.WriteStatusLine(response.BadRequest)
			h := response.GetDefaultHeaders(len(BadRequestTemplate))
			w.WriteHeaders(h)
			w.WriteBody([]byte(BadRequestTemplate))

		case "/myproblem":
			w.WriteStatusLine(response.InternalServerError)
			h := response.GetDefaultHeaders(len(InternalServerErrorTemplate))
			w.WriteHeaders(h)
			w.WriteBody([]byte(InternalServerErrorTemplate))

		default:
			w.WriteStatusLine(response.OK)
			h := response.GetDefaultHeaders(len(OkTemplate))
			w.WriteHeaders(h)
			w.WriteBody([]byte(OkTemplate))
		}
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
