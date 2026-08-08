// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileservices"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

func TestDataPlaneFileShareAvailabilityPoller(t *testing.T) {
	accountId := commonids.NewStorageAccountID("12345678-1234-1234-1234-123456789012", "resource-group", "storageaccount")

	testCases := []struct {
		name           string
		statusCode     int
		errorCode      string
		responseError  error
		expectedStatus pollers.PollingStatus
		expectError    bool
	}{
		{
			name:           "available",
			statusCode:     http.StatusOK,
			expectedStatus: pollers.PollingStatusSucceeded,
		},
		{
			name:           "not yet available",
			statusCode:     http.StatusNotFound,
			errorCode:      fileServiceResourceNotFoundErrorCode,
			responseError:  errors.New(fileServiceResourceNotFoundErrorCode),
			expectedStatus: pollers.PollingStatusInProgress,
		},
		{
			name:          "different not found error",
			statusCode:    http.StatusNotFound,
			errorCode:     "ParentResourceNotFound",
			responseError: errors.New("ParentResourceNotFound"),
			expectError:   true,
		},
		{
			name:          "permanent error",
			statusCode:    http.StatusForbidden,
			errorCode:     "AuthorizationFailed",
			responseError: errors.New("AuthorizationFailed"),
			expectError:   true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var responseOData *odata.OData
			if testCase.errorCode != "" {
				responseOData = &odata.OData{
					Error: &odata.Error{Code: &testCase.errorCode},
				}
			}

			poller := DataPlaneFileShareAvailabilityPoller{
				client: &mockFileServicesClient{
					response: fileservices.GetServicePropertiesOperationResponse{
						HttpResponse: &http.Response{StatusCode: testCase.statusCode},
						OData:        responseOData,
					},
					err: testCase.responseError,
				},
				storageAccountId: accountId,
			}

			result, err := poller.Poll(context.Background())
			if testCase.expectError {
				if err == nil {
					t.Fatal("expected a permanent error, got nil")
				}
				var pollingFailedError pollers.PollingFailedError
				if !errors.As(err, &pollingFailedError) {
					t.Fatalf("expected a polling failure, got %T: %v", err, err)
				}
				if result != nil {
					t.Fatalf("expected no poll result, got %#v", result)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if result == nil {
				t.Fatal("expected a poll result, got nil")
			}
			if result.Status != testCase.expectedStatus {
				t.Fatalf("expected status %q, got %q", testCase.expectedStatus, result.Status)
			}
		})
	}
}

type mockFileServicesClient struct {
	response fileservices.GetServicePropertiesOperationResponse
	err      error
}

func (m mockFileServicesClient) GetServiceProperties(context.Context, commonids.StorageAccountId) (fileservices.GetServicePropertiesOperationResponse, error) {
	return m.response, m.err
}
