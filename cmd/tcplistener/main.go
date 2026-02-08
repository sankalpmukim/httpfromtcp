package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sankalpmukim/httpfromtcp/internal/request"
	"github.com/sankalpmukim/httpfromtcp/internal/utils"
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

		utils.PrintRequest(*request)
	}

}
