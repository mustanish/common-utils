package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestIsSuccess(t *testing.T) {
	util := &HTTPUtil{}

	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{199, false}, // Below 200
		{200, true},  // Success
		{201, true},  // Created
		{299, true},  // Last 2xx
		{300, false}, // Redirect
		{404, false}, // Not found
		{500, false}, // Server error
	}

	for _, tc := range testCases {
		resp := &http.Response{StatusCode: tc.statusCode}
		result := util.IsSuccess(resp)
		if result != tc.expected {
			t.Errorf("For status %d, expected %v, got %v", tc.statusCode, tc.expected, result)
		}
	}
}

func TestReadBody(t *testing.T) {
	body := "hello world"
	resp := &http.Response{Body: io.NopCloser(bytes.NewBufferString(body))}
	util := &HTTPUtil{}
	data, err := util.ReadBody(resp)
	if err != nil || string(data) != body {
		t.Errorf("Expected body '%s', got '%s', err: %v", body, string(data), err)
	}
}

func TestDecodeJSON(t *testing.T) {
	type Foo struct {
		Bar string `json:"bar"`
	}
	obj := Foo{Bar: "baz"}
	buf, _ := json.Marshal(obj)
	resp := &http.Response{Body: io.NopCloser(bytes.NewBuffer(buf))}
	util := &HTTPUtil{}
	var out Foo
	if err := util.DecodeJSON(resp, &out); err != nil || out.Bar != "baz" {
		t.Errorf("Expected Bar to be 'baz', got '%s', err: %v", out.Bar, err)
	}
}

func TestGetHeader(t *testing.T) {
	resp := &http.Response{Header: http.Header{"X-Test": []string{"value"}}}
	util := &HTTPUtil{}
	if val := util.GetHeader(resp, "X-Test"); val != "value" {
		t.Errorf("Expected header 'X-Test' to be 'value', got '%s'", val)
	}
}

func TestCloseResponse(t *testing.T) {
	resp := &http.Response{Body: io.NopCloser(bytes.NewBufferString("test"))}
	util := &HTTPUtil{}
	util.CloseResponse(resp)
}

func TestRetryExhaustedError_Error(t *testing.T) {
	err := &RetryExhaustedError{
		LastError:  errors.New("fail"),
		LastStatus: 500,
		Attempts:   3,
		URL:        "http://test",
		Method:     "GET",
	}
	msg := err.Error()
	if msg == "" {
		t.Error("Expected error message to be non-empty")
	}
}

func TestHTTPUtil_Get_Success(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()
	resp, err := util.Get(context.Background(), server.URL, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	data, _ := util.ReadBody(resp)
	if string(data) != "ok" {
		t.Errorf("Expected body 'ok', got '%s'", string(data))
	}
}

func TestHTTPUtil_Post_Success(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("created"))
	}))
	defer server.Close()
	resp, err := util.Post(context.Background(), server.URL, bytes.NewBufferString("data"), nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	data, _ := util.ReadBody(resp)
	if string(data) != "created" {
		t.Errorf("Expected body 'created', got '%s'", string(data))
	}
}

func TestHTTPUtil_Put_Success(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("updated"))
	}))
	defer server.Close()
	resp, err := util.Put(context.Background(), server.URL, bytes.NewBufferString("data"), nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	data, _ := util.ReadBody(resp)
	if string(data) != "updated" {
		t.Errorf("Expected body 'updated', got '%s'", string(data))
	}
}

func TestHTTPUtil_Patch_Success(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}

		// Verify body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}
		expectedBody := "patch_data"
		if string(body) != expectedBody {
			t.Errorf("Expected body '%s', got '%s'", expectedBody, string(body))
		}

		w.WriteHeader(200)
		_, _ = w.Write([]byte("patched"))
	}))
	defer server.Close()

	resp, err := util.Patch(context.Background(), server.URL, bytes.NewBufferString("patch_data"), nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	data, _ := util.ReadBody(resp)
	if string(data) != "patched" {
		t.Errorf("Expected body 'patched', got '%s'", string(data))
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestHTTPUtil_Patch_WithHeaders(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Errorf("Expected X-Custom-Header 'test-value', got '%s'", r.Header.Get("X-Custom-Header"))
		}

		w.WriteHeader(200)
		_, _ = w.Write([]byte("success"))
	}))
	defer server.Close()

	headers := map[string]string{
		"Content-Type":    "application/json",
		"X-Custom-Header": "test-value",
	}

	resp, err := util.Patch(context.Background(), server.URL, bytes.NewBufferString(`{"field":"value"}`), headers)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestHTTPUtil_Patch_ErrorHandling(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		expectedStatus int
		shouldError    bool
	}{
		{
			name: "404 Not Found",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
				_, _ = w.Write([]byte("not found"))
			},
			expectedStatus: 404,
			shouldError:    false,
		},
		{
			name: "400 Bad Request",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("bad request"))
			},
			expectedStatus: 400,
			shouldError:    false,
		},
		{
			name: "422 Unprocessable Entity",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(422)
				_, _ = w.Write([]byte("validation error"))
			},
			expectedStatus: 422,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			resp, err := util.Patch(context.Background(), server.URL, bytes.NewBufferString("data"), nil)

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if resp != nil && resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestHTTPUtil_Patch_InvalidURL(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	// Use a malformed URL that will fail immediately
	_, err := util.Patch(context.Background(), "://invalid-url", nil, nil)
	if err == nil {
		t.Error("Expected error for invalid URL, got none")
	}
}

func TestHTTPUtil_Patch_ContextCancellation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow server
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := util.Patch(ctx, server.URL, bytes.NewBufferString("data"), nil)
	if err == nil {
		t.Error("Expected context timeout error, got none")
	}
}

func TestHTTPUtil_Patch_NilBody(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}
		if len(body) != 0 {
			t.Errorf("Expected empty body, got %s", string(body))
		}

		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	resp, err := util.Patch(context.Background(), server.URL, nil, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestHTTPUtil_Delete_Success(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer server.Close()
	resp, err := util.Delete(context.Background(), server.URL, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestSetRetryHookAndSuccessHook(t *testing.T) {
	logger := logrus.New()
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	successCalled := false
	util.SetSuccessHook(func(resp *http.Response, options RequestOptions) {
		successCalled = true
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()
	_, err := util.Get(context.Background(), server.URL, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !successCalled {
		t.Error("Expected success hook to be called")
	}
}

func TestShouldRetry(t *testing.T) {
	logger := logrus.New()
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	resp := &http.Response{StatusCode: http.StatusInternalServerError}
	if !util.shouldRetry(resp, nil) {
		t.Error("Expected shouldRetry to return true for retryable status")
	}
	resp.StatusCode = 200
	if util.shouldRetry(resp, nil) {
		t.Error("Expected shouldRetry to return false for non-retryable status")
	}
	if !util.shouldRetry(nil, errors.New("fail")) {
		t.Error("Expected shouldRetry to return true for error")
	}
}

func TestHTTPUtil_EmptyMethod(t *testing.T) {
	logger := logrus.New()
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	opts := RequestOptions{
		Method: "",
		URL:    "http://example.com",
	}
	_, err := util.doRequest(opts)
	if err == nil || !strings.Contains(err.Error(), "method cannot be empty") {
		t.Errorf("Expected method validation error, got: %v", err)
	}
}

func TestHTTPUtil_EmptyURL(t *testing.T) {
	logger := logrus.New()
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	opts := RequestOptions{
		Method: "GET",
		URL:    "",
	}
	_, err := util.doRequest(opts)
	if err == nil || !strings.Contains(err.Error(), "URL cannot be empty") {
		t.Errorf("Expected URL validation error, got: %v", err)
	}
}

func TestHTTPUtil_RetryLogic(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	util.MaxRetries = 2
	util.InitialWait = 10 * time.Millisecond
	util.MaxWait = 50 * time.Millisecond

	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	_, err := util.Get(context.Background(), server.URL, nil)
	if err != nil {
		t.Errorf("Expected success after retries, got error: %v", err)
	}
	if retryCount != 3 {
		t.Errorf("Expected 3 attempts (1 + 2 retries), got %d", retryCount)
	}
}

func TestHTTPUtil_RetryExhausted(t *testing.T) {
	logger := logrus.New()
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	util.MaxRetries = 1
	util.InitialWait = 1 * time.Millisecond

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := util.Get(context.Background(), server.URL, nil)
	if err == nil {
		t.Error("Expected error after retry exhaustion")
	}

	var retryErr *RetryExhaustedError
	if !errors.As(err, &retryErr) {
		t.Errorf("Expected RetryExhaustedError, got %T", err)
	}
}

func TestHTTPUtil_ContextCancellation(t *testing.T) {
	logger := logrus.New()
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := util.Get(ctx, "http://example.com", nil)
	if err == nil {
		t.Error("Expected context cancellation error")
	}
}

func TestHTTPUtil_RateLimitWithRetryAfter(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	util.MaxRetries = 1
	util.InitialWait = 1 * time.Millisecond

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	start := time.Now()
	_, err := util.Get(context.Background(), server.URL, nil)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Expected success after rate limit retry, got: %v", err)
	}
	if duration < 800*time.Millisecond {
		t.Errorf("Expected to wait for Retry-After header, but completed too quickly: %v", duration)
	}
}

func TestDecodeJSON_InvalidJSON(t *testing.T) {
	testCases := []struct {
		name string
		json string
	}{
		{"missing quote", `{"invalid": json}`},
		{"missing comma", `{"a": "b" "c": "d"}`},
		{"trailing comma", `{"valid": "json",}`},
		{"empty body", ``},
	}

	util := &HTTPUtil{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Body: io.NopCloser(bytes.NewBufferString(tc.json))}
			var result map[string]any
			err := util.DecodeJSON(resp, &result)
			if err == nil {
				t.Errorf("Expected JSON decode error for %s", tc.name)
			}
		})
	}
}

func TestCloseResponse_NilResponse(t *testing.T) {
	util := &HTTPUtil{}
	util.CloseResponse(nil)
}

func TestCloseResponse_NilBody(t *testing.T) {
	util := &HTTPUtil{}
	resp := &http.Response{Body: nil}
	util.CloseResponse(resp)
}

func TestGetHeader_CaseInsensitive(t *testing.T) {
	util := &HTTPUtil{}
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}

	// Go HTTP headers are case-insensitive
	if val := util.GetHeader(resp, "content-type"); val != "application/json" {
		t.Errorf("Expected case-insensitive header lookup to work")
	}
}

func TestHTTPUtil_WithCustomHeaders(t *testing.T) {
	logger := logrus.New()
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom") != "test-value" {
			t.Errorf("Expected custom header, got: %s", r.Header.Get("X-Custom"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	headers := map[string]string{
		"X-Custom": "test-value",
	}

	_, err := util.Get(context.Background(), server.URL, headers)
	if err != nil {
		t.Errorf("Expected success with custom headers, got: %v", err)
	}
}

func TestNewHTTPUtil_DefaultSettings(t *testing.T) {
	logger := logrus.New()
	client := NewHTTPUtil(logger, nil)
	util := client.(*HTTPUtil)

	if util.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries=5, got %d", util.MaxRetries)
	}
	if util.InitialWait != 5*time.Second {
		t.Errorf("Expected InitialWait=5s, got %v", util.InitialWait)
	}
	if util.MaxWait != 60*time.Second {
		t.Errorf("Expected MaxWait=60s, got %v", util.MaxWait)
	}
}

func TestHTTPUtil_AllMethods(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	testCases := []struct {
		method   string
		testFunc func(string) (*http.Response, error)
	}{
		{"GET", func(url string) (*http.Response, error) {
			return util.Get(context.Background(), url, nil)
		}},
		{"POST", func(url string) (*http.Response, error) {
			return util.Post(context.Background(), url, bytes.NewBufferString("data"), nil)
		}},
		{"PUT", func(url string) (*http.Response, error) {
			return util.Put(context.Background(), url, bytes.NewBufferString("data"), nil)
		}},
		{"DELETE", func(url string) (*http.Response, error) {
			return util.Delete(context.Background(), url, nil)
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tc.method {
					t.Errorf("Expected method %s, got %s", tc.method, r.Method)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			_, err := tc.testFunc(server.URL)
			if err != nil {
				t.Errorf("Expected no error for %s, got: %v", tc.method, err)
			}
		})
	}
}

func TestHTTPUtil_RequestBody(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)

	expectedBody := "test request body"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != expectedBody {
			t.Errorf("Expected body '%s', got '%s'", expectedBody, string(body))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	_, err := util.Post(context.Background(), server.URL, bytes.NewBufferString(expectedBody), nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestHTTPUtil_TimeoutHandling(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger, nil).(*HTTPUtil)
	util.Client.Timeout = 100 * time.Millisecond
	util.MaxRetries = 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	_, err := util.Get(context.Background(), server.URL, nil)
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestReadBody_EmptyBody(t *testing.T) {
	resp := &http.Response{Body: io.NopCloser(bytes.NewBufferString(""))}
	util := &HTTPUtil{}
	data, err := util.ReadBody(resp)
	if err != nil {
		t.Errorf("Expected no error for empty body, got: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty data, got: %s", string(data))
	}
}

func TestDefaultHTTPConfig(t *testing.T) {
	config := DefaultHTTPConfig()

	// Test default values
	if config.ClientTimeout != 10*time.Minute {
		t.Errorf("Expected ClientTimeout to be 10 minutes, got %v", config.ClientTimeout)
	}
	if config.DisableCompression != false {
		t.Errorf("Expected DisableCompression to be false, got %v", config.DisableCompression)
	}
	if config.ForceAttemptHTTP2 != true {
		t.Errorf("Expected ForceAttemptHTTP2 to be true, got %v", config.ForceAttemptHTTP2)
	}
	if config.MaxIdleConnsPerHost != 20 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 20, got %v", config.MaxIdleConnsPerHost)
	}
	if config.MaxIdleConns != 100 {
		t.Errorf("Expected MaxIdleConns to be 100, got %v", config.MaxIdleConns)
	}
	if config.IdleConnTimeout != 90*time.Second {
		t.Errorf("Expected IdleConnTimeout to be 90 seconds, got %v", config.IdleConnTimeout)
	}
	if config.TLSHandshakeTimeout != 30*time.Second {
		t.Errorf("Expected TLSHandshakeTimeout to be 30 seconds, got %v", config.TLSHandshakeTimeout)
	}
	if config.ExpectContinueTimeout != 1*time.Second {
		t.Errorf("Expected ExpectContinueTimeout to be 1 second, got %v", config.ExpectContinueTimeout)
	}
	if config.ResponseHeaderTimeout != 60*time.Second {
		t.Errorf("Expected ResponseHeaderTimeout to be 60 seconds, got %v", config.ResponseHeaderTimeout)
	}
	if config.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries to be 5, got %v", config.MaxRetries)
	}
	if config.InitialWait != 5*time.Second {
		t.Errorf("Expected InitialWait to be 5 seconds, got %v", config.InitialWait)
	}
	if config.MaxWait != 60*time.Second {
		t.Errorf("Expected MaxWait to be 60 seconds, got %v", config.MaxWait)
	}

	// Test retry status codes
	expectedStatuses := []int{
		http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}
	if len(config.RetryOnStatus) != len(expectedStatuses) {
		t.Errorf("Expected %d retry status codes, got %d", len(expectedStatuses), len(config.RetryOnStatus))
	}
	for i, expected := range expectedStatuses {
		if config.RetryOnStatus[i] != expected {
			t.Errorf("Expected retry status code %d at index %d, got %d", expected, i, config.RetryOnStatus[i])
		}
	}
}

func TestNewHTTPUtil_NilConfig(t *testing.T) {
	logger := logrus.New()
	client := NewHTTPUtil(logger, nil)

	if client == nil {
		t.Fatal("Expected client to be created with nil config")
	}

	// Should use default configuration
	httpUtil := client.(*HTTPUtil)
	if httpUtil.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries to be 5 (default), got %v", httpUtil.MaxRetries)
	}
	if httpUtil.InitialWait != 5*time.Second {
		t.Errorf("Expected InitialWait to be 5 seconds (default), got %v", httpUtil.InitialWait)
	}
}

func TestNewHTTPUtil_CustomConfig(t *testing.T) {
	logger := logrus.New()

	// Create custom configuration
	config := &HTTPConfig{
		ClientTimeout:         5 * time.Minute,
		DisableCompression:    true,
		ForceAttemptHTTP2:     false,
		MaxIdleConnsPerHost:   50,
		MaxIdleConns:          200,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		MaxRetries:            3,
		InitialWait:           2 * time.Second,
		MaxWait:               30 * time.Second,
		RetryOnStatus: []int{
			http.StatusInternalServerError,
			http.StatusBadGateway,
		},
	}

	client := NewHTTPUtil(logger, config)
	httpUtil := client.(*HTTPUtil)

	// Test custom values were applied
	if httpUtil.Client.Timeout != config.ClientTimeout {
		t.Errorf("Expected ClientTimeout to be %v, got %v", config.ClientTimeout, httpUtil.Client.Timeout)
	}
	if httpUtil.MaxRetries != config.MaxRetries {
		t.Errorf("Expected MaxRetries to be %v, got %v", config.MaxRetries, httpUtil.MaxRetries)
	}
	if httpUtil.InitialWait != config.InitialWait {
		t.Errorf("Expected InitialWait to be %v, got %v", config.InitialWait, httpUtil.InitialWait)
	}
	if httpUtil.MaxWait != config.MaxWait {
		t.Errorf("Expected MaxWait to be %v, got %v", config.MaxWait, httpUtil.MaxWait)
	}

	// Test transport settings
	transport := httpUtil.Client.Transport.(*http.Transport)
	if transport.DisableCompression != config.DisableCompression {
		t.Errorf("Expected DisableCompression to be %v, got %v", config.DisableCompression, transport.DisableCompression)
	}
	if transport.ForceAttemptHTTP2 != config.ForceAttemptHTTP2 {
		t.Errorf("Expected ForceAttemptHTTP2 to be %v, got %v", config.ForceAttemptHTTP2, transport.ForceAttemptHTTP2)
	}
	if transport.MaxIdleConnsPerHost != config.MaxIdleConnsPerHost {
		t.Errorf("Expected MaxIdleConnsPerHost to be %v, got %v", config.MaxIdleConnsPerHost, transport.MaxIdleConnsPerHost)
	}
	if transport.TLSHandshakeTimeout != config.TLSHandshakeTimeout {
		t.Errorf("Expected TLSHandshakeTimeout to be %v, got %v", config.TLSHandshakeTimeout, transport.TLSHandshakeTimeout)
	}

	// Test retry status codes
	if len(httpUtil.RetryOnStatus) != len(config.RetryOnStatus) {
		t.Errorf("Expected %d retry status codes, got %d", len(config.RetryOnStatus), len(httpUtil.RetryOnStatus))
	}
	for i, expected := range config.RetryOnStatus {
		if httpUtil.RetryOnStatus[i] != expected {
			t.Errorf("Expected retry status code %d at index %d, got %d", expected, i, httpUtil.RetryOnStatus[i])
		}
	}
}

func TestNewHTTPUtil_DefaultConfiguration(t *testing.T) {
	logger := logrus.New()

	// Test with nil config uses defaults
	client := NewHTTPUtil(logger, nil)
	httpUtil := client.(*HTTPUtil)

	// Should have default values
	if httpUtil.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries to be 5 (default), got %v", httpUtil.MaxRetries)
	}
	if httpUtil.Client.Timeout != 10*time.Minute {
		t.Errorf("Expected ClientTimeout to be 10 minutes (default), got %v", httpUtil.Client.Timeout)
	}

	transport := httpUtil.Client.Transport.(*http.Transport)
	if transport.MaxIdleConnsPerHost != 20 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 20 (default), got %v", transport.MaxIdleConnsPerHost)
	}
}

func TestNewHTTPUtil_PartialOverride(t *testing.T) {
	logger := logrus.New()

	// Test partial override - only set some properties
	client := NewHTTPUtil(logger, &HTTPConfig{
		ClientTimeout:       3 * time.Minute,
		MaxRetries:          2,
		MaxIdleConnsPerHost: 30,
		// Other properties not set - should use defaults
	})

	httpUtil := client.(*HTTPUtil)

	// Check overridden values
	if httpUtil.Client.Timeout != 3*time.Minute {
		t.Errorf("Expected ClientTimeout to be 3 minutes, got %v", httpUtil.Client.Timeout)
	}
	if httpUtil.MaxRetries != 2 {
		t.Errorf("Expected MaxRetries to be 2, got %v", httpUtil.MaxRetries)
	}

	transport := httpUtil.Client.Transport.(*http.Transport)
	if transport.MaxIdleConnsPerHost != 30 {
		t.Errorf("Expected MaxIdleConnsPerHost to be 30, got %v", transport.MaxIdleConnsPerHost)
	}

	// Check default values for non-overridden fields
	if transport.MaxIdleConns != 100 {
		t.Errorf("Expected MaxIdleConns to remain default (100), got %v", transport.MaxIdleConns)
	}
	if httpUtil.InitialWait != 5*time.Second {
		t.Errorf("Expected InitialWait to remain default (5s), got %v", httpUtil.InitialWait)
	}
	if transport.TLSHandshakeTimeout != 30*time.Second {
		t.Errorf("Expected TLSHandshakeTimeout to remain default (30s), got %v", transport.TLSHandshakeTimeout)
	}
}
