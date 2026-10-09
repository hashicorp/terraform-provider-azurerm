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

type legacySmbMountOwnershipTransport func(*http.Request) (*http.Response, error)

func (t legacySmbMountOwnershipTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return t(request)
}

func TestStorageMoverLegacySmbMountReadinessOwnership(t *testing.T) {
	for _, resource := range []sdk.ResourceWithUpdate{
		storagemover.StorageMoverSourceEndpointResource{},
		storagemover.StorageMoverTargetEndpointResource{},
	} {
		t.Run(resource.ResourceType(), func(t *testing.T) {
			for _, operation := range []string{"import", "read", "update", "delete"} {
				t.Run(operation, func(t *testing.T) {
					client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
					if err != nil {
						t.Fatal(err)
					}
					client.Client.AuthorizeRequest = nil
					client.Client.DisableRetries = true
					reads, writes := 0, 0
					id := endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "rg", "mover", "endpoint")
					client.Client.SetTransport(legacySmbMountOwnershipTransport(func(request *http.Request) (*http.Response, error) {
						if request.Method != http.MethodGet {
							writes++
							return nil, fmt.Errorf("unexpected mutation: %s", request.Method)
						}
						reads++
						if request.URL.Path != id.ID() {
							t.Errorf("expected GET path %q, got %q", id.ID(), request.URL.Path)
						}
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": []string{"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"properties":{"endpointType":"SmbMount","host":"server.example","shareName":"share"}}`)),
							Request:    request,
						}, nil
					}))
					metadata := sdk.NewResourceMetaData(&clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}, resource)
					metadata.SetID(id)
					ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
					defer cancel()
					switch operation {
					case "import":
						_, err = sdk.WrappedResource(resource).Importer.StateContext(ctx, metadata.ResourceData, metadata.Client)
					case "read":
						err = resource.Read().Func(ctx, metadata)
					case "update":
						err = resource.Update().Func(ctx, metadata)
					case "delete":
						err = resource.Delete().Func(ctx, metadata)
					}
					if err == nil || !strings.Contains(err.Error(), "expected ") {
						t.Errorf("expected endpoint-type rejection, got %v", err)
					}
					if reads != 1 || writes != 0 {
						t.Errorf("expected 1 read and no mutations, got reads=%d writes=%d", reads, writes)
					}
				})
			}
		})
	}
}
