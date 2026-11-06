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

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/client"
)

func TestStorageMoverNfsFileShareTargetEndpointList(t *testing.T) {
	parent := endpoints.NewStorageMoverID("00000000-0000-0000-0000-000000000000", "rg", "mover")
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account"
	targets := make([]endpoints.Endpoint, 2)
	for i, name := range []string{"first", "second"} {
		targets[i] = endpoints.Endpoint{
			Id:   pointer.To(endpoints.NewEndpointID(parent.SubscriptionId, parent.ResourceGroupName, parent.StorageMoverName, name).ID()),
			Name: pointer.To(name),
			Properties: endpoints.AzureStorageNfsFileShareEndpointProperties{
				FileShareName:            "share",
				StorageAccountResourceId: accountID,
			},
		}
	}
	foreign := endpoints.Endpoint{Properties: endpoints.AzureStorageBlobContainerEndpointProperties{
		StorageAccountResourceId: accountID,
		BlobContainerName:        "container",
	}}
	invalidID := targets[0]
	invalidID.Id = pointer.To("invalid")
	invalidAccount := targets[0]
	invalidAccount.Properties = endpoints.AzureStorageNfsFileShareEndpointProperties{
		FileShareName:            "share",
		StorageAccountResourceId: "invalid",
	}

	for _, testCase := range []struct {
		name      string
		items     []endpoints.Endpoint
		status    int
		wantNames []string
		wantError string
		stopFirst bool
	}{
		{name: "mixed variants", items: []endpoints.Endpoint{foreign, targets[0], {Properties: nil}, targets[1]}, wantNames: []string{"first", "second"}},
		{name: "only other variants", items: []endpoints.Endpoint{foreign}},
		{name: "empty", items: []endpoints.Endpoint{}},
		{name: "consumer stops", items: targets, wantNames: []string{"first"}, stopFirst: true},
		{name: "invalid endpoint ID", items: []endpoints.Endpoint{invalidID, targets[1]}, wantError: "parsing Storage Mover NFS File Share Target Endpoint ID"},
		{name: "invalid storage account ID", items: []endpoints.Endpoint{invalidAccount, targets[1]}, wantError: "encoding"},
		{name: "API error", status: http.StatusForbidden, wantError: "listing"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			endpointsClient, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			endpointsClient.Client.AuthorizeRequest = nil
			endpointsClient.Client.DisableRetries = true
			requests := 0
			endpointsClient.Client.SetTransport(nfsTargetEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				if request.Method != http.MethodGet || request.URL.Path != parent.ID()+"/endpoints" {
					return nil, fmt.Errorf("unexpected request: %s %s", request.Method, request.URL.Path)
				}
				status := testCase.status
				if status == 0 {
					status = http.StatusOK
				}
				body, err := json.Marshal(map[string]any{"value": testCase.items})
				if err != nil {
					return nil, err
				}
				if status != http.StatusOK {
					body = []byte(`{"error":{"code":"AuthorizationFailed","message":"denied"}}`)
				}
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(string(body))),
					Request:    request,
				}, nil
			}))

			resource := storagemover.StorageMoverNfsFileShareTargetEndpointListResource{}
			var configSchema list.ListResourceSchemaResponse
			resource.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &configSchema)
			resourceSchema := schema.Schema{
				Attributes: map[string]schema.Attribute{},
				Blocks: map[string]schema.Block{
					"timeouts": schema.SingleNestedBlock{Attributes: map[string]schema.Attribute{
						"create": schema.StringAttribute{Optional: true},
						"read":   schema.StringAttribute{Optional: true},
						"update": schema.StringAttribute{Optional: true},
						"delete": schema.StringAttribute{Optional: true},
					}},
				},
			}
			for _, name := range []string{"id", "name", "storage_mover_id", "file_share_name", "storage_account_id", "description"} {
				resourceSchema.Attributes[name] = schema.StringAttribute{Computed: true}
			}
			identitySchema := schema.Schema{Attributes: map[string]schema.Attribute{}}
			for _, name := range []string{"subscription_id", "resource_group_name", "storage_mover_name", "name"} {
				identitySchema.Attributes[name] = schema.StringAttribute{Required: true}
			}
			if !resourceSchema.Type().TerraformType(ctx).Equal(resource.ResourceFunc().ProtoSchema(ctx)().ValueType()) {
				t.Fatal("test result schema does not match the registered resource")
			}
			if !identitySchema.Type().TerraformType(ctx).Equal(resource.ResourceFunc().ProtoIdentitySchema(ctx)().ValueType()) {
				t.Fatal("test identity schema does not match the registered resource")
			}
			request := list.ListRequest{
				Config: tfsdk.Config{
					Schema: configSchema.Schema,
					Raw: tftypes.NewValue(configSchema.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
						"storage_mover_id": tftypes.NewValue(tftypes.String, parent.ID()),
					}),
				},
				IncludeResource:        true,
				ResourceSchema:         resourceSchema,
				ResourceIdentitySchema: identitySchema,
			}
			var stream list.ListResultsStream
			resource.List(ctx, request, &stream, sdk.ResourceMetadata{
				Client: &clients.Client{StorageMover: &client.Client{EndpointsClient: endpointsClient}},
			})
			if stream.Results == nil {
				t.Fatal("List did not return a result iterator")
			}
			results, errors := 0, 0
			stream.Results(func(result list.ListResult) bool {
				if result.Diagnostics.HasError() {
					errors++
					if testCase.wantError == "" || !strings.Contains(fmt.Sprint(result.Diagnostics), testCase.wantError) {
						t.Fatalf("unexpected diagnostics: %v", result.Diagnostics)
					}
					return true
				}
				if results >= len(testCase.wantNames) {
					t.Fatalf("unexpected result: %q", result.DisplayName)
				}
				name := testCase.wantNames[results]
				results++
				if result.DisplayName != name {
					t.Fatalf("expected display name %q, got %q", name, result.DisplayName)
				}
				for key, expected := range map[string]string{
					"subscription_id":     parent.SubscriptionId,
					"resource_group_name": parent.ResourceGroupName,
					"storage_mover_name":  parent.StorageMoverName,
					"name":                name,
				} {
					var actual string
					if diags := result.Identity.GetAttribute(ctx, path.Root(key), &actual); diags.HasError() {
						t.Fatal(diags)
					}
					if actual != expected {
						t.Errorf("identity %s: expected %q, got %q", key, expected, actual)
					}
				}
				for key, expected := range map[string]string{
					"id":                 endpoints.NewEndpointID(parent.SubscriptionId, parent.ResourceGroupName, parent.StorageMoverName, name).ID(),
					"name":               name,
					"storage_mover_id":   parent.ID(),
					"file_share_name":    "share",
					"storage_account_id": accountID,
					"description":        "",
				} {
					var actual string
					if diags := result.Resource.GetAttribute(ctx, path.Root(key), &actual); diags.HasError() {
						t.Fatal(diags)
					}
					if actual != expected {
						t.Errorf("resource %s: expected %q, got %q", key, expected, actual)
					}
				}
				return !testCase.stopFirst
			})
			if results != len(testCase.wantNames) || errors > 1 || (errors == 1) != (testCase.wantError != "") {
				t.Fatalf("unexpected result counts: successful=%d errors=%d", results, errors)
			}
			if requests != 1 {
				t.Fatalf("expected one List request, got %d", requests)
			}
		})
	}
}
