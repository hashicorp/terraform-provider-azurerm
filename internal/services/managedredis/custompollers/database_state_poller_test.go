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
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

func TestDatabaseStatePoller_Success(t *testing.T) {
	id := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "default")

	mockClient := &mockDatabasesClient{
		getResponse: &databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &databases.Database{
				Properties: &databases.DatabaseProperties{
					ResourceState: pointer.To(databases.ResourceStateRunning),
				},
			},
		},
	}

	pollerType := &dbStatePoller{
		client: mockClient,
		id:     id,
	}

	result, err := pollerType.Poll(context.Background())
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

func TestDatabaseStatePoller_InProgress(t *testing.T) {
	testCases := []struct {
		name  string
		state databases.ResourceState
	}{
		{"Creating", databases.ResourceStateCreating},
		{"Updating", databases.ResourceStateUpdating},
		{"Enabling", databases.ResourceStateEnabling},
		{"Deleting", databases.ResourceStateDeleting},
		{"Disabling", databases.ResourceStateDisabling},
		{"Moving", databases.ResourceStateMoving},
		{"Scaling", databases.ResourceStateScaling},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "default")

			mockClient := &mockDatabasesClient{
				getResponse: &databases.GetOperationResponse{
					HttpResponse: &http.Response{StatusCode: 200},
					Model: &databases.Database{
						Properties: &databases.DatabaseProperties{
							ResourceState: pointer.To(tc.state),
						},
					},
				},
			}

			pollerType := &dbStatePoller{
				client: mockClient,
				id:     id,
			}

			result, err := pollerType.Poll(context.Background())
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

func TestDatabaseStatePoller_UnexpectedState(t *testing.T) {
	id := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "default")

	mockClient := &mockDatabasesClient{
		getResponse: &databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &databases.Database{
				Properties: &databases.DatabaseProperties{
					ResourceState: pointer.To(databases.ResourceState("UnexpectedState")),
				},
			},
		},
	}

	pollerType := &dbStatePoller{
		client: mockClient,
		id:     id,
	}

	result, err := pollerType.Poll(context.Background())

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

func TestDatabaseStatePoller_GetError(t *testing.T) {
	id := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "default")

	mockClient := &mockDatabasesClient{
		getError: fmt.Errorf("API error"),
	}

	pollerType := &dbStatePoller{
		client: mockClient,
		id:     id,
	}

	result, err := pollerType.Poll(context.Background())

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

func TestDatabaseStatePoller_StateTransitions(t *testing.T) {
	id := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "default")

	callCount := 0
	mockClient := &statefulMockDatabaseClient{
		callCount: &callCount,
	}

	pollerType := NewDBStatePoller(mockClient, id)

	result1, err1 := pollerType.Poll(context.Background())
	if err1 != nil {
		t.Fatalf("first poll error: %v", err1)
	}
	if result1.Status != pollers.PollingStatusInProgress {
		t.Fatalf("expected in progress, got %s", result1.Status)
	}

	result2, err2 := pollerType.Poll(context.Background())
	if err2 != nil {
		t.Fatalf("second poll error: %v", err2)
	}
	if result2.Status != pollers.PollingStatusInProgress {
		t.Fatalf("expected in progress, got %s", result2.Status)
	}

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

type mockDatabasesClient struct {
	getResponse *databases.GetOperationResponse
	getError    error
}

func (m *mockDatabasesClient) Get(ctx context.Context, id databases.DatabaseId) (databases.GetOperationResponse, error) {
	if m.getError != nil && m.getResponse != nil {
		return *m.getResponse, m.getError
	}
	if m.getError != nil {
		return databases.GetOperationResponse{}, m.getError
	}
	if m.getResponse != nil {
		return *m.getResponse, nil
	}
	return databases.GetOperationResponse{}, fmt.Errorf("no mock response configured")
}

type statefulMockDatabaseClient struct {
	callCount *int
}

func (m *statefulMockDatabaseClient) Get(ctx context.Context, id databases.DatabaseId) (databases.GetOperationResponse, error) {
	*m.callCount++

	var state databases.ResourceState
	switch *m.callCount {
	case 1:
		state = databases.ResourceStateCreating
	case 2:
		state = databases.ResourceStateUpdating
	case 3:
		state = databases.ResourceStateRunning
	default:
		state = databases.ResourceStateRunning
	}

	return databases.GetOperationResponse{
		HttpResponse: &http.Response{StatusCode: 200},
		Model: &databases.Database{
			Properties: &databases.DatabaseProperties{
				ResourceState: pointer.To(state),
			},
		},
	}, nil
}
