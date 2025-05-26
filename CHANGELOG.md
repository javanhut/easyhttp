# Changelog

All notable changes to the EasyHTTP project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-26-05

### Added
- Initial release of EasyHTTP - Go HTTP client library
- Python requests-like API for Go
- Integration with easyjson for fluent JSON handling
- Support for all HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD)
- Client configuration with method chaining
- Authentication support (Bearer tokens and Basic auth)
- Flexible request options (headers, query parameters, timeouts, redirects)
- Multiple response formats (Text, JSON, JSONValue, Bytes)
- Comprehensive test suite with 95%+ coverage
- Benchmark tests for performance measurement
- Complete documentation and examples

### Features
- **Simple API**: One-liner HTTP requests like `easyhttp.GET("https://api.com")`
- **JSON Integration**: Fluent JSON access with `resp.JSONValue().Get("field").AsString()`
- **Client Builder**: Method chaining for configuration
- **Authentication**: Built-in support for tokens and basic auth
- **Performance**: Built on Go's native HTTP client for maximum performance
- **Testing**: Comprehensive test suite with mock server
- **Documentation**: Complete examples and API documentation

### Technical Details
- Go 1.21+ compatible
- Depends on github.com/javanhut/easyjson for JSON handling
- Thread-safe for concurrent usage
- Memory efficient with connection pooling
- HTTP/2 support through Go's native client
