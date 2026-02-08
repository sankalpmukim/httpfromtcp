package utils

import (
	"fmt"

	"github.com/sankalpmukim/httpfromtcp/internal/request"
)

func PrintRequest(request request.Request) {
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
