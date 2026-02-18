package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
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
				toSendHeaders := (headers.ConvertInbuiltHeadersToOurHeaders(resp.Header))
				toSendHeaders["Connection"] = "close"

				if len(resp.TransferEncoding) != 0 {
					toSendHeaders["Transfer-Encoding"] = "chunked"
					toSendHeaders["Trailer"] = "X-Content-SHA256, X-Content-Length"
					delete(toSendHeaders, "content-length")
					w.WriteHeaders(toSendHeaders)
					fmt.Println("chunked encoding mode")
					hasher := sha256.New()
					var totalLength int64

					buf := make([]byte, 32*1024)

					for {
						n, err := resp.Body.Read(buf)
						if n > 0 {
							chunk := buf[:n]

							hasher.Write(chunk)
							totalLength += int64(n)

							w.WriteChunkedBody(chunk)
						}

						if err == io.EOF {
							break
						}
					}
					trailers := headers.NewHeaders()
					hashString := hex.EncodeToString(hasher.Sum(nil))

					trailers["X-Content-SHA256"] = hashString
					trailers["X-Content-Length"] = strconv.FormatInt(totalLength, 10)
					w.WriteTrailers(trailers)
				} else {
					w.WriteHeaders(toSendHeaders)
					fmt.Println("Not chunked encoding mode")
					body, _ := io.ReadAll(resp.Body)
					w.WriteBody(body)
				}
				fmt.Printf("Protocol: %s\n", resp.Proto)
			} else if req.RequestLine.RequestTarget == "/video" {
				videoFileContents, err := os.ReadFile("assets/vim.mp4")
				if errors.Is(err, os.ErrNotExist) {
					w.WriteStatusLine(response.NotFound)
					w.WriteHeaders(response.GetDefaultHeaders(0))
				}
				h := response.GetDefaultHeaders(len(videoFileContents))
				h["Content-Type"] = "video/mp4"
				w.WriteStatusLine(response.OK)
				w.WriteHeaders(h)
				w.WriteBody(videoFileContents)
			} else {
				// 	w.WriteStatusLine(response.OK)
				// 	h := response.GetDefaultHeaders(len(OkTemplate))
				// 	h["content-type"] = "text/html"
				// 	w.WriteHeaders(h)
				// 	w.WriteBody([]byte(OkTemplate))
				fileName := fmt.Sprintf("static-file-server%s", req.RequestLine.RequestTarget)
				if strings.HasSuffix(fileName, "/") {
					fileName += "index"
				}
				file, err := os.ReadFile(fileName)
				ext := filepath.Ext(req.RequestLine.RequestTarget)
				if errors.Is(err, os.ErrNotExist) {
					// try for html file
					file, err = os.ReadFile(fmt.Sprintf("%s.html", fileName))
				}
				if errors.Is(err, os.ErrNotExist) {
					w.WriteStatusLine(response.NotFound)
					w.WriteHeaders(response.GetDefaultHeaders(0))
				} else {
					if ext == "" {
						ext = "html"
					}
					mimeType := mime.TypeByExtension(ext)
					w.WriteStatusLine(response.OK)
					h := response.GetDefaultHeaders(len(file))
					h["Content-Type"] = mimeType
					w.WriteHeaders(h)
					w.WriteBody(file)
				}
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
