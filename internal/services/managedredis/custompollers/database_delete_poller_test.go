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
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type mockDatabaseClient struct {
	getResponse *databases.GetOperationResponse
	getError    error
}

func (m *mockDatabaseClient) Get(ctx context.Context, id databases.DatabaseId) (databases.GetOperationResponse, error) {
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

func TestDatabaseDeletePoller_DatabaseNotFound(t *testing.T) {
	databaseId := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "test-db")
	clusterId := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	mockDbClient := &mockDatabaseClient{
		getResponse: &databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 404},
		},
		getError: fmt.Errorf("database not found"),
	}

	mockClusterClient := &mockRedisEnterpriseClient{
		getResponse: &redisenterprise.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &redisenterprise.Cluster{
				Properties: &redisenterprise.ClusterProperties{
					ResourceState: pointer.To(redisenterprise.ResourceStateRunning),
				},
			},
		},
	}

	pollerType := &databaseDeletePoller{
		databaseClient: mockDbClient,
		clusterClient:  mockClusterClient,
		databaseId:     databaseId,
		clusterId:      clusterId,
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

func TestDatabaseDeletePoller_ClusterNotFound(t *testing.T) {
	databaseId := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "test-db")
	clusterId := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	mockDbClient := &mockDatabaseClient{
		getResponse: &databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &databases.Database{
				Properties: &databases.DatabaseProperties{},
			},
		},
	}

	mockClusterClient := &mockRedisEnterpriseClient{
		getResponse: &redisenterprise.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 404},
		},
		getError: fmt.Errorf("cluster not found"),
	}

	pollerType := &databaseDeletePoller{
		databaseClient: mockDbClient,
		clusterClient:  mockClusterClient,
		databaseId:     databaseId,
		clusterId:      clusterId,
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

func TestDatabaseDeletePoller_BothExist(t *testing.T) {
	databaseId := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "test-db")
	clusterId := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	mockDbClient := &mockDatabaseClient{
		getResponse: &databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &databases.Database{
				Properties: &databases.DatabaseProperties{},
			},
		},
	}

	mockClusterClient := &mockRedisEnterpriseClient{
		getResponse: &redisenterprise.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &redisenterprise.Cluster{
				Properties: &redisenterprise.ClusterProperties{
					ResourceState: pointer.To(redisenterprise.ResourceStateRunning),
				},
			},
		},
	}

	pollerType := &databaseDeletePoller{
		databaseClient: mockDbClient,
		clusterClient:  mockClusterClient,
		databaseId:     databaseId,
		clusterId:      clusterId,
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
}

func TestDatabaseDeletePoller_ClusterGetError(t *testing.T) {
	databaseId := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "test-db")
	clusterId := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	mockDbClient := &mockDatabaseClient{
		getResponse: &databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
		},
	}

	mockClusterClient := &mockRedisEnterpriseClient{
		getError: fmt.Errorf("cluster API error"),
	}

	pollerType := &databaseDeletePoller{
		databaseClient: mockDbClient,
		clusterClient:  mockClusterClient,
		databaseId:     databaseId,
		clusterId:      clusterId,
	}

	result, err := pollerType.Poll(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
	expectedError := fmt.Sprintf("retrieving cluster %s: cluster API error", clusterId)
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestDatabaseDeletePoller_DatabaseGetError(t *testing.T) {
	databaseId := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "test-db")
	clusterId := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	mockDbClient := &mockDatabaseClient{
		getError: fmt.Errorf("database API error"),
	}

	mockClusterClient := &mockRedisEnterpriseClient{
		getResponse: &redisenterprise.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &redisenterprise.Cluster{
				Properties: &redisenterprise.ClusterProperties{
					ResourceState: pointer.To(redisenterprise.ResourceStateRunning),
				},
			},
		},
	}

	pollerType := &databaseDeletePoller{
		databaseClient: mockDbClient,
		clusterClient:  mockClusterClient,
		databaseId:     databaseId,
		clusterId:      clusterId,
	}

	result, err := pollerType.Poll(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
	expectedError := fmt.Sprintf("retrieving database %s: database API error", databaseId)
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestDatabaseDeletePoller_DeletionFlow(t *testing.T) {
	databaseId := databases.NewDatabaseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster", "test-db")
	clusterId := redisenterprise.NewRedisEnterpriseID("12345678-1234-1234-1234-123456789012", "test-rg", "test-cluster")

	callCount := 0
	mockDbClient := &statefulMockDatabaseClient{
		callCount: &callCount,
	}

	mockClusterClient := &mockRedisEnterpriseClient{
		getResponse: &redisenterprise.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &redisenterprise.Cluster{
				Properties: &redisenterprise.ClusterProperties{
					ResourceState: pointer.To(redisenterprise.ResourceStateRunning),
				},
			},
		},
	}

	pollerType := NewDatabaseDeletePoller(mockDbClient, mockClusterClient, databaseId, clusterId)

	// Test first call - database should exist and return in progress
	result1, err1 := pollerType.Poll(context.Background())
	if err1 != nil {
		t.Fatalf("first poll error: %v", err1)
	}
	if result1.Status != pollers.PollingStatusInProgress {
		t.Fatalf("expected in progress, got %s", result1.Status)
	}

	// Test second call - database should exist and return in progress
	result2, err2 := pollerType.Poll(context.Background())
	if err2 != nil {
		t.Fatalf("second poll error: %v", err2)
	}
	if result2.Status != pollers.PollingStatusInProgress {
		t.Fatalf("expected in progress, got %s", result2.Status)
	}

	// Test third call - database should be deleted and return success
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

type statefulMockDatabaseClient struct {
	callCount *int
}

func (m *statefulMockDatabaseClient) Get(ctx context.Context, id databases.DatabaseId) (databases.GetOperationResponse, error) {
	*m.callCount++

	switch *m.callCount {
	case 1, 2:
		// Database exists, deletion in progress
		return databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 200},
			Model: &databases.Database{
				Properties: &databases.DatabaseProperties{},
			},
		}, nil
	case 3:
		// Database not found, deletion successful
		return databases.GetOperationResponse{
			HttpResponse: &http.Response{StatusCode: 404},
		}, fmt.Errorf("database not found")
	default:
		return databases.GetOperationResponse{}, fmt.Errorf("unexpected call count: %d", *m.callCount)
	}
}
