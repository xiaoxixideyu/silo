package platform

import (
	"context"
	"fmt"
	config "silo/pkg/platform/config"
	"silo/pkg/platform/request"
	"silo/pkg/platform/response"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// PlatformClient platform client
type PlatformClient struct {
	client *resty.Client
	config config.PlatformConfig
}

// NewPlatformClient create new platform client
func NewPlatformClient(cfg config.PlatformConfig) (*PlatformClient, error) {
	// Validate required configuration
	if cfg.Addr == "" {
		return nil, fmt.Errorf("platform address cannot be empty")
	}

	// Create Resty client
	client := resty.New()

	// Set base URL
	client.SetBaseURL(cfg.Addr)

	// Set timeout
	if cfg.Timeout > 0 {
		client.SetTimeout(cfg.Timeout)
	} else {
		return nil, fmt.Errorf("timeout must be greater than 0")
	}

	// Set retry
	if cfg.RetryCount > 0 {
		client.SetRetryCount(cfg.RetryCount)
	} else {
		return nil, fmt.Errorf("retry count must be greater than 0")
	}

	if cfg.RetryWaitTime > 0 {
		client.SetRetryWaitTime(cfg.RetryWaitTime)
	} else {
		return nil, fmt.Errorf("retry wait time must be greater than 0")
	}

	// Set debug mode
	client.SetDebug(cfg.Debug)

	// Set default request headers
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "application/json")

	// Set authentication information
	if cfg.APIKey != "" {
		client.SetHeader("X-API-Key", cfg.APIKey)
	}
	if cfg.AuthToken != "" {
		client.SetAuthToken(cfg.AuthToken)
	}

	// Set request and response middleware
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.Debug().
			Str("method", req.Method).
			Str("url", req.URL).
			Msg("Sending HTTP request")
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log.Debug().
			Int("status_code", resp.StatusCode()).
			Str("url", resp.Request.URL).
			Dur("duration", resp.Time()).
			Msg("Received HTTP response")
		return nil
	})

	return &PlatformClient{
		client: client,
		config: cfg,
	}, nil
}

// doRequest generic method to execute HTTP requests
func (c *PlatformClient) doRequest(ctx context.Context, method, path string, body any, queryParams map[string]string) ([]byte, error) {
	req := c.client.R().SetContext(ctx)

	// Add query parameters
	for key, value := range queryParams {
		req.SetQueryParam(key, value)
	}

	// Set request body
	if body != nil {
		req.SetBody(body)
	}

	var resp *resty.Response
	var err error

	// Send request based on method
	switch method {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "DELETE":
		resp, err = req.Delete(path)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		log.Error().Err(err).Str("method", method).Str("path", path).Msg("HTTP request failed")
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}

	if !resp.IsSuccess() {
		log.Error().
			Int("status_code", resp.StatusCode()).
			Str("response", string(resp.Body())).
			Msg("HTTP request returned error status code")
		return nil, fmt.Errorf("%s request returned error status code: %d, response: %s", method, resp.StatusCode(), string(resp.Body()))
	}

	return resp.Body(), nil
}

// SetHeader set request header
func (c *PlatformClient) SetHeader(key, value string) {
	c.client.SetHeader(key, value)
}

// SetAuthToken set authentication token
func (c *PlatformClient) SetAuthToken(token string) {
	c.client.SetAuthToken(token)
}

// GetConfig get client configuration
func (c *PlatformClient) GetConfig() config.PlatformConfig {
	return c.config
}

// GetRestyClient get underlying Resty client (for advanced usage)
func (c *PlatformClient) GetRestyClient() *resty.Client {
	return c.client
}

// ===== Example related APIs =====

// GetExamples get example list
func (c *PlatformClient) GetExamples(ctx context.Context, req request.GetExamplesRequest) (*response.GetExamplesResponse, error) {
	queryParams := map[string]string{
		"page":  fmt.Sprintf("%d", req.Page),
		"limit": fmt.Sprintf("%d", req.Limit),
	}
	if req.Search != "" {
		queryParams["search"] = req.Search
	}

	respBody, err := c.doRequest(ctx, "GET", "/api/examples", nil, queryParams)
	if err != nil {
		return nil, err
	}

	var result response.GetExamplesResponse
	if err := response.ParseResponse(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetExample get single example information
func (c *PlatformClient) GetExample(ctx context.Context, exampleID int64) (*response.Example, error) {
	path := fmt.Sprintf("/api/examples/%d", exampleID)

	respBody, err := c.doRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result response.Example
	if err := response.ParseResponse(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateExample create example
func (c *PlatformClient) CreateExample(ctx context.Context, req request.CreateExampleRequest) (*response.CreateExampleResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", "/api/examples", req, nil)
	if err != nil {
		return nil, err
	}

	var result response.CreateExampleResponse
	if err := response.ParseResponse(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateExample update example
func (c *PlatformClient) UpdateExample(ctx context.Context, exampleID int64, req request.UpdateExampleRequest) (*response.UpdateExampleResponse, error) {
	path := fmt.Sprintf("/api/examples/%d", exampleID)

	respBody, err := c.doRequest(ctx, "PUT", path, req, nil)
	if err != nil {
		return nil, err
	}

	var result response.UpdateExampleResponse
	if err := response.ParseResponse(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteExample delete example
func (c *PlatformClient) DeleteExample(ctx context.Context, exampleID int64) (*response.DeleteExampleResponse, error) {
	path := fmt.Sprintf("/api/examples/%d", exampleID)

	respBody, err := c.doRequest(ctx, "DELETE", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var result response.DeleteExampleResponse
	if err := response.ParseResponse(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Reverse reverse string
func (c *PlatformClient) Reverse(ctx context.Context, req request.ReverseRequest) (*response.ReverseResponse, error) {
	respBody, err := c.doRequest(ctx, "POST", "/api/examples/reverse", req, nil)
	if err != nil {
		return nil, err
	}

	var result response.ReverseResponse
	if err := response.ParseResponse(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ===== Generic API methods =====

// Get generic GET request
func (c *PlatformClient) Get(ctx context.Context, path string, queryParams map[string]string) ([]byte, error) {
	return c.doRequest(ctx, "GET", path, nil, queryParams)
}

// Post generic POST request
func (c *PlatformClient) Post(ctx context.Context, path string, body any) ([]byte, error) {
	return c.doRequest(ctx, "POST", path, body, nil)
}

// Put generic PUT request
func (c *PlatformClient) Put(ctx context.Context, path string, body any) ([]byte, error) {
	return c.doRequest(ctx, "PUT", path, body, nil)
}

// Delete generic DELETE request
func (c *PlatformClient) Delete(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, "DELETE", path, nil, nil)
}

// PostForm send form POST request
func (c *PlatformClient) PostForm(ctx context.Context, path string, formData map[string]string) ([]byte, error) {
	req := c.client.R().
		SetContext(ctx).
		SetFormData(formData)

	resp, err := req.Post(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("POST form request failed")
		return nil, fmt.Errorf("POST form request failed: %w", err)
	}

	if !resp.IsSuccess() {
		log.Error().
			Int("status_code", resp.StatusCode()).
			Str("response", string(resp.Body())).
			Msg("POST form request returned error status code")
		return nil, fmt.Errorf("POST form request returned error status code: %d, response: %s", resp.StatusCode(), string(resp.Body()))
	}

	return resp.Body(), nil
}
