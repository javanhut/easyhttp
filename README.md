# EasyHTTP - Go HTTP Client Library

A simple, powerful HTTP client library for Go that combines the performance of Go's native HTTP client with the ease of use of Python's requests library.

## Features

- **Simple API**: Intuitive method chaining and easy-to-use functions
- **Python-like**: Familiar API for those coming from Python's requests
- **High Performance**: Built on Go's native `net/http` package with easyjson for flexible JSON processing
- **JSON Support**: Built-in JSON marshaling/unmarshaling using github.com/javanhut/easyjson with fluent query API
- **Authentication**: Support for Bearer tokens and Basic auth
- **Flexible**: Support for custom headers, timeouts, and request options
- **Response Helpers**: Easy methods to get text, JSON, or bytes from responses
- **Base URL Support**: Set a base URL for all requests
- **Query Parameters**: Easy parameter handling
- **Error Handling**: Clean error handling with Go idioms

## Installation

Add EasyHTTP to your Go project:

```bash
go get github.com/javanhut/easyhttp
```

The library also requires the easyjson dependency:
```bash
go get github.com/javanhut/easyjson
```

## Quick Start

### Simple GET Request

```go
package main

import (
    "fmt"
    "log"
    "github.com/javanhut/easyhttp"
)

func main() {
    resp, err := easyhttp.GET("https://httpbin.org/get")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Status:", resp.StatusCode)
    fmt.Println("Body:", resp.Text())
}
```

### POST with easyjson JSONValue

```go
// Create JSON using easyjson's fluent API
jsonObj := easyjson.NewObject()
jsonObj.Set("name", "John Doe")
jsonObj.Set("email", "john@example.com")

resp, err := easyhttp.POST("https://httpbin.org/post", &easyhttp.RequestOptions{
    JSON: jsonObj, // Pass easyjson JSONValue directly
})
```

### Using a Client with Base URL

```go
client := easyhttp.New().
    SetBaseURL("https://api.github.com").
    SetHeaders(map[string]string{
        "User-Agent": "MyApp/1.0",
    })

resp, err := client.GET("/user") // Requests https://api.github.com/user
```

## API Reference

### Creating Clients

#### Quick Functions
- `easyhttp.GET(url, opts...)` - Simple GET request
- `easyhttp.POST(url, opts...)` - Simple POST request
- `easyhttp.PUT(url, opts...)` - Simple PUT request
- `easyhttp.DELETE(url, opts...)` - Simple DELETE request
- `easyhttp.PATCH(url, opts...)` - Simple PATCH request
- `easyhttp.HEAD(url, opts...)` - Simple HEAD request

#### Client Creation
```go
client := easyhttp.New()
```

#### Client Configuration
```go
client := easyhttp.New().
    SetBaseURL("https://api.example.com").
    SetTimeout(30 * time.Second).
    SetHeaders(map[string]string{
        "User-Agent": "MyApp/1.0",
    }).
    SetAuth(&easyhttp.Auth{
        Token: "bearer-token",
    })
```

### Request Options

```go
type RequestOptions struct {
    Headers       map[string]string // Custom headers
    Params        map[string]string // Query parameters
    JSON          interface{}       // JSON body (auto-marshaled)
    Data          interface{}       // Raw body data
    Auth          *Auth            // Request-specific auth
    Timeout       time.Duration    // Request timeout
    AllowRedirect bool            // Allow/disable redirects
}
```

### Authentication

#### Bearer Token
```go
auth := &easyhttp.Auth{
    Token: "your-bearer-token",
}
```

#### Basic Auth
```go
auth := &easyhttp.Auth{
    Username: "username",
    Password: "password",
}
```

### Response Methods

```go
resp, err := easyhttp.GET("https://httpbin.org/json")

// Check if successful (2xx status)
if resp.OK() {
    // Get response as string
    text := resp.Text()
    
    // Get response as bytes
    bytes := resp.Bytes()
    
    // Parse JSON into struct (traditional way)
    var data MyStruct
    err := resp.JSON(&data)
    
    // Or use easyjson's fluent API (much easier!)
    jsonValue, err := resp.JSONValue()
    if err == nil {
        // No struct needed - access JSON directly!
        name := jsonValue.Get("name").AsString()
        age := jsonValue.Get("age").AsInt()
        city := jsonValue.Q("address", "city").AsString() // Query chaining!
    }
}
```

## Examples

### GET with Query Parameters

```go
resp, err := easyhttp.GET("https://httpbin.org/get", &easyhttp.RequestOptions{
    Params: map[string]string{
        "page": "1",
        "size": "10",
    },
})
```

### POST with Form Data

```go
resp, err := easyhttp.POST("https://httpbin.org/post", &easyhttp.RequestOptions{
    Data: "key1=value1&key2=value2",
    Headers: map[string]string{
        "Content-Type": "application/x-www-form-urlencoded",
    },
})
```

### Custom Headers

```go
resp, err := easyhttp.GET("https://api.github.com/user", &easyhttp.RequestOptions{
    Headers: map[string]string{
        "Authorization": "token your-github-token",
        "Accept":        "application/vnd.github.v3+json",
    },
})
```

### Timeout Configuration

```go
// Per-request timeout
resp, err := easyhttp.GET("https://httpbin.org/delay/5", &easyhttp.RequestOptions{
    Timeout: 3 * time.Second,
})

// Default client timeout
client := easyhttp.New().SetTimeout(10 * time.Second)
```

### Error Handling

```go
resp, err := easyhttp.GET("https://httpbin.org/status/404")
if err != nil {
    log.Fatal("Request failed:", err)
}

if !resp.OK() {
    fmt.Printf("HTTP Error: %d\n", resp.StatusCode)
}
```

### Working with JSON APIs - The Easy Way

```go
// Instead of defining structs, use easyjson's fluent API
resp, err := easyhttp.GET("https://jsonplaceholder.typicode.com/users/1")
if err != nil {
    log.Fatal(err)
}

if resp.OK() {
    jsonValue, err := resp.JSONValue()
    if err != nil {
        log.Fatal("Failed to parse JSON:", err)
    }
    
    // Access JSON data without struct definitions!
    name := jsonValue.Get("name").AsString()
    email := jsonValue.Get("email").AsString()
    city := jsonValue.Q("address", "city").AsString()  // Query chaining for nested access
    
    fmt.Printf("User: %s (%s) from %s\n", name, email, city)
    
    // Check if fields exist
    if jsonValue.Has("company") {
        company := jsonValue.Q("company", "name").AsString()
        fmt.Printf("Works at: %s\n", company)
    }
}
```

### Creating Complex JSON with easyjson

```go
// Build complex JSON structures easily
user := easyjson.NewObject()
user.Set("name", "Alice Johnson")
user.Set("email", "alice@example.com")

// Add nested objects
address := easyjson.NewObject()
address.Set("street", "123 Main St")
address.Set("city", "Boston")
user.Set("address", address.Raw())

// Add arrays
hobbies := easyjson.NewArrayFrom([]interface{}{"reading", "hiking", "coding"})
user.Set("hobbies", hobbies.Raw())

// Send it
resp, err := easyhttp.POST("https://api.example.com/users", &easyhttp.RequestOptions{
    JSON: user,
})
```

### Building an API Client

```go
type APIClient struct {
    *easyhttp.Client
    baseURL string
}

func NewAPIClient(baseURL, token string) *APIClient {
    client := easyhttp.New().
        SetBaseURL(baseURL).
        SetAuth(&easyhttp.Auth{Token: token}).
        SetHeaders(map[string]string{
            "Content-Type": "application/json",
            "User-Agent":   "MyAPI-Client/1.0",
        })
    
    return &APIClient{
        Client:  client,
        baseURL: baseURL,
    }
}

func (c *APIClient) GetUser(id int) (*User, error) {
    resp, err := c.GET(fmt.Sprintf("/users/%d", id))
    if err != nil {
        return nil, err
    }
    
    if !resp.OK() {
        return nil, fmt.Errorf("API error: %d", resp.StatusCode)
    }
    
    var user User
    if err := resp.JSON(&user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (c *APIClient) CreateUser(user *User) (*User, error) {
    resp, err := c.POST("/users", &easyhttp.RequestOptions{
        JSON: user,
    })
    if err != nil {
        return nil, err
    }
    
    if !resp.OK() {
        return nil, fmt.Errorf("API error: %d", resp.StatusCode)
    }
    
    var createdUser User
    if err := resp.JSON(&createdUser); err != nil {
        return nil, err
    }
    
    return &createdUser, nil
}
```

## Performance

EasyHTTP is built on top of Go's native `net/http` package with easyjson for JSON operations, giving you excellent performance:

- **HTTP Performance**: Connection pooling, reuse, HTTP/2 support, efficient memory usage
- **JSON Flexibility**: Uses github.com/javanhut/easyjson which provides a fluent, Python-like API for JSON manipulation
- **Concurrent request handling**
- **Minimal overhead** while providing a convenient API

### easyjson Benefits

The library uses a specialized easyjson implementation that focuses on ease of use:

- **Fluent API**: Access JSON without defining structs
- **Query chaining**: `json.Q("user", "address", "city").AsString()`
- **Python-like**: Similar to Python's dict/list access patterns
- **Type conversion**: Automatic conversion with `.AsString()`, `.AsInt()`, `.AsBool()`, etc.
- **Path access**: `json.Path("user.address.city")` for dot notation

```go
// No struct definitions needed!
jsonValue, err := resp.JSONValue()
if err == nil {
    // Much easier than traditional Go JSON handling
    name := jsonValue.Get("name").AsString()
    age := jsonValue.Get("age").AsInt()
    city := jsonValue.Q("address", "city").AsString()
    
    // Check existence
    if jsonValue.Has("optional_field") {
        value := jsonValue.Get("optional_field").AsString()
    }
}
```

## Comparison with Standard Library

### Standard net/http
```go
// Standard library approach
client := &http.Client{}
req, err := http.NewRequest("POST", "https://api.example.com/users", bytes.NewReader(jsonData))
if err != nil {
    return err
}
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer "+token)

resp, err := client.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
    return err
}

var result User
if err := json.Unmarshal(body, &result); err != nil {
    return err
}
```

### EasyHTTP
```go
// EasyHTTP approach
resp, err := easyhttp.POST("https://api.example.com/users", &easyhttp.RequestOptions{
    JSON: userData,
    Auth: &easyhttp.Auth{Token: token},
})
if err != nil {
    return err
}

var result User
if err := resp.JSON(&result); err != nil {
    return err
}
```

## Dependencies

- **Go 1.18+**
- **github.com/javanhut/easyjson** - High-performance JSON library

The library has minimal dependencies to maintain performance and reduce bloat.

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linting (`make lint`)
6. Commit your changes (`git commit -am 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/javanhut/easyhttp.git
cd easyhttp

# Install development tools
make dev-setup

# Install dependencies
make deps

# Run tests
make test

# Run all checks (format, lint, test)
make all
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
