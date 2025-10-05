# Contributing to Common Utils

We welcome contributions from the community! Whether you're fixing bugs, adding features, improving documentation, or helping with testing, your contributions make this project better for everyone.

## Quick Start

### Prerequisites
- **Go 1.19+** - [Install Go](https://golang.org/doc/install)
- **Git** - For version control

### Setting Up Development Environment

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/common-utils.git
cd common-utils

# 3. Add upstream remote
git remote add upstream https://github.com/mustanish/common-utils.git

# 4. Verify everything works
go test ./...
```

## How to Contribute

### Reporting Issues

**Bug Reports:**
- Use GitHub Issues with clear title and description
- Include Go version, OS, and steps to reproduce
- Provide minimal code example if possible
- Check existing issues first to avoid duplicates

**Feature Requests:**
- Describe the problem you're trying to solve
- Explain how the feature would benefit users
- Consider if it fits the library's scope (common utilities)

### Code Contributions

#### 1. Create a Branch
```bash
git checkout -b feature/amazing-feature
# or
git checkout -b fix/bug-description
```

#### 2. Make Your Changes
- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation as needed
- Keep changes focused and atomic

#### 3. Testing
```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...
```

#### 4. Commit Your Changes
```bash
# Use clear, descriptive commit messages
git add .
git commit -m "feat: add new utility method for date parsing"

# Or for bug fixes
git commit -m "fix: handle nil pointer in assertion utility"
```

**Commit Message Format:**
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `test:` - Adding or modifying tests
- `refactor:` - Code refactoring
- `perf:` - Performance improvements

#### 5. Submit Pull Request
```bash
# Push to your fork
git push origin feature/amazing-feature

# Create Pull Request on GitHub
```

## Project Structure

```
common-utils/
├── assertionutil/          # Safe type assertions
│   ├── client.go          # Main implementation
│   └── client_test.go     # Tests
├── collectionutil/         # Collection operations
│   ├── client.go
│   └── client_test.go
├── dateutil/              # Date/time utilities
│   ├── client.go
│   └── client_test.go
├── httputil/              # HTTP client utilities
│   ├── client.go
│   ├── client_test.go
│   ├── errors.go
│   ├── request.go
│   └── response.go
├── scripts/               # Automation and utility scripts
│   └── check-version.sh   # Version consistency checker
├── CHANGELOG.md           # Version history
├── CONTRIBUTING.md        # This file
├── go.mod                 # Go module definition
├── LICENSE                # MIT license
├── README.md              # Project documentation
└── VERSION                # Current version (used by automation)
```

### VERSION File
The `VERSION` file serves as the single source of truth for the current version:

**Purpose:**
- **Automation**: Scripts read this for tagging releases
- **Consistency**: Keeps README, CHANGELOG, and git tags in sync
- **Build Scripts**: Used by release automation tools

**Usage Examples:**
```bash
# Read current version
VERSION=$(cat VERSION)

# Tag git release
git tag "v$(cat VERSION)"

# Update documentation
sed -i "s/v[0-9.]\+/v$VERSION/g" README.md
```

**Updating VERSION:**
- Only update when preparing a release
- Follow semantic versioning (MAJOR.MINOR.PATCH)
- Update CHANGELOG.md when changing VERSION
- Use `scripts/check-version.sh` to verify version consistency

## Coding Standards

### Go Style Guide
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `go fmt` for formatting
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use meaningful variable and function names

### Interface Design Pattern
All utilities follow this consistent pattern:

```go
// 1. Define interface
type NewUtilClient interface {
    DoSomething(param string) (result string, err error)
}

// 2. Implement struct
type NewUtil struct {
    // fields if needed
}

// 3. Constructor function
func NewNewUtil() NewUtilClient {
    return &NewUtil{}
}

// 4. Implement methods
func (u *NewUtil) DoSomething(param string) (string, error) {
    // implementation
}
```

### Testing Guidelines
- **Test Coverage**: Aim for >90% statement coverage
- **Test Organization**: Use subtests for multiple scenarios
- **Test Naming**: Use descriptive test names
- **Benchmarks**: Add benchmarks for performance-critical code

Example test structure:
```go
func TestNewUtilMethod(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "test", "expected", false},
        {"invalid input", "", "", true},
    }
    
    util := NewNewUtil()
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := util.DoSomething(tt.input)
            
            if tt.wantErr && err == nil {
                t.Error("expected error but got none")
            }
            
            if !tt.wantErr && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

## Adding a New Utility Package

### 1. Planning
- Ensure it fits the project scope (common, reusable utilities)
- Check it doesn't overlap with existing packages
- Consider the 80/20 rule (focus on commonly used operations)

### 2. Implementation Steps

```bash
# Create new package directory
mkdir newutil
cd newutil

# Create interface and implementation
touch client.go client_test.go
```

**client.go structure:**
```go
package newutil

// NewUtilClient defines the interface for the new utility
type NewUtilClient interface {
    // Add your methods here
    ProcessData(input string) (string, error)
}

// NewUtil provides the implementation
type NewUtil struct {
    // Add fields if needed
}

// NewNewUtil creates a new utility instance
func NewNewUtil() NewUtilClient {
    return &NewUtil{}
}

// ProcessData implements the interface method
func (u *NewUtil) ProcessData(input string) (string, error) {
    // Your implementation here
}
```

### 3. Testing
```go
package newutil

import "testing"

func TestNewNewUtil(t *testing.T) {
    util := NewNewUtil()
    if util == nil {
        t.Error("NewNewUtil() returned nil")
    }
}

func TestProcessData(t *testing.T) {
    // Add comprehensive tests
}
```

### 4. Documentation
- Add package to main README.md
- Update examples section
- Add godoc comments for public functions

## Code Review Process

### What We Look For
- **Functionality**: Does it work as intended?
- **Tests**: Adequate test coverage and quality
- **Performance**: No unnecessary allocations or slow operations
- **API Design**: Clean, intuitive interface
- **Documentation**: Clear godoc comments
- **Compatibility**: Maintains backward compatibility

### Review Checklist
- [ ] Code follows Go conventions
- [ ] Tests pass and provide good coverage
- [ ] Documentation is clear and complete
- [ ] No breaking changes (unless major version)
- [ ] Performance is acceptable
- [ ] Error handling is appropriate

## Documentation Guidelines

### README.md
- Keep examples concise but complete
- Update package table when adding new utilities
- Maintain consistent formatting

### Code Documentation
```go
// ProcessData processes the input string and returns a formatted result.
// It returns an error if the input is empty or contains invalid characters.
func (u *NewUtil) ProcessData(input string) (string, error) {
    // implementation
}
```

## Release Process

### Version Management
- Version follows [Semantic Versioning](https://semver.org/)
- Update `VERSION` file for releases
- Update `CHANGELOG.md` with changes

### Release Types
- **MAJOR** (x.0.0): Breaking changes
- **MINOR** (x.y.0): New features, backward compatible
- **PATCH** (x.y.z): Bug fixes, backward compatible

### Manual Release Steps
1. Update `VERSION` file with new version
2. Update `CHANGELOG.md` with changes
3. Update version references in `README.md`
4. Run `scripts/check-version.sh` to verify consistency
5. Commit changes: `git commit -m "release: v2.1.0"`
6. Create git tag: `git tag v2.1.0`
7. Push changes and tag: `git push origin main --tags`
8. Create GitHub release from tag

## Getting Help

### Communication Channels
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and general discussion
- **Pull Request Reviews**: Code-specific discussions

### Questions?
- Check existing issues and discussions first
- Provide context and examples when asking questions
- Be respectful and patient

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions
- Special thanks for major features or improvements

## Development Workflow

### Daily Development
```bash
# Stay up to date
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/new-feature

# Make changes, test, commit
go test ./...
git add .
git commit -m "feat: add new feature"

# Push and create PR
git push origin feature/new-feature
```

### Before Submitting
- [ ] All tests pass
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated (if applicable)

---

Thank you for contributing to Common Utils! 🙏

Your contributions help make Go development easier for everyone.