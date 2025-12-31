package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChannel := make(chan string)
	go func() {
		defer close(linesChannel)
		defer f.Close()

		line := ""
		for {
			bytes := make([]byte, 8)
			_, err := f.Read(bytes)
			if err == io.EOF {
				break
			}
			readContent := string(bytes)
			indexOfNewLine := strings.Index(readContent, "\n")
			if indexOfNewLine == -1 {
				line += readContent
			} else {
				line += readContent[:indexOfNewLine]

				linesChannel <- line

				line = readContent[indexOfNewLine+1:]
			}
		}
	}()
	return linesChannel
}

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

		linesChannel := getLinesChannel(connection)
		for v := range linesChannel {
			fmt.Printf("read: %s\n", v)
		}
		fmt.Println("The connection has been closed.")
	}

}
