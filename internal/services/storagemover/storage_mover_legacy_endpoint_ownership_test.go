// Copyright IBM Corp. 2014, 2026
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

type legacyEndpointOwnershipTransport func(*http.Request) (*http.Response, error)

func (f legacyEndpointOwnershipTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestStorageMoverLegacyEndpointOwnership(t *testing.T) {
	for _, test := range []struct {
		name      string
		resource  sdk.Resource
		validBody string
	}{
		{
			name:      "nfs-source",
			resource:  storagemover.StorageMoverSourceEndpointResource{},
			validBody: `{"properties":{"endpointType":"NfsMount","host":"server.example","export":"/share","nfsVersion":"NFSv3"}}`,
		},
		{
			name:      "blob-target",
			resource:  storagemover.StorageMoverTargetEndpointResource{},
			validBody: `{"properties":{"endpointType":"AzureStorageBlobContainer","storageAccountResourceId":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account","blobContainerName":"container"}}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, operation := range []string{"import-foreign", "read-foreign", "update-foreign", "delete-foreign", "read-valid", "update-valid", "delete-valid", "delete-not-found"} {
				t.Run(operation, func(t *testing.T) {
					client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
					if err != nil {
						t.Fatal(err)
					}
					client.Client.AuthorizeRequest = nil
					client.Client.DisableRetries = true
					body := `{"properties":{"endpointType":"AzureStorageNfsFileShare","fileShareName":"share","storageAccountResourceId":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account"}}`
					status := http.StatusOK
					if strings.HasSuffix(operation, "-valid") {
						body = test.validBody
					}
					if operation == "delete-not-found" {
						status = http.StatusNotFound
						body = `{"error":{"code":"ResourceNotFound","message":"not found"}}`
					}
					reads, writes := 0, 0
					client.Client.SetTransport(legacyEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
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
					metadata := sdk.NewResourceMetaData(&clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}, test.resource)
					metadata.SetID(endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "rg", "mover", "endpoint"))
					ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
					defer cancel()
					switch operation {
					case "import-foreign":
						_, err = sdk.WrappedResource(test.resource).Importer.StateContext(ctx, metadata.ResourceData, metadata.Client)
					case "read-foreign", "read-valid":
						err = test.resource.Read().Func(ctx, metadata)
					case "update-foreign", "update-valid":
						updatable, ok := test.resource.(sdk.ResourceWithUpdate)
						if !ok {
							t.Fatal("test resource must support update")
						}
						err = updatable.Update().Func(ctx, metadata)
					case "delete-foreign", "delete-valid", "delete-not-found":
						err = test.resource.Delete().Func(ctx, metadata)
					}
					wantError := strings.HasSuffix(operation, "-foreign")
					if wantError && (err == nil || !strings.Contains(err.Error(), "expected ")) {
						t.Errorf("expected endpoint-type rejection, got %v", err)
					}
					if !wantError && err != nil {
						t.Errorf("expected success, got %v", err)
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
		})
	}
}
