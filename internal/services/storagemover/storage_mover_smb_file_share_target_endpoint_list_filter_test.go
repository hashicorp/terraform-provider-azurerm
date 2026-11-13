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
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover"
	moverclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/client"
)

func TestStorageMoverSmbFileShareTargetEndpointListFiltersEndpointTypes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	subscriptionID := "00000000-0000-0000-0000-000000000000"
	storageMoverID := storagemovers.NewStorageMoverID(subscriptionID, "rg", "mover").ID()
	body := fmt.Sprintf(`{"value":[
  {"id":"%[1]s/endpoints/blob","name":"blob","properties":{"endpointType":"AzureStorageBlobContainer","storageAccountResourceId":"/subscriptions/%[2]s/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account","blobContainerName":"container"}},
  {"id":"%[1]s/endpoints/smb-first","name":"smb-first","properties":{"endpointType":"AzureStorageSmbFileShare","storageAccountResourceId":"/subscriptions/%[2]s/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account","fileShareName":"share"}},
  {"id":"%[1]s/endpoints/nfs","name":"nfs","properties":{"endpointType":"NfsMount","host":"192.0.2.1","export":"/share"}},
  {"id":"%[1]s/endpoints/smb-second","name":"smb-second","properties":{"endpointType":"AzureStorageSmbFileShare","storageAccountResourceId":"/subscriptions/%[2]s/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/account","fileShareName":"share"}}
]}`, storageMoverID, subscriptionID)

	client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
	if err != nil {
		t.Fatal(err)
	}
	client.Client.AuthorizeRequest = nil
	client.Client.DisableRetries = true
	requests := 0
	client.Client.SetTransport(smbTargetEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
		requests++
		if request.Method != http.MethodGet || request.URL.Path != storageMoverID+"/endpoints" || request.URL.Query().Get("api-version") != "2025-07-01" {
			return nil, fmt.Errorf("unexpected request: %s %s", request.Method, request.URL)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    request,
		}, nil
	}))

	r := storagemover.StorageMoverSmbFileShareTargetEndpointListResource{}
	configSchema := list.ListResourceSchemaResponse{}
	r.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &configSchema)
	resourceSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                 schema.StringAttribute{Computed: true},
			"name":               schema.StringAttribute{Required: true},
			"storage_mover_id":   schema.StringAttribute{Required: true},
			"storage_account_id": schema.StringAttribute{Required: true},
			"file_share_name":    schema.StringAttribute{Required: true},
			"description":        schema.StringAttribute{Optional: true},
		},
		Blocks: map[string]schema.Block{
			"timeouts": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"create": schema.StringAttribute{Optional: true},
					"read":   schema.StringAttribute{Optional: true},
					"update": schema.StringAttribute{Optional: true},
					"delete": schema.StringAttribute{Optional: true},
				},
			},
		},
	}
	identitySchema := identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"name":                identityschema.StringAttribute{RequiredForImport: true},
			"subscription_id":     identityschema.StringAttribute{RequiredForImport: true},
			"resource_group_name": identityschema.StringAttribute{RequiredForImport: true},
			"storage_mover_name":  identityschema.StringAttribute{RequiredForImport: true},
		},
	}
	if !resourceSchema.Type().TerraformType(ctx).Equal(r.ResourceFunc().ProtoSchema(ctx)().ValueType()) {
		t.Fatal("test resource schema does not match the resource's protocol schema")
	}
	if !identitySchema.Type().TerraformType(ctx).Equal(r.ResourceFunc().ProtoIdentitySchema(ctx)().ValueType()) {
		t.Fatal("test identity schema does not match the resource's protocol schema")
	}

	request := list.ListRequest{
		Config: tfsdk.Config{
			Schema: configSchema.Schema,
			Raw: tftypes.NewValue(configSchema.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"storage_mover_id": tftypes.NewValue(tftypes.String, storageMoverID),
			}),
		},
		IncludeResource:        true,
		ResourceSchema:         resourceSchema,
		ResourceIdentitySchema: identitySchema,
	}
	stream := list.ListResultsStream{}
	r.List(ctx, request, &stream, sdk.ResourceMetadata{
		Client: &clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}},
	})

	remaining := map[string]bool{"smb-first": true, "smb-second": true}
	count := 0
	for result := range stream.Results {
		count++
		if len(result.Diagnostics) != 0 {
			t.Fatalf("unexpected list diagnostics: %v", result.Diagnostics)
		}
		if !remaining[result.DisplayName] {
			t.Fatalf("unexpected or duplicate endpoint: %q", result.DisplayName)
		}
		if result.Identity == nil {
			t.Fatal("missing endpoint identity")
		}
		expected := tftypes.NewValue(identitySchema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"name":                tftypes.NewValue(tftypes.String, result.DisplayName),
			"subscription_id":     tftypes.NewValue(tftypes.String, subscriptionID),
			"resource_group_name": tftypes.NewValue(tftypes.String, "rg"),
			"storage_mover_name":  tftypes.NewValue(tftypes.String, "mover"),
		})
		if !result.Identity.Raw.Equal(expected) {
			t.Errorf("expected identity %s, got %s", expected, result.Identity.Raw)
		}
		delete(remaining, result.DisplayName)
	}
	if count != 2 || len(remaining) != 0 {
		t.Errorf("expected both SMB endpoints only, got %d results; missing: %v", count, remaining)
	}
	if requests != 1 {
		t.Errorf("expected one List GET and no other requests, got %d requests", requests)
	}
}
