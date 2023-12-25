package client

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	client http.Client
	Logger *zap.Logger // TODO: move to pkg and make obfuscation
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	query := req.URL.Query()
	logParams := make([]zap.Field, 0, len(query))
	for k, v := range query {
		logParams = append(logParams, zap.String(k, strings.Join(v, ", ")))
	}
	logHeaders := make([]zap.Field, 0, len(req.Header))
	for k, v := range req.Header {
		logHeaders = append(logHeaders, zap.String(k, strings.Join(v, ", ")))
	}

	reqBody, _ := req.GetBody()
	logBody, _ := io.ReadAll(reqBody)
	_ = reqBody.Close()

	c.Logger.Debug(
		"starting request",
		zap.String("url", req.URL.String()),
		zap.String("host", req.URL.Host),
		zap.String("path", req.URL.Path),
		zap.Dict("params", logParams...), // FIXME: obfuscate
		zap.String("method", req.Method),
		zap.Dict("headers", logHeaders...),      // FIXME: obfuscate
		zap.ByteString("request body", logBody), // FIXME: obfuscate
	)

	resp, err := c.client.Do(req)
	if err != nil {
		c.Logger.Error(
			"got error in request",
			zap.String("url", req.URL.String()),
			zap.String("host", req.URL.Host),
			zap.String("path", req.URL.Path),
			zap.Dict("params", logParams...), // FIXME: obfuscate
			zap.String("method", req.Method),
			zap.Dict("headers", logHeaders...), // FIXME: obfuscate
			zap.ByteString("body", logBody),    // FIXME: obfuscate
			zap.Error(err),
		)
		return nil, err
	}

	logRespBody, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewBuffer(logRespBody))

	c.Logger.Debug(
		"got response",
		zap.String("url", req.URL.String()),
		zap.String("host", req.URL.Host),
		zap.String("path", req.URL.Path),
		zap.Dict("params", logParams...), // FIXME: obfuscate
		zap.String("method", req.Method),
		zap.Dict("headers", logHeaders...),      // FIXME: obfuscate
		zap.ByteString("request body", logBody), // FIXME: obfuscate
		zap.Int("response status code", resp.StatusCode),
		zap.String("response Content-Type", resp.Header.Get("Content-Type")),
		zap.ByteString("response body", logRespBody), // FIXME: obfuscate
	)

	return resp, nil
}

func NewClient(options ...ClientOption) *Client {
	c := new(Client)
	c.client = http.Client{
		Transport: http.DefaultTransport,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

type ClientOption func(*Client)

func WithClientTimeout(dur time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = dur
	}
}
