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

type smbTargetEndpointOwnershipTransport func(*http.Request) (*http.Response, error)

func (f smbTargetEndpointOwnershipTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestStorageMoverSmbFileShareTargetEndpointOwnership(t *testing.T) {
	for _, operation := range []string{"import", "read", "update", "delete", "read-not-found", "delete-not-found"} {
		t.Run(operation, func(t *testing.T) {
			client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			writes := 0
			notFound := strings.HasSuffix(operation, "-not-found")
			client.Client.SetTransport(smbTargetEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet {
					writes++
					return nil, fmt.Errorf("unexpected mutation: %s", request.Method)
				}
				status := http.StatusOK
				if notFound {
					status = http.StatusNotFound
				}
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"properties":{"endpointType":"AzureStorageBlobContainer","storageAccountResourceId":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account","blobContainerName":"container"}}`)),
					Request:    request,
				}, nil
			}))
			resource := storagemover.StorageMoverSmbFileShareTargetEndpointResource{}
			metadata := sdk.NewResourceMetaData(&clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}, resource)
			metadata.SetID(endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "rg", "mover", "endpoint"))
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			switch operation {
			case "import":
				_, err = sdk.WrappedResource(resource).Importer.StateContext(ctx, metadata.ResourceData, metadata.Client)
			case "read", "read-not-found":
				err = resource.Read().Func(ctx, metadata)
			case "update":
				err = resource.Update().Func(ctx, metadata)
			case "delete", "delete-not-found":
				err = resource.Delete().Func(ctx, metadata)
			}
			if !notFound && (err == nil || !strings.Contains(err.Error(), "expected")) {
				t.Errorf("expected endpoint-type rejection, got %v", err)
			}
			if notFound && err != nil {
				t.Errorf("expected missing endpoint to succeed, got %v", err)
			}
			if operation == "read-not-found" && metadata.ResourceData.Id() != "" {
				t.Error("missing endpoint was not removed from state")
			}
			if writes != 0 {
				t.Errorf("expected no mutations, got %d", writes)
			}
		})
	}
}
