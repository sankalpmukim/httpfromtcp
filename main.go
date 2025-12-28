package main

import (
	"fmt"
	"io"
	"log"
	"os"
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
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error while opening messages.txt")
	}
	linesChannel := getLinesChannel(f)
	for v := range linesChannel {
		fmt.Printf("read: %s\n", v)
	}
	fmt.Println("read: end")
}
