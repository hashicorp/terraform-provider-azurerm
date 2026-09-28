// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover"
	moverclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/client"
)

type nfsTargetEndpointOwnershipTransport func(*http.Request) (*http.Response, error)

func (f nfsTargetEndpointOwnershipTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestStorageMoverNfsFileShareTargetEndpointOwnership(t *testing.T) {
	for _, operation := range []string{"import", "read", "update", "delete", "import-valid", "read-valid", "update-valid", "delete-valid", "read-not-found", "delete-not-found"} {
		t.Run(operation, func(t *testing.T) {
			client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			reads, writes := 0, 0
			body := `{"properties":{"endpointType":"AzureStorageBlobContainer","storageAccountResourceId":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account","blobContainerName":"container"}}`
			status := http.StatusOK
			if strings.HasSuffix(operation, "-valid") {
				body = `{"properties":{"endpointType":"AzureStorageNfsFileShare","fileShareName":"share","storageAccountResourceId":"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg/providers/microsoft.storage/storageaccounts/account"}}`
			}
			if strings.HasSuffix(operation, "-not-found") {
				status = http.StatusNotFound
				body = `{"error":{"code":"ResourceNotFound","message":"not found"}}`
			}
			client.Client.SetTransport(nfsTargetEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
				responseStatus, responseBody := status, body
				if request.Method != http.MethodGet {
					writes++
					switch {
					case operation == "update-valid" && request.Method == http.MethodPut:
					case operation == "delete-valid" && request.Method == http.MethodDelete:
						responseStatus, responseBody = http.StatusNoContent, ""
					default:
						return nil, fmt.Errorf("unexpected mutation: %s", request.Method)
					}
				} else {
					reads++
					if operation == "delete-valid" && writes > 0 {
						responseStatus = http.StatusNotFound
						responseBody = `{"error":{"code":"ResourceNotFound","message":"not found"}}`
					}
				}
				return &http.Response{
					StatusCode: responseStatus,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(responseBody)),
					Request:    request,
				}, nil
			}))
			resource := storagemover.StorageMoverNfsFileShareTargetEndpointResource{}
			metadata := sdk.NewResourceMetaData(&clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}, resource)
			metadata.SetID(endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "rg", "mover", "endpoint"))
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			switch operation {
			case "import", "import-valid":
				states, importErr := sdk.WrappedResource(resource).Importer.StateContext(ctx, metadata.ResourceData, metadata.Client)
				err = importErr
				if err == nil && (len(states) != 1 || states[0] == nil || states[0].Id() != metadata.ResourceData.Id()) {
					t.Fatal("expected one imported state with the endpoint ID")
				}
			case "read", "read-valid", "read-not-found":
				err = resource.Read().Func(ctx, metadata)
			case "update", "update-valid":
				err = resource.Update().Func(ctx, metadata)
			case "delete", "delete-valid", "delete-not-found":
				err = resource.Delete().Func(ctx, metadata)
			}
			wantError := !strings.HasSuffix(operation, "-valid") && !strings.HasSuffix(operation, "-not-found")
			if wantError && (err == nil || !strings.Contains(err.Error(), "expected")) {
				t.Errorf("expected endpoint-type rejection, got %v", err)
			}
			if !wantError && err != nil {
				t.Errorf("expected success, got %v", err)
			}
			if operation == "import-valid" || operation == "read-valid" {
				identity, err := metadata.ResourceData.Identity()
				if err != nil {
					t.Fatal(err)
				}
				for key, expected := range map[string]string{
					"subscription_id":     "00000000-0000-0000-0000-000000000000",
					"resource_group_name": "rg",
					"storage_mover_name":  "mover",
					"name":                "endpoint",
				} {
					if actual := identity.Get(key); actual != expected {
						t.Errorf("identity %s: expected %q, got %v", key, expected, actual)
					}
				}
			}
			if operation == "read-valid" && metadata.ResourceData.Get("storage_account_id") != "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account" {
				t.Errorf("storage account ID was not normalized")
			}
			if operation == "read-not-found" && metadata.ResourceData.Id() != "" {
				t.Fatal("missing endpoint was not removed from state")
			}
			wantReads, wantWrites := 1, 0
			if operation == "update-valid" || operation == "delete-valid" {
				wantWrites = 1
			}
			if operation == "delete-valid" {
				wantReads = 2
			}
			if reads != wantReads || writes != wantWrites {
				t.Errorf("expected %d reads and %d mutations, got reads=%d writes=%d", wantReads, wantWrites, reads, writes)
			}
		})
	}
}
