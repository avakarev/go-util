// Package httputil implements http helpers
package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// BaseClient implements base http client
type BaseClient struct {
	BaseURL string
	Header  http.Header
	Cookies []http.Cookie
	Timeout time.Duration
}

func (c *BaseClient) url(path string) string {
	if path == "" {
		return c.BaseURL
	}
	if strings.HasPrefix(path, "?") {
		return c.BaseURL + path
	}
	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		return path
	}
	return strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

// Req performs new http request with given method, path and optional json body
func (c *BaseClient) Req(method string, path string, payload any) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}

	r, err := http.NewRequest(method, c.url(path), body)
	if err != nil {
		return nil, err
	}

	for key, values := range c.Header {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}

	for i := range c.Cookies {
		r.AddCookie(&c.Cookies[i])
	}

	return (&http.Client{
		Timeout: c.Timeout,
	}).Do(r)
}

// Get makes get request and returns response body
func (c *BaseClient) Get(path string) ([]byte, int, error) {
	resp, err := c.Req(http.MethodGet, path, nil)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	bytes, err := io.ReadAll(resp.Body)
	return bytes, resp.StatusCode, err
}

func (c *BaseClient) decodeJSON(resp *http.Response, destPtr any) (int, error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	if destPtr == nil {
		return resp.StatusCode, nil
	}

	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(destPtr)
}

// GetJSON makes get request and decodes response body as json
func (c *BaseClient) GetJSON(path string, destPtr any) (int, error) {
	resp, err := c.Req(http.MethodGet, path, nil)
	if err != nil {
		return 0, err
	}
	return c.decodeJSON(resp, destPtr)
}

// PostJSON makes post request and decodes response body as json
func (c *BaseClient) PostJSON(path string, payload any, destPtr any) (int, error) {
	resp, err := c.Req(http.MethodPost, path, payload)
	if err != nil {
		return 0, err
	}
	return c.decodeJSON(resp, destPtr)
}

// PutJSON makes put request and decodes response body as json
func (c *BaseClient) PutJSON(path string, payload any, destPtr any) (int, error) {
	resp, err := c.Req(http.MethodPut, path, payload)
	if err != nil {
		return 0, err
	}
	return c.decodeJSON(resp, destPtr)
}

// DelJSON makes delete request and decodes response body as json
func (c *BaseClient) DelJSON(path string, payload any, destPtr any) (int, error) {
	resp, err := c.Req(http.MethodDelete, path, payload)
	if err != nil {
		return 0, err
	}
	return c.decodeJSON(resp, destPtr)
}
