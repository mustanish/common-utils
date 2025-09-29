package httputil

import "fmt"

// RetryExhaustedError is returned when all retry attempts are exhausted
// It contains details about the last error, status code, number of attempts, URL, and method.
// This error type is useful for understanding why a request ultimately failed after retries.
type RetryExhaustedError struct {
	LastError  error
	LastStatus int
	Attempts   int
	URL        string
	Method     string
}

// Error implements the error interface for RetryExhaustedError
func (e *RetryExhaustedError) Error() string {
	return fmt.Sprintf("retry exhausted after %d attempts for %s %s: HTTP %d: %v", e.Attempts, e.Method, e.URL, e.LastStatus, e.LastError)
}
