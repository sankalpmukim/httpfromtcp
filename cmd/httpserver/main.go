package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sankalpmukim/httpfromtcp/internal/headers"
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
			h["content-type"] = "text/html"
			w.WriteHeaders(h)
			w.WriteBody([]byte(BadRequestTemplate))

		case "/myproblem":
			w.WriteStatusLine(response.InternalServerError)
			h := response.GetDefaultHeaders(len(InternalServerErrorTemplate))
			h["content-type"] = "text/html"
			w.WriteHeaders(h)
			w.WriteBody([]byte(InternalServerErrorTemplate))

		default:
			if after, ok := strings.CutPrefix(req.RequestLine.RequestTarget, "/httpbin"); ok {
				httpBinTarget := fmt.Sprintf("https://httpbin.org%s", after)
				var reqBody io.Reader = nil
				if req.Headers.Get("content-length") != "" {
					reqBody = bytes.NewReader(req.Body)
				}
				binRequest, _ := http.NewRequest(req.RequestLine.Method, httpBinTarget, reqBody)
				for k, v := range req.Headers {
					binRequest.Header.Set(k, v)
				}
				tr := &http.Transport{
					TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper), // Disable HTTP/2
				}

				// Create a custom Client using the custom Transport
				client := &http.Client{Transport: tr}

				resp, _ := client.Do(binRequest)
				w.WriteStatusLine(response.StatusCode(resp.StatusCode))
				w.WriteHeaders(headers.ConvertInbuiltHeadersToOurHeaders(resp.Header))
				// NOTE: First part of this if condition will never be true because of
				// how go's http client library works.
				// TODO: Implement it using resp.TransferEncoding
				if resp.Header.Get("transfer-encoding") != "" || strings.Contains(httpBinTarget, "stream") {
					fmt.Println("chunked encoding mode")
					scanner := bufio.NewScanner(resp.Body)
					for scanner.Scan() {
						line := scanner.Bytes()
						w.WriteChunkedBody(line)
					}
					w.WriteChunkedBodyDone()
				} else {
					fmt.Println("Not chunked encoding mode")
					body, _ := io.ReadAll(resp.Body)
					w.WriteBody(body)
				}
				fmt.Printf("Protocol: %s\n", resp.Proto)
			} else {
				w.WriteStatusLine(response.OK)
				h := response.GetDefaultHeaders(len(OkTemplate))
				h["content-type"] = "text/html"
				w.WriteHeaders(h)
				w.WriteBody([]byte(OkTemplate))
			}
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
