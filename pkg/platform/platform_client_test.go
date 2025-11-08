package platform

import (
	config2 "silo/pkg/platform/config"
	"testing"
	"time"

	"silo/pkg/platform/request"
)

func TestNewPlatformClient(t *testing.T) {
	// Test with valid configuration
	config := config2.PlatformConfig{
		Addr:          "http://localhost:8081",
		Timeout:       30 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 1 * time.Second,
		Debug:         false,
	}

	client, err := NewPlatformClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if client == nil {
		t.Fatal("Failed to create client")
	}

	clientConfig := client.GetConfig()
	if clientConfig.Addr != "http://localhost:8081" {
		t.Errorf("Expected address to be http://localhost:8081, got %s", clientConfig.Addr)
	}
}

func TestNewPlatformClientWithCustomConfig(t *testing.T) {
	customConfig := config2.PlatformConfig{
		Addr:          "https://api.example.com",
		Timeout:       10 * time.Second,
		RetryCount:    5,
		RetryWaitTime: 2 * time.Second,
		Debug:         true,
	}

	client, err := NewPlatformClient(customConfig)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if client == nil {
		t.Fatal("Failed to create client")
	}

	config := client.GetConfig()
	if config.Addr != "https://api.example.com" {
		t.Errorf("Expected address to be https://api.example.com, got %s", config.Addr)
	}

	if config.Timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10s, got %v", config.Timeout)
	}

	if config.RetryCount != 5 {
		t.Errorf("Expected retry count to be 5, got %d", config.RetryCount)
	}
}

func TestNewPlatformClientWithInvalidConfig(t *testing.T) {
	// Test empty address
	invalidConfig := config2.PlatformConfig{
		Addr: "",
	}
	_, err := NewPlatformClient(invalidConfig)
	if err == nil {
		t.Error("Expected empty address config to return error")
	}

	// Test zero timeout
	invalidConfig = config2.PlatformConfig{
		Addr:    "http://localhost:8081",
		Timeout: 0,
	}
	_, err = NewPlatformClient(invalidConfig)
	if err == nil {
		t.Error("Expected zero timeout config to return error")
	}

	// Test zero retry count
	invalidConfig = config2.PlatformConfig{
		Addr:       "http://localhost:8081",
		Timeout:    30 * time.Second,
		RetryCount: 0,
	}
	_, err = NewPlatformClient(invalidConfig)
	if err == nil {
		t.Error("Expected zero retry count config to return error")
	}

	// Test zero retry wait time
	invalidConfig = config2.PlatformConfig{
		Addr:          "http://localhost:8081",
		Timeout:       30 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 0,
	}
	_, err = NewPlatformClient(invalidConfig)
	if err == nil {
		t.Error("Expected zero retry wait time config to return error")
	}
}

func TestRequestStructures(t *testing.T) {
	// Test request structures
	createReq := request.CreateExampleRequest{
		Info: "Test Example",
		Val:  100,
	}

	if createReq.Info != "Test Example" {
		t.Errorf("Expected example info to be 'Test Example', got '%s'", createReq.Info)
	}

	if createReq.Val != 100 {
		t.Errorf("Expected example val to be 100, got %d", createReq.Val)
	}

	updateReq := request.UpdateExampleRequest{
		Info: "Updated example info",
		Val:  200,
	}

	if updateReq.Info != "Updated example info" {
		t.Errorf("Expected updated example info to be 'Updated example info', got '%s'", updateReq.Info)
	}

	if updateReq.Val != 200 {
		t.Errorf("Expected updated example val to be 200, got %d", updateReq.Val)
	}

	getExamplesReq := request.GetExamplesRequest{
		Page:  1,
		Limit: 10,
	}

	if getExamplesReq.Page != 1 {
		t.Errorf("Expected page to be 1, got %d", getExamplesReq.Page)
	}

	reverseReq := request.ReverseRequest{
		Info: "hello world",
	}

	if reverseReq.Info != "hello world" {
		t.Errorf("Expected reverse info to be 'hello world', got '%s'", reverseReq.Info)
	}
}

func TestClientMethods(t *testing.T) {
	config := config2.PlatformConfig{
		Addr:          "http://localhost:8081",
		Timeout:       30 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 1 * time.Second,
		Debug:         false,
	}
	client, err := NewPlatformClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test if methods exist and are callable
	// Since these are methods not function values, we just need to ensure they can be called without panic

	// Test generic methods
	_ = client.Get
	_ = client.Post
	_ = client.Put
	_ = client.Delete
	_ = client.PostForm

	// Test example related methods
	_ = client.GetExamples
	_ = client.GetExample
	_ = client.CreateExample
	_ = client.UpdateExample
	_ = client.DeleteExample
	_ = client.Reverse

	// Test helper methods
	_ = client.SetHeader
	_ = client.SetAuthToken
	_ = client.GetConfig
	_ = client.GetRestyClient
}

func TestSetHeader(t *testing.T) {
	config := config2.PlatformConfig{
		Addr:          "http://localhost:8081",
		Timeout:       30 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 1 * time.Second,
		Debug:         false,
	}
	client, err := NewPlatformClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test setting request header
	client.SetHeader("X-Custom-Header", "custom-value")

	// Since we cannot directly access internal client headers, here we only test that the method doesn't panic
	// Actual header setting tests require integration tests
}

func TestSetAuthToken(t *testing.T) {
	config := config2.PlatformConfig{
		Addr:          "http://localhost:8081",
		Timeout:       30 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 1 * time.Second,
		Debug:         false,
	}
	client, err := NewPlatformClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test setting authentication token
	client.SetAuthToken("test-token")

	// Since we cannot directly access internal client token, here we only test that the method doesn't panic
	// Actual token setting tests require integration tests
}
