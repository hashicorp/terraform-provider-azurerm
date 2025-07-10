// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type mockRedisEnterpriseClient struct {
	getResponse *redisenterprise.GetOperationResponse
	getError    error
}

func (m *mockRedisEnterpriseClient) Get(ctx context.Context, id redisenterprise.RedisEnterpriseId) (redisenterprise.GetOperationResponse, error) {
	if m.getError != nil {
		return redisenterprise.GetOperationResponse{}, m.getError
	}
	if m.getResponse != nil {
		return *m.getResponse, nil
	}
	return redisenterprise.GetOperationResponse{}, fmt.Errorf("no mock response configured")
}

func TestClusterStatePoller_Success(t *testing.T) {
	// Arrange
	id := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	mockClient := &mockRedisEnterpriseClient{
		getResponse: &redisenterprise.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &redisenterprise.Cluster{
				Properties: &redisenterprise.ClusterProperties{
					ResourceState: pointer.To(redisenterprise.ResourceStateRunning),
				},
			},
		},
	}

	pollerType := &clusterStatePoller{
		client: mockClient,
		id:     id,
	}

	// Act
	result, err := pollerType.Poll(context.Background())

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Status != pollers.PollingStatusSucceeded {
		t.Fatalf("expected status %s, got %s", pollers.PollingStatusSucceeded, result.Status)
	}
}

func TestClusterStatePoller_InProgress(t *testing.T) {
	testCases := []struct {
		name  string
		state redisenterprise.ResourceState
	}{
		{"Creating", redisenterprise.ResourceStateCreating},
		{"Updating", redisenterprise.ResourceStateUpdating},
		{"Enabling", redisenterprise.ResourceStateEnabling},
		{"Deleting", redisenterprise.ResourceStateDeleting},
		{"Disabling", redisenterprise.ResourceStateDisabling},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			id := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

			mockClient := &mockRedisEnterpriseClient{
				getResponse: &redisenterprise.GetOperationResponse{
					HttpResponse: &http.Response{StatusCode: 200},
					Model: &redisenterprise.Cluster{
						Properties: &redisenterprise.ClusterProperties{
							ResourceState: pointer.To(tc.state),
						},
					},
				},
			}

			pollerType := &clusterStatePoller{
				client: mockClient,
				id:     id,
			}

			// Act
			result, err := pollerType.Poll(context.Background())

			// Assert
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if result == nil {
				t.Fatal("expected result, got nil")
			}
			if result.Status != pollers.PollingStatusInProgress {
				t.Fatalf("expected status %s, got %s", pollers.PollingStatusInProgress, result.Status)
			}
			if result.PollInterval != 15*time.Second {
				t.Fatalf("expected poll interval %v, got %v", 15*time.Second, result.PollInterval)
			}
		})
	}
}

func TestClusterStatePoller_UnexpectedState(t *testing.T) {
	// Arrange
	id := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	mockClient := &mockRedisEnterpriseClient{
		getResponse: &redisenterprise.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &redisenterprise.Cluster{
				Properties: &redisenterprise.ClusterProperties{
					ResourceState: pointer.To(redisenterprise.ResourceState("UnexpectedState")),
				},
			},
		},
	}

	pollerType := &clusterStatePoller{
		client: mockClient,
		id:     id,
	}

	// Act
	result, err := pollerType.Poll(context.Background())

	// Assert
	if err == nil {
		t.Fatal("expected error for unexpected state, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
	expectedError := fmt.Sprintf("unexpected resource state for %s: UnexpectedState", id)
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestClusterStatePoller_GetError(t *testing.T) {
	// Arrange
	id := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	mockClient := &mockRedisEnterpriseClient{
		getError: fmt.Errorf("API error"),
	}

	pollerType := &clusterStatePoller{
		client: mockClient,
		id:     id,
	}

	// Act
	result, err := pollerType.Poll(context.Background())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
	expectedError := fmt.Sprintf("retrieving %s: API error", id)
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestClusterStatePoller_StateTransitions(t *testing.T) {
	id := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	callCount := 0
	mockClient := &statefulMockClient{
		callCount: &callCount,
	}

	pollerType := NewClusterStatePoller(mockClient, id)

	// Test first call - should be "Creating" and return in progress
	result1, err1 := pollerType.Poll(context.Background())
	if err1 != nil {
		t.Fatalf("first poll error: %v", err1)
	}
	if result1.Status != pollers.PollingStatusInProgress {
		t.Fatalf("expected in progress, got %s", result1.Status)
	}

	// Test second call - should be "Updating" and return in progress
	result2, err2 := pollerType.Poll(context.Background())
	if err2 != nil {
		t.Fatalf("second poll error: %v", err2)
	}
	if result2.Status != pollers.PollingStatusInProgress {
		t.Fatalf("expected in progress, got %s", result2.Status)
	}

	// Test third call - should be "Running" and return success
	result3, err3 := pollerType.Poll(context.Background())
	if err3 != nil {
		t.Fatalf("third poll error: %v", err3)
	}
	if result3.Status != pollers.PollingStatusSucceeded {
		t.Fatalf("expected succeeded, got %s", result3.Status)
	}

	if callCount != 3 {
		t.Fatalf("expected 3 calls, got %d", callCount)
	}
}

type statefulMockClient struct {
	callCount *int
}

func (m *statefulMockClient) Get(ctx context.Context, id redisenterprise.RedisEnterpriseId) (redisenterprise.GetOperationResponse, error) {
	*m.callCount++

	var state redisenterprise.ResourceState
	switch *m.callCount {
	case 1:
		state = redisenterprise.ResourceStateCreating
	case 2:
		state = redisenterprise.ResourceStateUpdating
	case 3:
		state = redisenterprise.ResourceStateRunning
	default:
		state = redisenterprise.ResourceStateRunning
	}

	return redisenterprise.GetOperationResponse{
		HttpResponse: &http.Response{StatusCode: 200},
		Model: &redisenterprise.Cluster{
			Properties: &redisenterprise.ClusterProperties{
				ResourceState: pointer.To(state),
			},
		},
	}, nil
}
