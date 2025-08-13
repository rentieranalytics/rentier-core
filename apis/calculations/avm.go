package calculations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (c *Client) CalculateAVM(
	ctx context.Context,
	body any,
	ropts ...RequestOption,
) (*AVMResponse, error) {
	var out AVMResponse
	_, err := c.doJSON(ctx, http.MethodPost, DefaultCalculateAVMPath, body, &out, ropts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) doJSON(
	ctx context.Context,
	method, path string,
	body any,
	out any,
	ropts ...RequestOption,
) ([]byte, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := *c.baseURL
	u.Path = strings.TrimRight(c.baseURL.Path, "/") + path

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, fmt.Errorf("encode request: %w", err)
		}
	}

	attempts := max(c.retryMax, 1)

	var lastErr error
	for i := range attempts {
		req, err := http.NewRequestWithContext(ctx, method, u.String(), io.Reader(buf))
		if err != nil {
			return nil, fmt.Errorf("new request: %w", err)
		}
		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		// default headers
		for k, vs := range c.defaultHeaders {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
		// auth
		if c.authHeader != "" && c.authValue != "" {
			req.Header.Set(c.authHeader, c.authValue)
		}
		// per request options
		for _, f := range ropts {
			f(req)
		}

		res, err := c.hc.Do(req)
		if err != nil {
			// network/timeout—retry
			lastErr = err
			if !shouldRetryHTTP(nil, i, attempts) {
				break
			}
			backoff(i, c.retryWaitMin, c.retryWaitMax)
			continue
		}
		defer res.Body.Close()
		b, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return nil, fmt.Errorf("read response: %w", readErr)
		}

		if res.StatusCode >= 200 && res.StatusCode < 300 {
			if out != nil && len(b) > 0 {
				if err := json.Unmarshal(b, out); err != nil {
					return nil, fmt.Errorf("decode response: %w", err)
				}
			}
			return b, nil
		}

		// Attempt to parse message for better errors
		var emsg struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(b, &emsg)
		ae := &APIError{
			StatusCode: res.StatusCode,
			Body:       b,
			Message:    strings.TrimSpace(emsg.Message),
		}

		// 429/5xx (except 501) -> retry
		if shouldRetryHTTP(res, i, attempts) {
			lastErr = ae
			backoff(i, c.retryWaitMin, c.retryWaitMax)
			continue
		}
		return nil, ae
	}
	return nil, fmt.Errorf("request failed after %d attempts: %w", attempts, lastErr)
}

func shouldRetryHTTP(res *http.Response, attempt, max int) bool {
	if attempt >= max-1 {
		return false
	}
	if res == nil {
		return true // network error
	}
	if res.StatusCode == http.StatusTooManyRequests {
		return true
	}
	if res.StatusCode >= 500 && res.StatusCode != http.StatusNotImplemented {
		return true
	}
	return false
}

func backoff(attempt int, min, max time.Duration) {
	if min <= 0 {
		min = 250 * time.Millisecond
	}
	if max <= 0 {
		max = 2 * time.Second
	}
	d := min * (1 << attempt)
	if d > max {
		d = max
	}
	time.Sleep(d)
}
