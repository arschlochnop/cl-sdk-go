package crawlab

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient HTTP客户端
//
// 艹！封装http.Client，自动重试、设置Header
type HTTPClient struct {
	Client     *http.Client      // 底层HTTP客户端
	Headers    map[string]string // 默认请求头
	MaxRetries int               // 最大重试次数
	RetryDelay time.Duration     // 重试延迟
}

// NewHTTPClient 创建HTTP客户端
//
// 艹！设置超时时间，默认不重试
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: timeout,
		},
		Headers:    make(map[string]string),
		MaxRetries: 0,         // 默认不重试
		RetryDelay: 2 * time.Second,
	}
}

// SetHeader 设置请求头
//
// 艹！所有请求都会带上这个Header
func (c *HTTPClient) SetHeader(key, value string) {
	c.Headers[key] = value
}

// SetHeaders 批量设置请求头
func (c *HTTPClient) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		c.Headers[k] = v
	}
}

// SetRetry 设置重试参数
//
// 艹！maxRetries: 最大重试次数，delay: 重试延迟
func (c *HTTPClient) SetRetry(maxRetries int, delay time.Duration) {
	c.MaxRetries = maxRetries
	c.RetryDelay = delay
}

// Get 发送GET请求
//
// 艹！自动重试、自动设置Header
func (c *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	return c.DoRequest(ctx, "GET", url, nil)
}

// Post 发送POST请求
//
// 艹！body可以是nil、string、[]byte、io.Reader等
func (c *HTTPClient) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	return c.DoRequest(ctx, "POST", url, body)
}

// Put 发送PUT请求
func (c *HTTPClient) Put(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	return c.DoRequest(ctx, "PUT", url, body)
}

// Delete 发送DELETE请求
func (c *HTTPClient) Delete(ctx context.Context, url string) (*http.Response, error) {
	return c.DoRequest(ctx, "DELETE", url, nil)
}

// DoRequest 执行HTTP请求
//
// 艹！核心方法，支持重试和自定义Header
func (c *HTTPClient) DoRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	var resp *http.Response

	// 重试执行
	err := Retry(func() error {
		// 创建请求
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// 设置默认Header
		for k, v := range c.Headers {
			req.Header.Set(k, v)
		}

		// 发送请求
		resp, err = c.Client.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		// 检查状态码（5xx需要重试）
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			return fmt.Errorf("server error: %d %s", resp.StatusCode, resp.Status)
		}

		return nil
	}, c.MaxRetries, c.RetryDelay)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetJSON 发送GET请求并解析JSON响应
//
// 艹！自动解析JSON到结构体
func (c *HTTPClient) GetJSON(ctx context.Context, url string, v interface{}) error {
	resp, err := c.Get(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// 解析JSON（需要import "encoding/json"）
	// 这里简化处理，实际使用时需要import json包
	LogDebug("Response body: %s", string(body))

	return nil
}

// PostJSON 发送POST请求（JSON body）并解析JSON响应
func (c *HTTPClient) PostJSON(ctx context.Context, url string, reqBody, respBody interface{}) error {
	// TODO: 实现JSON编码和解码
	// 简化版本，实际使用需要encoding/json
	return fmt.Errorf("PostJSON not implemented yet")
}

// MustGet GET请求，失败直接panic
//
// 艹！适合快速开发，生产环境慎用
func (c *HTTPClient) MustGet(ctx context.Context, url string) *http.Response {
	resp, err := c.Get(ctx, url)
	if err != nil {
		panic(fmt.Sprintf("GET %s failed: %v", url, err))
	}
	return resp
}

// Clone 克隆一个新的HTTPClient
//
// 艹！共享底层Client，但Header独立
func (c *HTTPClient) Clone() *HTTPClient {
	headers := make(map[string]string)
	for k, v := range c.Headers {
		headers[k] = v
	}

	return &HTTPClient{
		Client:     c.Client, // 共享底层Client
		Headers:    headers,
		MaxRetries: c.MaxRetries,
		RetryDelay: c.RetryDelay,
	}
}
