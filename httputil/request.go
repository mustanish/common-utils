package httputil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// RequestOptions holds options for the HTTP request
// It includes method, URL, body, headers, context, and timeout.
// This struct is used to encapsulate all parameters needed for making an HTTP request.
type RequestOptions struct {
	Method  string
	URL     string
	Body    io.Reader
	Headers map[string]string
	Context context.Context
}

// Get sends an HTTP GET request
func (h *HTTPUtil) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return h.doRequest(RequestOptions{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
		Context: ctx,
	})
}

// Post sends an HTTP POST request
func (h *HTTPUtil) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return h.doRequest(RequestOptions{
		Method:  http.MethodPost,
		URL:     url,
		Body:    body,
		Headers: headers,
		Context: ctx,
	})
}

// Put sends an HTTP PUT request
func (h *HTTPUtil) Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return h.doRequest(RequestOptions{
		Method:  http.MethodPut,
		URL:     url,
		Body:    body,
		Headers: headers,
		Context: ctx,
	})
}

// Delete sends an HTTP DELETE request
func (h *HTTPUtil) Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return h.doRequest(RequestOptions{
		Method:  http.MethodDelete,
		URL:     url,
		Headers: headers,
		Context: ctx,
	})
}

// doRequest performs an HTTP request with retry logic
func (h *HTTPUtil) doRequest(opts RequestOptions) (*http.Response, error) {
	var err error
	var bodyBytes []byte
	var resp *http.Response

	if opts.Method == "" {
		return nil, fmt.Errorf("method cannot be empty")
	}
	if opts.URL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	if opts.Body != nil {
		bodyBytes, err = io.ReadAll(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	// Log request start
	h.Logger.WithFields(logrus.Fields{"method": opts.Method, "url": opts.URL, "max_retries": h.MaxRetries}).Debug("Starting HTTP request")
	currentWait := h.InitialWait

	// Execute request with retries
	for attempt := 0; attempt <= h.MaxRetries; attempt++ {

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(opts.Context, opts.Method, opts.URL, bodyReader)
		if err != nil {
			h.Logger.WithFields(logrus.Fields{"error": err, "method": opts.Method, "url": opts.URL, "attempt": attempt + 1}).Error("Failed to create request")
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}

		resp, err = h.Client.Do(req)
		if err == nil && !h.shouldRetry(resp, err) {
			h.SuccessHook(resp, opts)
			return resp, nil
		}

		if attempt >= h.MaxRetries {
			break
		}

		h.RetryHook(attempt, resp, err)

		// Calculate wait time with exponential backoff and jitter
		jitter := time.Duration(rand.Float64() * float64(currentWait) * 0.1)
		waitTime := currentWait + jitter

		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			rateLimitWait := 60 * time.Second

			h.Logger.WithFields(logrus.Fields{"status": resp.StatusCode, "url": opts.URL}).Warn("Received 429 Too Many Requests")
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, parseErr := strconv.Atoi(retryAfter); parseErr == nil {
					rateLimitWait = time.Duration(seconds) * time.Second
				}
			}

			h.Logger.WithFields(logrus.Fields{"wait_time": rateLimitWait}).Info("Respecting Retry-After header wait time")
			if rateLimitWait > waitTime {
				waitTime = rateLimitWait
			}
		}

		h.Logger.WithFields(logrus.Fields{"wait_time": waitTime}).Info("Waiting before next retry")
		select {
		case <-opts.Context.Done():
			h.Logger.WithError(opts.Context.Err()).Warn("Request cancelled during retry wait")
			return nil, fmt.Errorf("context cancelled during retry: %w", opts.Context.Err())
		case <-time.After(waitTime):
			currentWait = time.Duration(math.Min(float64(currentWait)*1.5, float64(h.MaxWait)))
		}
	}

	h.Logger.WithFields(logrus.Fields{
		"method":  opts.Method,
		"url":     opts.URL,
		"retries": h.MaxRetries,
		"error":   err,
		"status": func() int {
			if resp != nil {
				return resp.StatusCode
			}
			return 0
		}(),
	}).Error("Request failed after all retries")

	return resp, &RetryExhaustedError{
		URL:      opts.URL,
		Method:   opts.Method,
		Attempts: h.MaxRetries + 1,
		LastStatus: func() int {
			if resp != nil {
				return resp.StatusCode
			}
			return 0
		}(),
		LastError: func() error {
			if err != nil {
				return err
			}
			return fmt.Errorf("unknown error after %d attempts", h.MaxRetries+1)
		}(),
	}
}
