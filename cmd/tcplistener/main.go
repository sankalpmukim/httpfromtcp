package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sankalpmukim/httpfromtcp/internal/request"
)

func main() {
	port := "42069"
	netListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("error while listening on port " + port)
	}
	defer netListener.Close()
	for {
		connection, err := netListener.Accept()
		if err != nil {
			log.Fatal("Error while accepting connections")
			break
		}
		fmt.Printf("A Connection has been accepted. %v\n", connection)

		request, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Println("Error in RequestFromReader")
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", request.RequestLine.Method)
		fmt.Printf("- Target: %v\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range request.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(request.Body))

		fmt.Println("The HTTP request has been parsed.")
	}

}
