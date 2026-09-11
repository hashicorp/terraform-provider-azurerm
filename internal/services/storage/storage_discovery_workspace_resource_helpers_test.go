// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
)

func TestStorageDiscoveryWorkspaceCreateImportCheck(t *testing.T) {
	id := storagediscoveryworkspaces.NewProviderStorageDiscoveryWorkspaceID("00000000-0000-0000-0000-000000000000", "test", "test")
	resource := StorageDiscoveryWorkspaceResource{}

	testCases := []struct {
		name           string
		allowOverwrite bool
		getStatus      int
		getRequests    int32
		putRequests    int32
		errorContains  string
	}{
		{
			name:          "existing workspace requires import by default",
			getStatus:     http.StatusOK,
			getRequests:   1,
			errorContains: "needs to be imported",
		},
		{
			name:        "missing workspace is created",
			getStatus:   http.StatusNotFound,
			getRequests: 1,
			putRequests: 1,
		},
		{
			name:           "explicit overwrite skips the presence check",
			allowOverwrite: true,
			getStatus:      http.StatusOK,
			putRequests:    1,
		},
		{
			name:          "presence check errors prevent creation",
			getStatus:     http.StatusForbidden,
			getRequests:   1,
			errorContains: "checking for presence of existing",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var getRequests, putRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != id.ID() {
					t.Errorf("expected request path %q, got %q", id.ID(), req.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				switch req.Method {
				case http.MethodGet:
					getRequests.Add(1)
					w.WriteHeader(testCase.getStatus)
				case http.MethodPut:
					putRequests.Add(1)
					var payload storagediscoveryworkspaces.StorageDiscoveryWorkspace
					if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
						t.Errorf("decoding create payload: %v", err)
					} else if payload.Properties == nil || len(payload.Properties.Scopes) != 1 ||
						!reflect.DeepEqual(payload.Properties.Scopes[0].ResourceTypes, []storagediscoveryworkspaces.StorageDiscoveryResourceType{
							storagediscoveryworkspaces.StorageDiscoveryResourceTypeMicrosoftPointStorageStorageAccounts,
						}) {
						t.Errorf("unexpected create payload: %#v", payload.Properties)
					}
					w.WriteHeader(http.StatusOK)
				default:
					t.Errorf("unexpected request method %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
				}

				if err := json.NewEncoder(w).Encode(storagediscoveryworkspaces.StorageDiscoveryWorkspace{Location: "westus2"}); err != nil {
					t.Errorf("encoding response: %v", err)
				}
			}))
			defer server.Close()

			client, err := storagediscoveryworkspaces.NewStorageDiscoveryWorkspacesClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
			if err != nil {
				t.Fatalf("building client: %v", err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true

			metadata := sdk.NewResourceMetaData(&clients.Client{
				Account: &clients.ResourceManagerAccount{SubscriptionId: id.SubscriptionId},
				Features: features.UserFeatures{
					SkipImportCheckOnCreateAndAllowOverwritingExistingResources: testCase.allowOverwrite,
				},
				Storage: &storageclient.Client{StorageDiscoveryWorkspacesClient: client},
			}, resource)
			if err := metadata.Encode(&StorageDiscoveryWorkspaceModel{
				Name:              "test",
				ResourceGroupName: "test",
				Location:          "westus2",
				WorkspaceRoots:    []string{"/subscriptions/00000000-0000-0000-0000-000000000000"},
				Sku:               string(storagediscoveryworkspaces.StorageDiscoverySkuStandard),
				Scope: []StorageDiscoveryScopeModel{{
					DisplayName:   "TestScope",
					ResourceTypes: []string{"Microsoft.Storage/storageAccounts"},
				}},
			}); err != nil {
				t.Fatalf("encoding resource data: %v", err)
			}

			ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
			defer cancel()
			err = resource.Create().Func(ctx, metadata)
			if testCase.errorContains != "" {
				if err == nil || !strings.Contains(err.Error(), testCase.errorContains) {
					t.Errorf("expected error containing %q, got %v", testCase.errorContains, err)
				}
				if actual := metadata.ResourceData.Id(); actual != "" {
					t.Errorf("expected no resource ID on failure, got %q", actual)
				}
			} else {
				if err != nil {
					t.Errorf("creating workspace: %v", err)
				}
				if actual := metadata.ResourceData.Id(); actual != id.ID() {
					t.Errorf("expected resource ID %q, got %q", id.ID(), actual)
				}
			}
			if actual := getRequests.Load(); actual != testCase.getRequests {
				t.Errorf("expected %d GET requests, got %d", testCase.getRequests, actual)
			}
			if actual := putRequests.Load(); actual != testCase.putRequests {
				t.Errorf("expected %d PUT requests, got %d", testCase.putRequests, actual)
			}
		})
	}
}

func TestStorageDiscoveryScopeResourceTypesConversion(t *testing.T) {
	testCases := []struct {
		name      string
		input     []string
		apiValues []storagediscoveryworkspaces.StorageDiscoveryResourceType
		expected  []string
	}{
		{
			name:      "nil",
			apiValues: []storagediscoveryworkspaces.StorageDiscoveryResourceType{},
			expected:  []string{},
		},
		{
			name:      "empty",
			input:     []string{},
			apiValues: []storagediscoveryworkspaces.StorageDiscoveryResourceType{},
			expected:  []string{},
		},
		{
			name:      "storage accounts",
			input:     []string{"Microsoft.Storage/storageAccounts"},
			apiValues: []storagediscoveryworkspaces.StorageDiscoveryResourceType{storagediscoveryworkspaces.StorageDiscoveryResourceTypeMicrosoftPointStorageStorageAccounts},
			expected:  []string{"Microsoft.Storage/storageAccounts"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			expanded := expandStorageDiscoveryScopes([]StorageDiscoveryScopeModel{{ResourceTypes: testCase.input}})
			if !reflect.DeepEqual(expanded[0].ResourceTypes, testCase.apiValues) {
				t.Errorf("expected expanded resource types %#v, got %#v", testCase.apiValues, expanded[0].ResourceTypes)
			}
			apiValues := testCase.apiValues
			if testCase.input == nil {
				apiValues = nil
			}
			flattened := flattenStorageDiscoveryScopes([]storagediscoveryworkspaces.StorageDiscoveryScope{{ResourceTypes: apiValues}})
			if !reflect.DeepEqual(flattened[0].ResourceTypes, testCase.expected) {
				t.Errorf("expected flattened resource types %#v, got %#v", testCase.expected, flattened[0].ResourceTypes)
			}
		})
	}
}

func TestStorageDiscoveryWorkspaceResourceTypesValidation(t *testing.T) {
	resource := sdk.WrappedResource(StorageDiscoveryWorkspaceResource{})
	testCases := []struct {
		name  string
		input []interface{}
		valid bool
	}{
		{name: "missing"},
		{name: "empty set", input: []interface{}{}},
		{name: "empty string", input: []interface{}{""}},
		{name: "invalid value", input: []interface{}{"Microsoft.Storage/storageAccounts/blobServices"}},
		{name: "incorrect casing", input: []interface{}{"microsoft.storage/storageaccounts"}},
		{name: "valid", input: []interface{}{"Microsoft.Storage/storageAccounts"}, valid: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := storageDiscoveryWorkspaceConfig([]interface{}{
				storageDiscoveryScopeConfig("TestScope", testCase.input, []interface{}{}, map[string]interface{}{}),
			})
			diagnostics := resource.Validate(terraform.NewResourceConfigRaw(config))
			if diagnostics.HasError() == testCase.valid {
				t.Fatalf("expected valid %t, got diagnostics: %v", testCase.valid, diagnostics)
			}
		})
	}
}
