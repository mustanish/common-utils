package httputil

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

// HTTPConfig holds configuration for HTTP transport settings
type HTTPConfig struct {
	// Client timeouts
	ClientTimeout time.Duration

	// Transport settings
	DisableCompression    bool
	ForceAttemptHTTP2     bool
	MaxIdleConnsPerHost   int
	MaxIdleConns          int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	ResponseHeaderTimeout time.Duration

	// Retry settings
	MaxRetries    int
	InitialWait   time.Duration
	MaxWait       time.Duration
	RetryOnStatus []int
}

// DefaultHTTPConfig returns default configuration
func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		ClientTimeout:         10 * time.Minute,
		DisableCompression:    false,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   20,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		MaxRetries:            5,
		InitialWait:           5 * time.Second,
		MaxWait:               60 * time.Second,
		RetryOnStatus: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
	}
}

// HTTPClient defines the interface for the custom HTTP client
type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
	Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
	Patch(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error)
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

// NewHTTPUtil creates a new HTTP client with configuration
// Pass nil for config to use all defaults, or pass config with only the properties you want to override
func NewHTTPUtil(logger *logrus.Logger, config *HTTPConfig) HTTPClient {
	defaults := DefaultHTTPConfig()

	if config != nil {
		if config.ClientTimeout != 0 {
			defaults.ClientTimeout = config.ClientTimeout
		}
		if config.MaxIdleConnsPerHost != 0 {
			defaults.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost
		}
		if config.MaxIdleConns != 0 {
			defaults.MaxIdleConns = config.MaxIdleConns
		}
		if config.IdleConnTimeout != 0 {
			defaults.IdleConnTimeout = config.IdleConnTimeout
		}
		if config.TLSHandshakeTimeout != 0 {
			defaults.TLSHandshakeTimeout = config.TLSHandshakeTimeout
		}
		if config.ExpectContinueTimeout != 0 {
			defaults.ExpectContinueTimeout = config.ExpectContinueTimeout
		}
		if config.ResponseHeaderTimeout != 0 {
			defaults.ResponseHeaderTimeout = config.ResponseHeaderTimeout
		}
		if config.MaxRetries != 0 {
			defaults.MaxRetries = config.MaxRetries
		}
		if config.InitialWait != 0 {
			defaults.InitialWait = config.InitialWait
		}
		if config.MaxWait != 0 {
			defaults.MaxWait = config.MaxWait
		}
		if config.RetryOnStatus != nil {
			defaults.RetryOnStatus = config.RetryOnStatus
		}

		defaults.DisableCompression = config.DisableCompression
		defaults.ForceAttemptHTTP2 = config.ForceAttemptHTTP2
	}

	client := &HTTPUtil{
		Client: &http.Client{
			Timeout: defaults.ClientTimeout,
			Transport: &http.Transport{
				DisableCompression:    defaults.DisableCompression,
				ForceAttemptHTTP2:     defaults.ForceAttemptHTTP2,
				MaxIdleConnsPerHost:   defaults.MaxIdleConnsPerHost,
				MaxIdleConns:          defaults.MaxIdleConns,
				IdleConnTimeout:       defaults.IdleConnTimeout,
				TLSHandshakeTimeout:   defaults.TLSHandshakeTimeout,
				ExpectContinueTimeout: defaults.ExpectContinueTimeout,
				ResponseHeaderTimeout: defaults.ResponseHeaderTimeout,
			},
		},
		MaxRetries:    defaults.MaxRetries,
		InitialWait:   defaults.InitialWait,
		MaxWait:       defaults.MaxWait,
		Logger:        logger,
		RetryOnStatus: defaults.RetryOnStatus,
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
