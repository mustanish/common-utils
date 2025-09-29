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
	util := NewHTTPUtil(logger).(*HTTPUtil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
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
	util := NewHTTPUtil(logger).(*HTTPUtil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("created"))
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
	util := NewHTTPUtil(logger).(*HTTPUtil)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("updated"))
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

func TestHTTPUtil_Delete_Success(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
	util := NewHTTPUtil(logger).(*HTTPUtil)

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
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
			var result map[string]interface{}
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
	util := NewHTTPUtil(logger).(*HTTPUtil)

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
	client := NewHTTPUtil(logger)
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
	util := NewHTTPUtil(logger).(*HTTPUtil)

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
	util := NewHTTPUtil(logger).(*HTTPUtil)

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
	util := NewHTTPUtil(logger).(*HTTPUtil)
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
