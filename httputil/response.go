package httputil

import (
	"encoding/json"
	"io"
	"net/http"
)

// ReadBody reads and returns the response body as bytes.
func (h *HTTPUtil) ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// DecodeJSON decodes the response body into the provided struct.
func (h *HTTPUtil) DecodeJSON(resp *http.Response, v any) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

// IsSuccess returns true if the response status code is 2xx.
func (h *HTTPUtil) IsSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// GetHeader returns the value of a specific header from the response.
func (h *HTTPUtil) GetHeader(resp *http.Response, key string) string {
	return resp.Header.Get(key)
}

// CloseResponse safely closes the response body and ignores any error.
func (h *HTTPUtil) CloseResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}
