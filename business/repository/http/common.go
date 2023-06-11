package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	JSON = "application/json" //default content type
)

type HttpRequest struct {
	URL     string        `json:"url,omitempty"`
	Headers http.Header   `json:"headers,omitempty"`
	Body    []byte        `json:"body,omitempty"`
	Timeout time.Duration `json:"timeout,omitempty" default:"5 * time.Minute"`
	Method  string        `json:"method,omitempty"`
}

// initiate a http call
func (request *HttpRequest) InitiateHttpCall(ctx context.Context, respObject interface{}) (int, error) {
	headers := addDefaultHeaders(ctx, request.Headers)
	reqBody := bytes.NewBuffer(request.Body)
	// req, err := http.NewRequest(request.Method, request.URL, reqBody)
	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, reqBody)
	if err != nil {
		// logging.Logger.WriteLogs(ctx, "error_creating_http_request_with_context", logging.ErrorLevel, logging.Fields{"error": err})
		return 0, err
	}
	req.Header = headers
	client := &http.Client{Timeout: request.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	err = json.Unmarshal(respBody, &respObject)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, nil
}

// Add  Request ID and session ID in headers
func addDefaultHeaders(c context.Context, headers http.Header) http.Header {
	if headers == nil {
		headers = http.Header{}
		headers.Set("Content-Type", JSON)
	}
	return headers
}
