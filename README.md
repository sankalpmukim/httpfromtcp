# HTTP from TCP

A learning project to build an HTTP server from scratch on top of TCP sockets in Go. This project explores the inner workings of the HTTP protocol by implementing request parsing, header handling, and response generation at the TCP level.

## Getting Started

### Prerequisites

- Go 1.25.5 or later

### Running Tests

Run all tests to verify the implementation:

```bash
go test -v ./internal/request/
go test ./...
```

### Starting the Server

Start the HTTP server on port 42069:

```bash
go run ./cmd/httpserver/
```

Press `Ctrl+C` to gracefully stop the server.

### Testing the Server

Once the server is running, test it with curl:

**POST request with JSON body:**

```bash
curl -X POST http://localhost:42069/coffee \
    -H 'Content-Type: application/json' \
    -d '{"type": "dark mode", "size": "medium"}'
```

**GET request:**

```bash
curl http://localhost:42069/use-neovim-btw
```

Both requests should receive a `200 OK` response with the body `Hello World!`.

## Project Structure

- `cmd/httpserver/` - Main server entry point
- `internal/request/` - HTTP request parsing logic
- `internal/headers/` - HTTP header parsing and handling
- `internal/server/` - TCP listener and connection handling
- `internal/utils/` - Utility functions

## Learning Outcomes

This project demonstrates:
- TCP socket programming in Go
- HTTP protocol specification (RFC 7230)
- Request line parsing
- Header parsing
- State machine pattern for parsing
- Concurrent connection handling with goroutines
- Graceful server shutdown

## Acknowledgments

This project was inspired by and built following the excellent guidance from the [Boot.dev Learn HTTP Protocol with Go course](https://www.boot.dev/courses/learn-http-protocol-golang). A huge thank you to the Boot.dev team for making their high-quality learning material available for free.

**Note:** No AI was used in the coding of this project. I did not purchase a subscription, so I did not have access to the solution files for any of the lessons. All code was written independently based on the course instructions and my understanding of the concepts.

If you're interested in learning how the HTTP protocol works under the hood or want to build projects like this yourself, I highly encourage you to check out the [Boot.dev course](https://www.boot.dev/courses/learn-http-protocol-golang). It's a fantastic resource for understanding web protocols from first principles.
