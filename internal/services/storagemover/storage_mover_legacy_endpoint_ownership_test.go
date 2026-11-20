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
	for _, resource := range []sdk.ResourceWithUpdate{
		storagemover.StorageMoverSourceEndpointResource{},
		storagemover.StorageMoverTargetEndpointResource{},
	} {
		t.Run(resource.ResourceType(), func(t *testing.T) {
			for _, operation := range []string{"import", "read", "update", "delete", "delete-not-found"} {
				t.Run(operation, func(t *testing.T) {
					client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
					if err != nil {
						t.Fatal(err)
					}
					client.Client.AuthorizeRequest = nil
					client.Client.DisableRetries = true
					writes := 0
					client.Client.SetTransport(legacyEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
						if request.Method != http.MethodGet {
							writes++
							return nil, fmt.Errorf("unexpected mutation: %s", request.Method)
						}
						status := http.StatusOK
						if operation == "delete-not-found" {
							status = http.StatusNotFound
						}
						return &http.Response{
							StatusCode: status,
							Header:     http.Header{"Content-Type": []string{"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"properties":{"endpointType":"AzureStorageNfsFileShare","fileShareName":"share","storageAccountResourceId":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account"}}`)),
							Request:    request,
						}, nil
					}))
					metadata := sdk.NewResourceMetaData(&clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}, resource)
					metadata.SetID(endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "rg", "mover", "endpoint"))
					ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
					defer cancel()
					switch operation {
					case "import":
						_, err = sdk.WrappedResource(resource).Importer.StateContext(ctx, metadata.ResourceData, metadata.Client)
					case "read":
						err = resource.Read().Func(ctx, metadata)
					case "update":
						err = resource.Update().Func(ctx, metadata)
					case "delete", "delete-not-found":
						err = resource.Delete().Func(ctx, metadata)
					}
					if operation != "delete-not-found" && (err == nil || !strings.Contains(err.Error(), "expected ")) {
						t.Errorf("expected endpoint-type rejection, got %v", err)
					}
					if operation == "delete-not-found" && err != nil {
						t.Errorf("expected missing endpoint to succeed, got %v", err)
					}
					if writes != 0 {
						t.Errorf("expected no mutations, got %d", writes)
					}
				})
			}
		})
	}
}
