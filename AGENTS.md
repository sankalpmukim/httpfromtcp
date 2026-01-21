# HTTP From TCP - Developer Guide

This document provides instructions and guidelines for agents and developers working on the `httpfromtcp` repository.

## 1. Build, Lint, and Test Commands

### Testing
*   **Run all tests:**
    ```bash
    go test ./...
    ```
*   **Run a single test:**
    Use the `-run` flag with a regex matching the test name.
    ```bash
    go test -v -run ^TestRequestLineParse$ ./internal/request
    ```
*   **Run tests with coverage:**
    ```bash
    go test -cover ./...
    ```

### Linting & Formatting
*   **Format code:**
    ```bash
    gofmt -s -w .
    ```
*   **Lint code:**
    ```bash
    go vet ./...
    ```

### Building
*   **Build the TCP listener:**
    ```bash
    go build -o tcplistener ./cmd/tcplistener
    ```
*   **Build the UDP sender:**
    ```bash
    go build -o udpsender ./cmd/udpsender
    ```

## 2. Code Style Guidelines

### General
*   **Language:** Go (1.25.5+)
*   **Formatting:** strictly adhere to `gofmt` standards.
*   **Idioms:** Follow effective Go idioms. Keep functions short and focused.

### Naming Conventions
*   **Exported definitions:** Use **PascalCase** (e.g., `RequestFromReader`, `RequestLine`).
*   **Unexported definitions:** Use **camelCase** (e.g., `parseRequestLine`, `bufferSize`).
*   **Variables:** Use short, concise names for short scopes (e.g., `r` for request, `err` for error). Use descriptive names for larger scopes.
*   **Constants:** Use camelCase or PascalCase depending on export status (e.g., `requestStateInitialized`).

### Imports
Group imports into two blocks separated by a newline:
1.  Standard library imports.
2.  Third-party and internal project imports.

```go
import (
	"errors"
	"io"
	"strings"

	"github.com/sankalpmukim/httpfromtcp/internal/headers"
	"github.com/stretchr/testify/assert"
)
```

### Error Handling
*   **Explicit Checks:** Always check for errors immediately after a function call returns one.
    ```go
    if err != nil {
        return nil, err
    }
    ```
*   **Wrapping:** Use `fmt.Errorf("context: %w", err)` to wrap errors when adding context, or custom errors if specific handling is needed.
*   **Sentinel Errors:** Use `errors.Is(err, io.EOF)` for checking specific error types.
*   **No Panics:** Avoid `panic` in library code (`internal/`). Return errors instead. `log.Fatal` is acceptable in `main` packages.

### Type Safety
*   Avoid `interface{}` unless absolutely necessary.
*   Use `struct` tags for serialization if needed (though not heavily used currently).
*   Define custom types for domain concepts (e.g., `type RequestState int`).

### Testing Guidelines
*   **Framework:** Use `github.com/stretchr/testify`.
    *   Use `require` package for checks that must pass for the test to proceed (e.g., no error returned).
    *   Use `assert` package for value comparisons.
*   **Structure:**
    *   Place unit tests in the same package as the code (e.g., `request_test.go` in `package request`).
    *   Use table-driven tests for multiple inputs/outputs logic.
    *   Use helper structs/mocks (like `chunkReader`) to simulate network behavior.

### Project Structure
*   `cmd/`: Application entry points. `main` packages go here.
*   `internal/`: Library code that is private to this project. Imported by `cmd`.
