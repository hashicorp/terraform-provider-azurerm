// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover"
	moverclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/client"
)

type nfsTargetEndpointOwnershipTransport func(*http.Request) (*http.Response, error)

func (f nfsTargetEndpointOwnershipTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestStorageMoverNfsFileShareTargetEndpointDescriptionUpdate(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		old         string
		description *string
		server      string
		expected    string
	}{
		{name: "add", description: pointer.To("added"), expected: "added"},
		{name: "change", old: "original", description: pointer.To("updated"), server: "original", expected: "updated"},
		{name: "remove", old: "original", server: "original"},
		{name: "unchanged", old: "original", description: pointer.To("original"), server: "server description", expected: "server description"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			id := endpoints.NewEndpointID("00000000-0000-0000-0000-000000000000", "rg", "mover", "endpoint")
			accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account"
			config := map[string]any{
				"name":               id.EndpointName,
				"storage_mover_id":   endpoints.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
				"file_share_name":    "share",
				"storage_account_id": accountID,
				"description":        testCase.old,
			}
			wrapped := sdk.WrappedResource(storagemover.StorageMoverNfsFileShareTargetEndpointResource{})
			createDiff, err := wrapped.Diff(ctx, nil, terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			attributes, err := createDiff.Apply(nil, wrapped.CoreConfigSchema())
			if err != nil {
				t.Fatal(err)
			}
			state := &terraform.InstanceState{ID: id.ID(), Attributes: attributes}
			delete(config, "description")
			if testCase.description != nil {
				config["description"] = *testCase.description
			}
			diff, err := wrapped.Diff(ctx, state, terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff.RequiresNew() {
				t.Fatal("description update must not replace the endpoint")
			}
			if testCase.name != "unchanged" {
				if change := diff.Attributes["description"]; change == nil || change.Old != testCase.old || change.New != testCase.expected {
					t.Fatalf("expected description diff %q -> %q, got %#v", testCase.old, testCase.expected, change)
				}
			}

			client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			remote := endpoints.Endpoint{
				Properties: endpoints.AzureStorageNfsFileShareEndpointProperties{
					FileShareName:            "share",
					StorageAccountResourceId: accountID,
					Description:              &testCase.server,
				},
			}
			reads, writes := 0, 0
			client.Client.SetTransport(nfsTargetEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path != id.ID() {
					return nil, fmt.Errorf("unexpected endpoint path: %s", request.URL.Path)
				}
				switch request.Method {
				case http.MethodGet:
					reads++
				case http.MethodPut:
					writes++
					if err := json.NewDecoder(request.Body).Decode(&remote); err != nil {
						return nil, err
					}
					properties, ok := remote.Properties.(endpoints.AzureStorageNfsFileShareEndpointProperties)
					if !ok {
						t.Fatalf("unexpected PUT discriminator: %T", remote.Properties)
					}
					if properties.Description == nil || *properties.Description != testCase.expected {
						t.Fatalf("expected explicit description %q in PUT, got %v", testCase.expected, properties.Description)
					}
					if properties.FileShareName != "share" || properties.StorageAccountResourceId != accountID {
						t.Fatal("description update changed endpoint storage properties")
					}
				default:
					return nil, fmt.Errorf("unexpected method: %s", request.Method)
				}
				body, err := json.Marshal(remote)
				if err != nil {
					return nil, err
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(string(body))),
					Request:    request,
				}, nil
			}))
			providerClient := &clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}}
			if testCase.name == "unchanged" {
				data := wrapped.Data(state)
				if diags := wrapped.UpdateContext(ctx, data, providerClient); diags.HasError() {
					t.Fatal(diags)
				}
				state = data.State()
			} else {
				updated, diags := wrapped.Apply(ctx, state, diff, providerClient)
				if diags.HasError() {
					t.Fatalf("applying description update: %v", diags)
				}
				state = updated
			}
			if state == nil || state.ID != id.ID() || state.Attributes["description"] != testCase.expected {
				t.Fatalf("unexpected state after update/read: %#v", state)
			}
			if testCase.name != "unchanged" {
				nextDiff, err := wrapped.Diff(ctx, state, terraform.NewResourceConfigRaw(config), nil)
				if err != nil {
					t.Fatal(err)
				}
				if nextDiff != nil {
					if !maps.Equal(nextDiff.Identity, state.Identity) {
						t.Fatalf("unexpected identity after update/read: got %v, expected %v", nextDiff.Identity, state.Identity)
					}
					// SDK Diff carries unchanged identity, which Empty treats as nonempty.
					changes := &terraform.InstanceDiff{
						Attributes:     nextDiff.Attributes,
						Destroy:        nextDiff.Destroy,
						DestroyDeposed: nextDiff.DestroyDeposed,
						DestroyTainted: nextDiff.DestroyTainted,
					}
					if !changes.Empty() {
						t.Fatalf("unexpected diff after update/read: %#v", nextDiff)
					}
				}
			}
			if reads != 2 || writes != 1 {
				t.Fatalf("expected update GET/PUT and subsequent Read, got %d GETs and %d PUTs", reads, writes)
			}
		})
	}
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
					Header: http.Header{
						"Content-Type":                   []string{"application/json"},
						sdkclient.SkipPollingDelayHeader: []string{"true"},
					},
					Body:    io.NopCloser(strings.NewReader(responseBody)),
					Request: request,
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
