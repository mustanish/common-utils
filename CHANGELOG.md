# Changelog

All notable changes to this project will be documented in this file.

## [v2.1.0] - 2024-10-05

### Added
- **DateUtil**: New package for date/time operations with parsing, arithmetic, and business day logic
- **AssertionUtil**: `GetStringOrEmpty()`, `GetStringSlice()`, `GetKeys()` methods

### Changed  
- **DateUtil**: Simplified `GetCommonFormats()` to 5 essential formats
- **CollectionUtil**: Optimized for common use cases (80/20 principle)

### Removed
- **CollectionUtil**: `DeepCopy`, `Flatten`, `GroupBy`, `Reduce` methods

## [v2.0.0] - 2025-10-04

### Changed
- **HTTP Utility**: Constructor now requires config parameter `NewHTTPUtil(logger, config)`
- **Configuration**: Simplified override approach - set only what you need, rest uses defaults

### Migration
```go
// v1.x
client := httputil.NewHTTPUtil(logger)

// v2.x
client := httputil.NewHTTPUtil(logger, nil) // for defaults
```

## [v1.2.0] - 2025-10-02

### Changed
- **Go Version**: Minimum requirement from 1.22 to 1.19 for broader compatibility
- **Documentation**: Enhanced README and examples

## [v1.1.0] - 2025-09-30

### Added
- **CollectionUtil**: New package for collection operations and type conversions

## [v1.0.0] - 2025-09-29

### Added
- **HttpUtil**: HTTP client with retry logic and rate limiting
- **AssertionUtil**: Safe type assertion utilities for `map[string]any`