// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover"
	moverclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/client"
)

type smbMountEndpointOwnershipTransport func(*http.Request) (*http.Response, error)

func (f smbMountEndpointOwnershipTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestStorageMoverSmbMountEndpointCredentialUpdate(t *testing.T) {
	const oldUsername = "https://example.vault.azure.net/secrets/old-username"
	const oldPassword = "https://example.vault.azure.net/secrets/old-password"
	const newUsername = "https://example.vault.azure.net/secrets/new-username"
	const newPassword = "https://example.vault.azure.net/secrets/new-password"

	resource := storagemover.StorageMoverSmbMountEndpointResource{}
	wrapped := sdk.WrappedResource(resource)
	id := endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "rg", "mover", "endpoint")
	old := schema.TestResourceDataRaw(t, wrapped.Schema, map[string]any{
		"name": "endpoint", "storage_mover_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.StorageMover/storageMovers/mover",
		"host": "server.example", "share_name": "share",
		"username_key_vault_secret_id": oldUsername, "password_key_vault_secret_id": oldPassword,
	})
	old.SetId(id.ID())
	config := terraform.NewResourceConfigRaw(map[string]any{
		"name": "endpoint", "storage_mover_id": old.Get("storage_mover_id"),
		"host": "server.example", "share_name": "share",
		"username_key_vault_secret_id": newUsername, "password_key_vault_secret_id": newPassword,
	})
	diff, err := wrapped.SimpleDiff(context.Background(), old.State(), config, nil)
	if err != nil {
		t.Fatal(err)
	}
	data, err := schema.InternalMap(wrapped.Schema).Data(old.State(), diff)
	if err != nil {
		t.Fatal(err)
	}
	if !data.HasChanges("username_key_vault_secret_id", "password_key_vault_secret_id") {
		t.Fatal("expected credential change in the update plan")
	}

	client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
	if err != nil {
		t.Fatal(err)
	}
	client.Client.AuthorizeRequest = nil
	client.Client.DisableRetries = true
	requests := 0
	client.Client.SetTransport(smbMountEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
		requests++
		body := `{"properties":{"endpointType":"SmbMount","host":"server.example","shareName":"share","credentials":{"type":"AzureKeyVaultSmb","usernameUri":"` + oldUsername + `","passwordUri":"` + oldPassword + `"}}}`
		if request.Method == http.MethodPut {
			payload := struct {
				Properties struct {
					Credentials endpoints.AzureKeyVaultSmbCredentials `json:"credentials"`
				} `json:"properties"`
			}{}
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if payload.Properties.Credentials.UsernameUri == nil || *payload.Properties.Credentials.UsernameUri != newUsername ||
				payload.Properties.Credentials.PasswordUri == nil || *payload.Properties.Credentials.PasswordUri != newPassword {
				t.Fatal("credential update did not send both new secret URIs")
			}
			body = `{"properties":{"endpointType":"SmbMount","host":"server.example","shareName":"share"}}`
		} else if request.Method != http.MethodGet {
			t.Fatalf("unexpected request: %s", request.Method)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    request,
		}, nil
	}))
	metadata := sdk.NewResourceMetaData(&clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}, resource)
	metadata.ResourceData = data
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := resource.Update().Func(ctx, metadata); err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("expected GET and PUT, got %d requests", requests)
	}
}

func TestStorageMoverSmbMountEndpointOwnership(t *testing.T) {
	for _, operation := range []string{"import", "read", "update", "delete", "import-valid", "read-valid", "update-valid", "delete-valid", "read-not-found", "delete-not-found"} {
		t.Run(operation, func(t *testing.T) {
			client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			reads, writes := 0, 0
			body := `{"properties":{"endpointType":"AzureStorageBlobContainer","storageAccountResourceId":"/subscriptions/00000000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account","blobContainerName":"container"}}`
			status := http.StatusOK
			if strings.HasSuffix(operation, "-valid") {
				body = `{"properties":{"endpointType":"SmbMount","host":"server.example","shareName":"share"}}`
			}
			if strings.HasSuffix(operation, "-not-found") {
				status = http.StatusNotFound
				body = `{"error":{"code":"ResourceNotFound","message":"not found"}}`
			}
			client.Client.SetTransport(smbMountEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
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
					Header: http.Header{
						"Content-Type":                   []string{"application/json"},
						sdkclient.SkipPollingDelayHeader: []string{"true"},
					},
					Body:    io.NopCloser(strings.NewReader(responseBody)),
					Request: request,
				}, nil
			}))
			resource := storagemover.StorageMoverSmbMountEndpointResource{}
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
			if operation == "read-valid" && metadata.ResourceData.Get("host") != "server.example" {
				t.Errorf("matching endpoint was not read into state")
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
