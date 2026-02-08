package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sankalpmukim/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server.ShuttingDown.Store(false)
	srv, err := server.Serve(port)
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
