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

func TestStorageMoverSmbMountReadinessOwnership(t *testing.T) {
	resource := storagemover.StorageMoverSmbMountEndpointResource{}
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
			client.Client.SetTransport(smbMountEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
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
					Body:       io.NopCloser(strings.NewReader(`{"properties":{"endpointType":"NfsMount","host":"server.example","export":"/share","nfsVersion":"NFSv3"}}`)),
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
}

func TestStorageMoverSmbMountReadinessImportParent(t *testing.T) {
	const prefix = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/"
	for _, test := range []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "foreign-provider",
			id:        prefix + "Microsoft.Storage/storageMovers/mover/endpoints/endpoint",
			wantError: true,
		},
		{
			name:      "foreign-parent-type",
			id:        prefix + "Microsoft.StorageMover/projects/mover/endpoints/endpoint",
			wantError: true,
		},
		{
			name:      "missing-parent",
			id:        prefix + "Microsoft.StorageMover/endpoints/endpoint",
			wantError: true,
		},
		{
			name:      "extra-child",
			id:        prefix + "Microsoft.StorageMover/storageMovers/mover/endpoints/endpoint/children/child",
			wantError: true,
		},
		{
			name: "valid-different-mover",
			id:   prefix + "Microsoft.StorageMover/storageMovers/another-mover/endpoints/endpoint",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			requests := 0
			client.Client.SetTransport(smbMountEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				if request.Method != http.MethodGet || request.URL.Path != test.id {
					return nil, fmt.Errorf("unexpected request: %s %s", request.Method, request.URL.Path)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"properties":{"endpointType":"SmbMount","host":"server.example","shareName":"share"}}`)),
					Request:    request,
				}, nil
			}))
			resource := storagemover.StorageMoverSmbMountEndpointResource{}
			metadata := sdk.NewResourceMetaData(&clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}, resource)
			metadata.ResourceData.SetId(test.id)
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			states, err := sdk.WrappedResource(resource).Importer.StateContext(ctx, metadata.ResourceData, metadata.Client)
			if test.wantError {
				if err == nil {
					t.Fatal("expected parent ID validation failure")
				}
				if requests != 0 {
					t.Errorf("invalid parent ID reached the API: %d requests", requests)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if requests != 1 || len(states) != 1 || states[0] == nil || states[0].Id() != test.id {
				t.Fatalf("expected one GET and one imported state for %q, got requests=%d states=%v", test.id, requests, states)
			}
			wantParent := prefix + "Microsoft.StorageMover/storageMovers/another-mover"
			if actual := states[0].Get("storage_mover_id"); actual != wantParent {
				t.Errorf("expected parent %q, got %v", wantParent, actual)
			}
			identity, err := states[0].Identity()
			if err != nil {
				t.Fatal(err)
			}
			if actual := identity.Get("storage_mover_name"); actual != "another-mover" {
				t.Errorf("expected imported parent identity another-mover, got %v", actual)
			}
		})
	}
}

func TestStorageMoverSmbMountEndpointCredentialValidation(t *testing.T) {
	arguments := storagemover.StorageMoverSmbMountEndpointResource{}.Arguments()
	for _, test := range []struct {
		name      string
		field     string
		value     string
		wantError bool
	}{
		{
			name:  "versionless-username",
			field: "username_key_vault_secret_id",
			value: "https://example.vault.azure.net/secrets/username",
		},
		{
			name:      "versioned-username",
			field:     "username_key_vault_secret_id",
			value:     "https://example.vault.azure.net/secrets/username/00000000000000000000000000000001",
			wantError: true,
		},
		{
			name:      "username-key-not-secret",
			field:     "username_key_vault_secret_id",
			value:     "https://example.vault.azure.net/keys/username",
			wantError: true,
		},
		{
			name:  "versionless-password",
			field: "password_key_vault_secret_id",
			value: "https://example.vault.azure.net/secrets/password",
		},
		{
			name:  "versioned-password",
			field: "password_key_vault_secret_id",
			value: "https://example.vault.azure.net/secrets/password/00000000000000000000000000000001",
		},
		{
			name:      "password-key-not-secret",
			field:     "password_key_vault_secret_id",
			value:     "https://example.vault.azure.net/keys/password",
			wantError: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, errors := arguments[test.field].ValidateFunc(test.value, test.field)
			if (len(errors) != 0) != test.wantError {
				t.Fatalf("expected validation error=%t, got %v", test.wantError, errors)
			}
		})
	}
}
