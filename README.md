# Golang Common Utilities

[![Go Version](https://img.shields.io/github/go-mod/go-version/mustanish/common-utils)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/mustanish/common-utils)](https://goreportcard.com/report/github.com/mustanish/common-utils)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/mustanish/common-utils.svg)](https://pkg.go.dev/github.com/mustanish/common-utils)

A collection of reusable utility packages for Go applications. This library provides common functionality that can be shared across multiple services and projects to reduce code duplication and improve consistency.

## Requirements

- **Go 1.22+** - This module requires Go version 1.22 or higher

## Table of Contents

- [Installation](#installation)
- [Available Utilities](#available-utilities)
  - [HTTP Utility](#http-utility)
  - [Assertion Utility](#assertion-utility)
- [Development](#development)
- [Contributing](#contributing)

## Installation

```bash
go get github.com/mustanish/common-utils
```

## Available Utilities

### HTTP Utility

The HTTP utility provides a robust HTTP client with retry logic, rate limiting, and comprehensive error handling.

#### Features

- **Automatic Retry Logic**: Configurable retry mechanism with exponential backoff
- **Rate Limiting**: Respects `Retry-After` headers for proper rate limiting
- **Context Support**: Full context cancellation and timeout support
- **Logging**: Structured logging with configurable levels
- **Hooks**: Customizable retry and success hooks for monitoring
- **Response Helpers**: Convenient methods for response handling

#### Quick Start

```go
package main

import (
    "context"
    "log"
    "bytes"
    
    "github.com/sirupsen/logrus"
    "github.com/mustanish/common-utils/httputil"
)

func main() {
    // Create logger
    logger := logrus.New()
    logger.SetLevel(logrus.InfoLevel)
    
    // Create HTTP client
    client := httputil.NewHTTPUtil(logger)
    
    // Make a GET request
    resp, err := client.Get(context.Background(), "https://api.example.com/data", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer client.CloseResponse(resp)
    
    // Check if response is successful
    if client.IsSuccess(resp) {
        // Read response body
        body, err := client.ReadBody(resp)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("Response: %s", string(body))
    }
}
```

#### Configuration

```go
// Create with custom configuration
client := httputil.NewHTTPUtil(logger)
httpClient := client.(*httputil.HTTPUtil)

// Configure retry settings
httpClient.MaxRetries = 3
httpClient.InitialWait = 2 * time.Second
httpClient.MaxWait = 30 * time.Second

// Set custom HTTP client timeout
httpClient.Client.Timeout = 10 * time.Second
```

#### Available Methods

##### HTTP Methods
- `Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error)`
- `Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)`
- `Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)`
- `Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error)`

##### Response Helpers
- `IsSuccess(resp *http.Response) bool` - Checks if status code is 2xx
- `ReadBody(resp *http.Response) ([]byte, error)` - Reads response body
- `DecodeJSON(resp *http.Response, v interface{}) error` - Decodes JSON response
- `GetHeader(resp *http.Response, key string) string` - Gets response header
- `CloseResponse(resp *http.Response)` - Safely closes response body

##### Hooks
- `SetRetryHook(hook func(attempt int, resp *http.Response, err error))` - Called on each retry
- `SetSuccessHook(hook func(resp *http.Response, options RequestOptions))` - Called on success

#### Advanced Usage

##### Custom Headers
```go
headers := map[string]string{
    "Authorization": "Bearer " + token,
    "Content-Type":  "application/json",
    "X-API-Version": "v1",
}

resp, err := client.Get(ctx, url, headers)
```

##### JSON Requests
```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{Name: "John", Email: "john@example.com"}
jsonData, _ := json.Marshal(user)

headers := map[string]string{
    "Content-Type": "application/json",
}

resp, err := client.Post(ctx, url, bytes.NewBuffer(jsonData), headers)
```

##### Response Processing
```go
resp, err := client.Get(ctx, url, nil)
if err != nil {
    return err
}
defer client.CloseResponse(resp)

if !client.IsSuccess(resp) {
    return fmt.Errorf("request failed with status: %d", resp.StatusCode)
}

var result map[string]interface{}
if err := client.DecodeJSON(resp, &result); err != nil {
    return fmt.Errorf("failed to decode JSON: %w", err)
}
```

##### Monitoring with Hooks
```go
// Set retry hook for monitoring
client.SetRetryHook(func(attempt int, resp *http.Response, err error) {
    if resp != nil {
        logger.Warnf("Retry attempt %d for %s: status %d", 
            attempt, resp.Request.URL, resp.StatusCode)
    } else {
        logger.Warnf("Retry attempt %d: %v", attempt, err)
    }
})

// Set success hook for metrics
client.SetSuccessHook(func(resp *http.Response, options httputil.RequestOptions) {
    duration := time.Since(options.StartTime)
    logger.Infof("Request to %s completed in %v with status %d", 
        options.URL, duration, resp.StatusCode)
})
```

#### Error Handling

The HTTP utility provides detailed error information:

```go
resp, err := client.Get(ctx, url, headers)
if err != nil {
    // Check if it's a retry exhausted error
    var retryErr *httputil.RetryExhaustedError
    if errors.As(err, &retryErr) {
        log.Printf("Request failed after %d attempts. Last status: %d", 
            retryErr.Attempts, retryErr.LastStatus)
    }
    return err
}
```

#### Testing

The package includes comprehensive test coverage. Run tests with:

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Or using go directly
go test ./...
```

### Assertion Utility

The Assertion utility provides safe type assertion and extraction utilities for working with `map[string]any` data structures commonly found in JSON parsing, configuration handling, and dynamic data processing.

#### Features

- **Safe Type Assertions**: Existence checks with proper type validation
- **Required Field Validation**: Descriptive error messages for missing fields
- **Multiple Type Support**: string, int, int64, float64, bool, map, slice
- **Default Value Fallbacks**: Convenient methods with fallback values
- **Numeric Type Conversion**: Cross-type numeric conversions with loss prevention
- **Nested Path Navigation**: Access deeply nested data structures safely
- **Bulk Validation**: Validate multiple required fields in one call
- **Zero-Allocation Design**: High-performance with minimal allocations
- **JSON Compatibility**: Handles JSON unmarshaling type quirks (float64 → int)

#### Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/mustanish/common-utils/assertionutil"
)

func main() {
    // Create assertion utility
    util := assertionutil.NewAssertionUtil()
    
    // Sample data (e.g., from JSON unmarshaling)
    data := map[string]any{
        "name":    "John Doe",
        "age":     30.0,  // JSON numbers are float64
        "active":  true,
        "profile": map[string]interface{}{
            "email": "john@example.com",
            "phone": "123-456-7890",
        },
        "tags": []interface{}{"admin", "user"},
    }
    
    // Basic extraction
    name, ok := util.GetString(data, "name")
    if ok {
        fmt.Printf("Name: %s\n", name)
    }
    
    // Get with default values
    age := util.GetIntWithDefault(data, "age", 0)
    status := util.GetStringWithDefault(data, "status", "inactive")
    
    fmt.Printf("Age: %d, Status: %s\n", age, status)
    
    // Required field validation
    email, err := util.GetStringRequired(data, "email")
    if err != nil {
        log.Printf("Error: %v", err)
    }
    
    // Nested access
    email, ok = util.GetNestedString(data, "profile", "email")
    if ok {
        fmt.Printf("Email: %s\n", email)
    }
    
    // Bulk validation
    err = util.ValidateRequired(data, "name", "age")
    if err != nil {
        log.Printf("Validation failed: %v", err)
    } else {
        fmt.Println("All required fields present")
    }
}
```

#### Available Methods

##### Basic Type Getters
- `GetString(m map[string]any, key string) (string, bool)` - Extract non-empty strings
- `GetInt(m map[string]any, key string) (int, bool)` - Extract integers (handles JSON float64)
- `GetInt64(m map[string]any, key string) (int64, bool)` - Extract int64 values
- `GetFloat64(m map[string]any, key string) (float64, bool)` - Extract float64 values
- `GetBool(m map[string]any, key string) (bool, bool)` - Extract boolean values
- `GetMap(m map[string]any, key string) (map[string]interface{}, bool)` - Extract nested maps
- `GetSlice(m map[string]any, key string) ([]interface{}, bool)` - Extract non-empty slices

##### Required Field Validation
- `GetStringRequired(m map[string]any, key string) (string, error)` - Required string with error

##### Default Value Methods
- `GetStringWithDefault(m map[string]any, key, defaultValue string) string`
- `GetIntWithDefault(m map[string]any, key string, defaultValue int) int`
- `GetFloat64WithDefault(m map[string]any, key string, defaultValue float64) float64`
- `GetBoolWithDefault(m map[string]any, key string, defaultValue bool) bool`

##### Numeric Conversion Utilities
- `GetNumericAsFloat64(m map[string]any, key string) (float64, bool)` - Convert any numeric type to float64
- `GetNumericAsInt(m map[string]any, key string) (int, bool)` - Convert numeric types to int (loss-safe)

##### Nested Path Navigation
- `GetNestedString(m map[string]any, path ...string) (string, bool)` - Navigate nested paths
- `GetNestedMap(m map[string]any, path ...string) (map[string]interface{}, bool)` - Access nested maps

##### Validation Utilities
- `HasKey(m map[string]any, key string) bool` - Check key existence
- `HasNonEmptyString(m map[string]any, key string) bool` - Check for valid string
- `ValidateRequired(m map[string]any, keys ...string) error` - Bulk validation

#### Advanced Usage

##### JSON Processing
```go
// Typical JSON unmarshaling scenario
var data map[string]any
json.Unmarshal(jsonBytes, &data)

util := assertionutil.NewAssertionUtil()

// Handle JSON number types safely
userID := util.GetIntWithDefault(data, "user_id", 0)  // Works with 123.0 from JSON
score := util.GetNumericAsFloat64(data, "score")      // Handles int, float64, etc.
```

##### Configuration Handling
```go
// Load configuration from various sources
config := map[string]any{
    "server": map[string]interface{}{
        "host": "localhost",
        "port": 8080.0,  // From JSON
        "ssl":  true,
    },
    "database": map[string]interface{}{
        "url":     "postgres://...",
        "timeout": 30.0,
    },
}

util := assertionutil.NewAssertionUtil()

// Extract nested configuration safely
host := util.GetStringWithDefault(config, "server.host", "0.0.0.0")
port := util.GetNestedString(config, "server", "port")
ssl, _ := util.GetNestedString(config, "server", "ssl")

// Validate required configuration
err := util.ValidateRequired(config, "server", "database")
if err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}
```

##### Dynamic Data Processing
```go
// Process dynamic data with type safety
processRecord := func(record map[string]any) error {
    util := assertionutil.NewAssertionUtil()
    
    // Validate required fields first
    if err := util.ValidateRequired(record, "id", "type", "timestamp"); err != nil {
        return fmt.Errorf("invalid record: %w", err)
    }
    
    // Extract with defaults
    id := util.GetStringWithDefault(record, "id", "")
    recordType := util.GetStringWithDefault(record, "type", "unknown")
    priority := util.GetIntWithDefault(record, "priority", 1)
    enabled := util.GetBoolWithDefault(record, "enabled", true)
    
    // Process nested metadata
    if metadata, ok := util.GetNestedMap(record, "metadata"); ok {
        // Further process metadata...
    }
    
    return nil
}
```

##### Bulk Operations
```go
// Validate multiple configurations
configs := []map[string]any{
    {"name": "service1", "port": 8080.0},
    {"name": "service2", "port": 8081.0},
    {"name": "service3"}, // Missing port
}

util := assertionutil.NewAssertionUtil()

for i, config := range configs {
    if err := util.ValidateRequired(config, "name", "port"); err != nil {
        log.Printf("Config %d invalid: %v", i, err)
        continue
    }
    
    name, _ := util.GetString(config, "name")
    port := util.GetIntWithDefault(config, "port", 8080)
    
    fmt.Printf("Service: %s, Port: %d\n", name, port)
}
```

#### Performance

The assertion utility is designed for high performance:

- **Basic operations**: ~8.6 ns/op with 0 allocations
- **Nested access**: ~29 ns/op with 0 allocations  
- **Bulk validation**: ~29.5 ns/op with 0 allocations
- **Test coverage**: 93.5% statement coverage

#### Testing

```bash
# Run assertion utility tests
go test ./assertionutil -v

# Run with coverage
go test ./assertionutil -cover

# Run benchmarks
go test ./assertionutil -bench=. -benchmem
```

## Development

### 🔧 Quick Start

```bash
# Clone and setup
git clone https://github.com/mustanish/common-utils.git
cd common-utils

# Run tests
go test ./...

# Build
go build ./...
```

### 📊 Test Coverage

Both packages maintain high test coverage:
- **httputil**: >93% coverage
- **assertionutil**: >93% coverage

```bash
# Generate coverage report
go test -cover ./...

# Detailed coverage with HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Future Utilities

This library is designed to grow with additional utilities. Each utility follows a consistent pattern:

```go
import "github.com/mustanish/common-utils/{utilname}"

// Create utility instance
util := utilname.New{UtilName}()

// Use utility methods
result, err := util.SomeMethod(params)
```

## Contributing

We welcome contributions from the community! Whether you're fixing bugs, adding features, or improving documentation, your contributions are valuable.

### 🚀 Getting Started

#### Prerequisites
- **Go 1.22+** - Ensure you have Go 1.22 or higher installed
- **Git** - For version control

#### Development Setup

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/common-utils.git
cd common-utils

# 3. Verify everything works
go test ./...
```

### 📝 How to Contribute

#### Reporting Issues
- **Bug Reports**: Use GitHub Issues with a clear description and reproduction steps
- **Feature Requests**: Describe the utility or feature you'd like to see

#### Submitting Changes

1. **Create a Branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

2. **Make Your Changes**
   - Follow Go conventions
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Your Changes**
   ```bash
   go test ./...
   go vet ./...
   ```

4. **Submit a Pull Request**
   - Create a Pull Request on GitHub with a clear description

### 🏗️ Adding a New Utility

#### Directory Structure
```
newutil/
├── client.go      # Main implementation
└── client_test.go # Tests
```

#### Implementation Pattern
```go
// NewUtilClient defines the interface
type NewUtilClient interface {
    DoSomething(param string) (result string, err error)
}

// NewUtil provides the implementation
type NewUtil struct {}

// NewNewUtil creates a new utility instance
func NewNewUtil() NewUtilClient {
    return &NewUtil{}
}
```

### 📋 Code Standards

- **Formatting**: Use `go fmt`
- **Testing**: Aim for >90% test coverage
- **Documentation**: Add godoc comments for public functions
- **Error Handling**: Return descriptive errors

## Support

- **Issues**: Report bugs and request features via GitHub issues
- **Discussions**: Join community discussions for questions and ideas
- **Documentation**: Check this README and godoc for detailed API documentation

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Made with ❤️ by [Mustanish](https://github.com/mustanish) and community contributors**