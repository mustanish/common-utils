# Common Utils

[![Go Version](https://img.shields.io/github/go-mod/go-version/mustanish/common-utils)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/mustanish/common-utils)](https://goreportcard.com/report/github.com/mustanish/common-utils)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/mustanish/common-utils.svg)](https://pkg.go.dev/github.com/mustanish/common-utils)
[![Latest Release](https://img.shields.io/github/v/release/mustanish/common-utils)](https://github.com/mustanish/common-utils/releases)

A collection of reusable utility packages for Go applications. This library provides common functionality that can be shared across multiple services and projects to reduce code duplication and improve consistency.

## Installation

### Latest Stable Version

```bash
go get github.com/mustanish/common-utils/v2@v2.1.0
```

### Latest Development Version

```bash
go get github.com/mustanish/common-utils/v2@latest
```

### Specific Version

```bash
# Install a specific version
go get github.com/mustanish/common-utils/v2@v2.1.0
```

## Quick Start

```go
import (
    "github.com/mustanish/common-utils/v2/httputil"
    "github.com/mustanish/common-utils/v2/assertionutil"
    "github.com/mustanish/common-utils/v2/collectionutil"
    "github.com/mustanish/common-utils/v2/dateutil"
)

// HTTP client with retry logic
httpClient := httputil.NewHTTPUtil(logger, nil)
resp, err := httpClient.Get(ctx, "https://api.example.com", nil)

// Safe type assertions
assertUtil := assertionutil.NewAssertionUtil()
name := assertUtil.GetStringOrEmpty(data, "name")
tags, ok := assertUtil.GetStringSlice(data, "tags")

// Collection operations
collectionUtil := collectionutil.NewCollectionUtil()
unique := collectionUtil.SliceUnique([]string{"a", "b", "a"})
peopleMap, _ := collectionUtil.ConvertToMap(people, keyFunc)

// Date operations
dateUtil := dateutil.NewDateUtil()
date, _ := dateUtil.Parse("2023-10-05")
tomorrow := dateUtil.AddDays(dateUtil.Today(), 1)
```

## Packages

| Package | Purpose | Key Methods |
|---------|---------|-------------|
| **httputil** | HTTP client with retry logic | `Get`, `Post`, `DecodeJSON` |
| **assertionutil** | Safe type extraction | `GetStringOrEmpty`, `GetStringSlice`, `GetInt` |
| **collectionutil** | Collection operations | `SliceUnique`, `ConvertToMap`, `MapFilter` |
| **dateutil** | Date/time utilities | `Parse`, `AddDays`, `IsAfter`, `NowUTC` |

## Features

### HttpUtil
- Automatic retry with exponential backoff
- Rate limiting and context support
- JSON request/response helpers

### AssertionUtil
- Safe type extraction from `map[string]any`
- No panic, no error handling needed for common cases
- `GetStringOrEmpty`, `GetStringSlice`, `GetInt`, etc.

### CollectionUtil  
- Type conversions (`ConvertToInteger`, `ConvertToBool`)
- Slice operations (`SliceUnique`, `SliceFilter`, `SliceContains`)
- Map operations (`MapFilter`, `ConvertToMap`)

### DateUtil
- Flexible parsing with auto-format detection
- Date arithmetic (`AddDays`, `AddMonths`, `AddYears`)
- Business day calculations
- 5 essential date formats (RFC3339, SimpleDateTime, USDate, etc.)

## Examples

<details>
<summary>HTTP Client</summary>

```go
client := httputil.NewHTTPUtil(logger, nil)
resp, err := client.Get(ctx, "https://api.example.com", headers)
if client.IsSuccess(resp) {
    var result map[string]any
    client.DecodeJSON(resp, &result)
}
```
</details>

<details>
<summary>Safe Type Assertions</summary>

```go
util := assertionutil.NewAssertionUtil()

// No error handling needed - returns empty string if missing/wrong type
username := util.GetStringOrEmpty(data, "username")
email := util.GetStringOrEmpty(data, "email")

// Safe slice extraction
if tags, ok := util.GetStringSlice(data, "tags"); ok {
    // Process string slice
}
```
</details>

<details>
<summary>Collection Operations</summary>

```go
util := collectionutil.NewCollectionUtil()

// Remove duplicates
unique := util.SliceUnique([]string{"a", "b", "a", "c"}) // ["a", "b", "c"]

// Convert slice to map
peopleMap, _ := util.ConvertToMap(people, func(item any) string {
    return item.(Person).Name
})

// Type conversion
age, _ := util.ConvertToInteger("25")
active, _ := util.ConvertToBool("yes") // true
```
</details>

<details>
<summary>Date Operations</summary>

```go
util := dateutil.NewDateUtil()

// Parse various formats automatically
date, _ := util.Parse("2023-10-05")
date, _ := util.Parse("10/05/2023")

// Date arithmetic
future := util.AddDays(util.Today(), 30)
if util.IsAfter(future, util.Today()) {
    // 30 days from now
}

// Business days
if util.IsBusinessDay(util.Today()) {
    nextBizDay := util.AddBusinessDays(util.Today(), 5)
}
```
</details>

## Requirements

- Go 1.19+
- See [pkg.go.dev](https://pkg.go.dev/github.com/mustanish/common-utils/v2) for full API documentation

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

**Quick start for contributors:**
```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/common-utils.git
cd common-utils

# Run tests
go test ./...

# Format code  
go fmt ./...
```

**Development workflow:**
- Fork the repository
- Create a feature branch (`git checkout -b feature/amazing-feature`)
- Make your changes with tests
- Submit a pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines, coding standards, and development setup.

See [CHANGELOG.md](CHANGELOG.md) for version history.

## License

MIT License - see [LICENSE](LICENSE) file.

---

**Made with ❤️ by [Mustanish](https://github.com/mustanish) and community contributors**
