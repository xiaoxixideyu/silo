package platform

import (
	"context"
	"fmt"
	config2 "silo/pkg/platform/config"
	"time"

	"silo/pkg/platform/request"
	"silo/pkg/platform/response"
)

// ExampleUsage demonstrates how to use PlatformClient
func ExampleUsage() {
	// Create configuration
	config := config2.PlatformConfig{
		Addr:          "http://localhost:8081",
		Timeout:       30 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 1 * time.Second,
		Debug:         true,
	}

	// Initialize client
	client, err := NewPlatformClient(config)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	// Usage example
	ctx := context.Background()

	// Get example list
	getExamplesReq := request.GetExamplesRequest{
		Page:  1,
		Limit: 10,
	}

	examplesResp, err := client.GetExamples(ctx, getExamplesReq)
	if err != nil {
		fmt.Printf("Failed to get example list: %v\n", err)
		return
	}
	fmt.Printf("Got %d examples:\n", examplesResp.Total)
	for _, example := range examplesResp.Examples {
		fmt.Printf("- ID: %d, Info: %s, Val: %d\n", example.ID, example.Info, example.Val)
	}

	// Create example
	createExampleReq := request.CreateExampleRequest{
		Info: "Test Example",
		Val:  100,
	}

	createResp, err := client.CreateExample(ctx, createExampleReq)
	if err != nil {
		fmt.Printf("Failed to create example: %v\n", err)
		return
	}
	fmt.Printf("Example created successfully, ID: %d\n", createResp.ID)

	// Get single example
	exampleID := createResp.ID
	example, err := client.GetExample(ctx, exampleID)
	if err != nil {
		fmt.Printf("Failed to get example: %v\n", err)
		return
	}
	fmt.Printf("Example info: ID=%d, Info=%s, Val=%d\n", example.ID, example.Info, example.Val)

	// Update example
	updateExampleReq := request.UpdateExampleRequest{
		Info: "Updated example info",
		Val:  200,
	}

	updateResp, err := client.UpdateExample(ctx, exampleID, updateExampleReq)
	if err != nil {
		fmt.Printf("Failed to update example: %v\n", err)
		return
	}
	fmt.Printf("Example updated successfully: %v\n", updateResp.Success)

	// Reverse string
	reverseReq := request.ReverseRequest{
		Info: "hello world",
	}

	reverseResp, err := client.Reverse(ctx, reverseReq)
	if err != nil {
		fmt.Printf("Reverse failed: %v\n", err)
		return
	}
	fmt.Printf("Reverse successful: %s\n", reverseResp.Info)

	// Delete example
	deleteResp, err := client.DeleteExample(ctx, exampleID)
	if err != nil {
		fmt.Printf("Failed to delete example: %v\n", err)
		return
	}
	fmt.Printf("Example deleted successfully: %v\n", deleteResp.Success)
}

// ExampleWithAuth demonstrates usage example with authentication
func ExampleWithAuth() {
	// Create configuration with authentication
	config := config2.PlatformConfig{
		Addr:       "https://api.example.com",
		Timeout:    30 * time.Second,
		RetryCount: 3,
		APIKey:     "your-api-key",
		AuthToken:  "your-auth-token",
		Debug:      true,
	}

	// Initialize client
	client, err := NewPlatformClient(config)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	ctx := context.Background()

	// Send authenticated request
	getExamplesReq := request.GetExamplesRequest{
		Page:  1,
		Limit: 10,
	}

	examplesResp, err := client.GetExamples(ctx, getExamplesReq)
	if err != nil {
		fmt.Printf("Authenticated request failed: %v\n", err)
		return
	}
	fmt.Printf("Authenticated request successful, got %d examples\n", examplesResp.Total)
}

// ExampleGenericUsage demonstrates how to use generic API
func ExampleGenericUsage() {
	config := config2.GetDefaultConfig()
	client, err := NewPlatformClient(config)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	ctx := context.Background()

	// Use generic GET request
	resp, err := client.Get(ctx, "/api/health", nil)
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		return
	}
	fmt.Printf("Health check response: %s\n", string(resp))

	// Use generic POST request
	customData := map[string]any{
		"action": "ping",
		"target": "server",
	}

	resp, err = client.Post(ctx, "/api/custom", customData)
	if err != nil {
		fmt.Printf("Custom request failed: %v\n", err)
		return
	}
	fmt.Printf("Custom request response: %s\n", string(resp))
}

// ExampleErrorHandling demonstrates error handling examples
func ExampleErrorHandling() {
	config := config2.GetDefaultConfig()
	client, err := NewPlatformClient(config)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	ctx := context.Background()

	// Try to get non-existent example
	_, err = client.GetExample(ctx, 99999)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
		// Handle specific error types here
	}

	// Try to create invalid example
	invalidExampleReq := request.CreateExampleRequest{
		Info: "", // Invalid empty info
		Val:  0,  // Invalid zero value
	}

	_, err = client.CreateExample(ctx, invalidExampleReq)
	if err != nil {
		fmt.Printf("Expected create example error: %v\n", err)
		// Handle specific error types here
	}

	// If you need to access raw response, you can use generic API
	resp, err := client.Get(ctx, "/api/examples/99999", nil)
	if err != nil {
		fmt.Printf("Generic request error: %v\n", err)
	} else {
		fmt.Printf("Raw response: %s\n", string(resp))
		// Can manually parse response
		commonResp, parseErr := response.ParseRawResponse(resp)
		if parseErr != nil {
			fmt.Printf("Failed to parse response: %v\n", parseErr)
		} else {
			fmt.Printf("Parsed response: Code=%d, Status=%s\n",
				commonResp.Code, commonResp.Status)
		}
	}
}
