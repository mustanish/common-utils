package httputil

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

// HTTPClient defines the interface for the custom HTTP client
type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
	Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	SetRetryHook(hook func(attempt int, resp *http.Response, err error))
	SetSuccessHook(hook func(resp *http.Response, options RequestOptions))

	// Response helpers
	ReadBody(resp *http.Response) ([]byte, error)
	DecodeJSON(resp *http.Response, v any) error
	IsSuccess(resp *http.Response) bool
	GetHeader(resp *http.Response, key string) string
	CloseResponse(resp *http.Response)
}

// HTTPUtil is a custom HTTP client with retry logic and enhanced logging
type HTTPUtil struct {
	Client         *http.Client
	MaxRetries     int
	InitialWait    time.Duration
	MaxWait        time.Duration
	Logger         *logrus.Logger
	RequestTimeout time.Duration
	RetryOnStatus  []int

	RetryHook   func(attempt int, resp *http.Response, err error)
	SuccessHook func(resp *http.Response, options RequestOptions)
}

// NewHTTPUtil creates a new HTTP client with default settings
func NewHTTPUtil(logger *logrus.Logger) HTTPClient {
	client := &HTTPUtil{
		Client: &http.Client{
			Timeout: 10 * time.Minute,
			Transport: &http.Transport{
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   20,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   30 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableCompression:    false,
				ForceAttemptHTTP2:     true,
				ResponseHeaderTimeout: 60 * time.Second,
			},
		},
		MaxRetries:  5,
		InitialWait: 5 * time.Second,
		MaxWait:     60 * time.Second,
		Logger:      logger,
		RetryOnStatus: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
	}

	// Set default hooks
	client.setDefaultHooks()

	return client
}

// SetRetryHook allows customizing the retry hook
func (h *HTTPUtil) SetRetryHook(hook func(attempt int, resp *http.Response, err error)) {
	h.RetryHook = hook
}

// SetSuccessHook allows customizing the success hook
func (h *HTTPUtil) SetSuccessHook(hook func(resp *http.Response, options RequestOptions)) {
	h.SuccessHook = hook
}

// setDefaultHooks configures the default hook implementations
func (h *HTTPUtil) setDefaultHooks() {
	h.RetryHook = func(attempt int, resp *http.Response, err error) {
		fields := logrus.Fields{
			"attempt": attempt + 1,
			"max":     h.MaxRetries,
			"wait":    h.InitialWait,
		}
		if err != nil {
			fields["error"] = err.Error()
		}
		if resp != nil {
			fields["status"] = resp.StatusCode
		}
		h.Logger.WithFields(fields).Warn("Request failed, retrying")
	}

	h.SuccessHook = func(resp *http.Response, options RequestOptions) {
		h.Logger.WithFields(logrus.Fields{"method": options.Method, "url": options.URL, "status": resp.StatusCode}).Info("Request completed successfully")
	}
}

// shouldRetry determines if a request should be retried
func (h *HTTPUtil) shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true // Always retry on errors
	}
	return funk.Contains(h.RetryOnStatus, resp.StatusCode)
}
