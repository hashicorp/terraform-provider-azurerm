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
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover"
	moverclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/client"
)

func TestStorageMoverSmbMountEndpointList(t *testing.T) {
	const subscriptionID = "00000000-0000-0000-0000-000000000000"
	parentID := endpoints.NewStorageMoverID(subscriptionID, "rg", "mover").ID()
	mixed := []string{"NfsMount", "SmbMount", "AzureStorageBlobContainer", "SmbMount"}

	for _, test := range []struct {
		name          string
		endpointTypes []string
		wantNames     []string
		stop          bool
	}{
		{
			name:          "mixed-variants",
			endpointTypes: mixed,
			wantNames:     []string{"endpoint-1", "endpoint-3"},
		},
		{
			name:          "foreign-only",
			endpointTypes: []string{"NfsMount", "AzureStorageBlobContainer"},
		},
		{
			name:          "stop-after-first",
			endpointTypes: mixed,
			wantNames:     []string{"endpoint-1"},
			stop:          true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			var items []string
			for index, endpointType := range test.endpointTypes {
				name := fmt.Sprintf("endpoint-%d", index)
				id := endpoints.NewEndpointID(subscriptionID, "rg", "mover", name)
				items = append(items, fmt.Sprintf(`{"id":%q,"name":%q,"properties":{"endpointType":%q,"host":%q,"shareName":%q}}`, id.ID(), name, endpointType, name+".example", "share-"+name))
			}
			body := `{"value":[` + strings.Join(items, ",") + `]}`

			client, err := endpoints.NewEndpointsClientWithBaseURI(environments.NewApiEndpoint("test", "https://mover.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			requests := 0
			client.Client.SetTransport(smbMountEndpointOwnershipTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				if request.Method != http.MethodGet || request.URL.Host != "mover.invalid" || request.URL.Path != parentID+"/endpoints" || request.URL.Query().Get("api-version") != "2025-07-01" {
					return nil, fmt.Errorf("unexpected request: %s %s", request.Method, request.URL)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    request,
				}, nil
			}))

			wrapper := sdk.FrameworkListResourceWrapper{
				FrameworkListWrappedResource: storagemover.StorageMoverSmbMountEndpointListResource{},
				ResourceMetadata: sdk.ResourceMetadata{
					Client: &clients.Client{StorageMover: &moverclient.Client{EndpointsClient: client}},
				},
			}
			request := smbMountEndpointListRequest(t, ctx, parentID)
			var stream list.ListResultsStream
			wrapper.List(ctx, request, &stream)
			if stream.Results == nil {
				t.Fatal("expected a List result stream")
			}

			results := 0
			stream.Results(func(result list.ListResult) bool {
				if result.Diagnostics.HasError() {
					t.Fatalf("unexpected List diagnostics: %v", result.Diagnostics)
				}
				if results >= len(test.wantNames) {
					t.Fatalf("unexpected extra result: %q", result.DisplayName)
				}
				name := test.wantNames[results]
				results++
				if result.DisplayName != name {
					t.Errorf("expected display name %q, got %q", name, result.DisplayName)
				}
				if result.Identity == nil || result.Resource == nil {
					t.Fatal("expected both identity and resource state")
				}
				for key, expected := range map[string]string{
					"name":                name,
					"subscription_id":     subscriptionID,
					"resource_group_name": "rg",
					"storage_mover_name":  "mover",
				} {
					var actual string
					if diags := result.Identity.GetAttribute(ctx, path.Root(key), &actual); diags.HasError() {
						t.Fatalf("reading identity %s: %v", key, diags)
					}
					if actual != expected {
						t.Errorf("identity %s: expected %q, got %q", key, expected, actual)
					}
				}
				for key, expected := range map[string]string{
					"id":               endpoints.NewEndpointID(subscriptionID, "rg", "mover", name).ID(),
					"name":             name,
					"storage_mover_id": parentID,
					"host":             name + ".example",
					"share_name":       "share-" + name,
				} {
					var actual string
					if diags := result.Resource.GetAttribute(ctx, path.Root(key), &actual); diags.HasError() {
						t.Fatalf("reading state %s: %v", key, diags)
					}
					if actual != expected {
						t.Errorf("state %s: expected %q, got %q", key, expected, actual)
					}
				}
				return !test.stop
			})
			if results != len(test.wantNames) {
				t.Errorf("expected %d results, got %d", len(test.wantNames), results)
			}
			if requests != 1 {
				t.Errorf("expected one parent-scoped GET and no other requests, got %d requests", requests)
			}
		})
	}
}

// Direct List calls bypass the framework server's conversion of SDKv2 schemas.
func smbMountEndpointListRequest(t *testing.T, ctx context.Context, parentID string) list.ListRequest {
	t.Helper()

	resource := storagemover.StorageMoverSmbMountEndpointListResource{}
	var config list.ListResourceSchemaResponse
	resource.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &config)
	if config.Diagnostics.HasError() {
		t.Fatal(config.Diagnostics)
	}
	legacy := resource.ResourceFunc()
	proto := legacy.ProtoSchema(ctx)()
	protoIdentity := legacy.ProtoIdentitySchema(ctx)()
	schema := resourceschema.Schema{
		Attributes: map[string]resourceschema.Attribute{},
		Blocks:     map[string]resourceschema.Block{},
	}
	for _, attribute := range proto.Block.Attributes {
		if !attribute.Type.Equal(tftypes.String) {
			t.Fatalf("expected a string resource attribute: %s", attribute.Name)
		}
		schema.Attributes[attribute.Name] = resourceschema.StringAttribute{
			Required:  attribute.Required,
			Optional:  attribute.Optional,
			Computed:  attribute.Computed,
			Sensitive: attribute.Sensitive,
		}
	}
	for _, block := range proto.Block.BlockTypes {
		if block.TypeName != "timeouts" || block.Nesting != tfprotov5.SchemaNestedBlockNestingModeSingle {
			t.Fatalf("unexpected resource block: %s", block.TypeName)
		}
		attributes := map[string]resourceschema.Attribute{}
		for _, attribute := range block.Block.Attributes {
			if !attribute.Type.Equal(tftypes.String) {
				t.Fatalf("expected a string timeout attribute: %s", attribute.Name)
			}
			attributes[attribute.Name] = resourceschema.StringAttribute{Optional: true}
		}
		schema.Blocks[block.TypeName] = resourceschema.SingleNestedBlock{Attributes: attributes}
	}
	identity := identityschema.Schema{Attributes: map[string]identityschema.Attribute{}}
	for _, attribute := range protoIdentity.IdentityAttributes {
		if !attribute.Type.Equal(tftypes.String) {
			t.Fatalf("expected a string identity attribute: %s", attribute.Name)
		}
		identity.Attributes[attribute.Name] = identityschema.StringAttribute{
			RequiredForImport: attribute.RequiredForImport,
			OptionalForImport: attribute.OptionalForImport,
		}
	}
	if !schema.Type().TerraformType(ctx).Equal(proto.ValueType()) || !identity.Type().TerraformType(ctx).Equal(protoIdentity.ValueType()) {
		t.Fatal("List request schemas do not match the wrapped resource")
	}

	return list.ListRequest{
		Config: tfsdk.Config{
			Schema: config.Schema,
			Raw: tftypes.NewValue(config.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"storage_mover_id": tftypes.NewValue(tftypes.String, parentID),
			}),
		},
		IncludeResource:        true,
		ResourceSchema:         schema,
		ResourceIdentitySchema: identity,
	}
}
