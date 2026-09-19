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

type smbMountEndpointOwnershipTransport func(*http.Request) (*http.Response, error)

func (f smbMountEndpointOwnershipTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestStorageMoverSmbMountEndpointOwnership(t *testing.T) {
	for _, operation := range []string{"import", "read", "update", "delete", "read-valid", "delete-not-found"} {
		t.Run(operation, func(t *testing.T) {
			client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			writes := 0
			body := `{"properties":{"endpointType":"AzureStorageBlobContainer","storageAccountResourceId":"/subscriptions/00000000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account","blobContainerName":"container"}}`
			status := http.StatusOK
			if operation == "read-valid" {
				body = `{"properties":{"endpointType":"SmbMount","host":"server.example","shareName":"share"}}`
			}
			if operation == "delete-not-found" {
				status = http.StatusNotFound
				body = `{"error":{"code":"ResourceNotFound","message":"not found"}}`
			}
			client.Client.SetTransport(smbMountEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet {
					writes++
					return nil, fmt.Errorf("unexpected mutation: %s", request.Method)
				}
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    request,
				}, nil
			}))
			resource := storagemover.StorageMoverSmbMountEndpointResource{}
			metadata := sdk.NewResourceMetaData(&clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}, resource)
			metadata.SetID(endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "rg", "mover", "endpoint"))
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			switch operation {
			case "import":
				_, err = sdk.WrappedResource(resource).Importer.StateContext(ctx, metadata.ResourceData, metadata.Client)
			case "read", "read-valid":
				err = resource.Read().Func(ctx, metadata)
			case "update":
				err = resource.Update().Func(ctx, metadata)
			case "delete", "delete-not-found":
				err = resource.Delete().Func(ctx, metadata)
			}
			wantError := operation != "read-valid" && operation != "delete-not-found"
			if wantError && (err == nil || !strings.Contains(err.Error(), "expected")) {
				t.Errorf("expected endpoint-type rejection, got %v", err)
			}
			if !wantError && err != nil {
				t.Errorf("expected success, got %v", err)
			}
			if operation == "read-valid" && metadata.ResourceData.Get("host") != "server.example" {
				t.Errorf("matching endpoint was not read into state")
			}
			if writes != 0 {
				t.Errorf("foreign endpoint received %d mutation requests", writes)
			}
		})
	}
}
