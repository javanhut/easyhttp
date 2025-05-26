# EasyHTTP Project Structure

This document outlines the recommended project structure for the EasyHTTP Go library.

## Directory Layout

```
easyhttp/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI/CD pipeline
├── example/
│   └── main.go                 # Example usage of the library
├── .gitignore                  # Git ignore rules
├── .golangci.yml              # Linting configuration
├── LICENSE                     # License file (MIT recommended)
├── Makefile                   # Build automation
├── README.md                  # Main documentation
├── easyhttp.go                # Main library code
├── easyhttp_test.go           # Test suite
├── go.mod                     # Go module definition
└── go.sum                     # Go module checksums
```

## Core Files

### `easyhttp.go`
The main library implementation containing:
- `Response` struct with convenience methods
- `Client` struct with configuration options
- HTTP method functions (GET, POST, PUT, DELETE, etc.)
- Request option handling
- Authentication support

### `easyhttp_test.go`
Comprehensive test suite including:
- Unit tests for all public functions
- Integration tests with test HTTP server
- Benchmark tests for performance measurement
- Error handling tests
- Authentication tests

### `example/main.go`
Example application demonstrating:
- Basic HTTP requests
- JSON handling with easyjson
- Client configuration
- Authentication usage
- Error handling patterns

## Configuration Files

### `.golangci.yml`
Linting configuration for code quality:
```yaml
linters-settings:
  govet:
    check-shadowing: true
  gocyclo:
    min-complexity: 15
  misspell:
    locale: US

linters:
  enable:
    - deadcode
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - structcheck
    - typecheck
    - unused
    - varcheck
    - misspell
    - gocyclo

run:
  timeout: 5m
```

### `.gitignore`
Standard Go project ignores:
```
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
easyhttp-example*

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out
coverage.html

# Dependency directories
vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
```

## Deployment Structure

### For Library Users
When someone imports your library, they only need:
```go
import "github.com/javanhut/easyhttp"
```

The library will be used as:
```go
package main

import (
    "fmt"
    "log"
    "github.com/javanhut/easyhttp"
)

func main() {
    resp, err := easyhttp.GET("https://api.example.com/users")
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.OK() {
        jsonValue, err := resp.JSONValue()
        if err == nil {
            name := jsonValue.Get("name").AsString()
            fmt.Println("Name:", name)
        }
    }
}
```

### Versioning
Use semantic versioning with Git tags:
- `v1.0.0` - Initial stable release
- `v1.1.0` - New features (backward compatible)
- `v1.0.1` - Bug fixes
- `v2.0.0` - Breaking changes

Tag releases with:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Development Workflow

### Local Development
1. Clone the repository
2. Run `make dev-setup` to install development tools
3. Run `make deps` to install dependencies
4. Make changes
5. Run `make all` to format, lint, and test
6. Run `make example` to test with the example

### Testing
- `make test` - Run all tests
- `make test-coverage` - Generate coverage report
- `make bench` - Run performance benchmarks
- `make test-verbose` - Run tests with race detection

### CI/CD Pipeline
GitHub Actions automatically:
1. Tests on Go 1.20, 1.21, and 1.22
2. Runs linting and security checks
3. Generates coverage reports
4. Builds example application
5. Uploads benchmark results

## Publishing

### To GitHub
1. Push to GitHub repository
2. Create releases with proper tags
3. Include changelog in release notes

### To Go Module Registry
Once pushed to GitHub with proper tags, the module is automatically available via:
```bash
go get github.com/javanhut/easyhttp@latest
```

## Best Practices

### Code Organization
- Keep the main library in a single file for simplicity
- Separate tests into `_test.go` files
- Use clear, descriptive function and variable names
- Include comprehensive documentation comments

### Documentation
- Every exported function/type should have a doc comment
- Include usage examples in doc comments
- Maintain an up-to-date README with examples
- Document breaking changes in releases

### Testing
- Aim for >90% test coverage
- Test both success and error cases
- Include benchmark tests for performance-critical code
- Use table-driven tests where appropriate

### Versioning
- Follow semantic versioning strictly
- Tag all releases
- Maintain a CHANGELOG.md file
- Never force-push to main branch

This structure ensures your library is professional, maintainable, and easy for others to adopt and contribute to.
