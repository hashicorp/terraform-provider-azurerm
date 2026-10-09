// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-10-01/agents"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestSreAgentListPopulatedResults(t *testing.T) {
	for _, includeResource := range []bool{false, true} {
		for _, stopEarly := range []bool{false, true} {
			t.Run(fmt.Sprintf("include_resource=%t/stop_early=%t", includeResource, stopEarly), func(t *testing.T) {
				calls := 0
				subscriptionID := "00000000-0000-0000-0000-000000000000"
				path := "/subscriptions/" + subscriptionID + "/providers/Microsoft.App/agents"
				md, _, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
					calls++
					if req.Method != http.MethodGet || req.URL.Path != path || req.URL.Query().Get("api-version") != "2026-10-01" {
						return nil, fmt.Errorf("unexpected list request %s %s", req.Method, req.URL)
					}
					if calls > 2 || (calls == 2 && req.URL.Query().Get("$skiptoken") != "next") {
						return nil, fmt.Errorf("unexpected page %d: %s", calls, req.URL)
					}
					id := agents.NewAgentID(subscriptionID, "example", fmt.Sprintf("agent-%d", calls))
					body := fmt.Sprintf(`{"value":[{"id":%q,"name":%q,"location":"East US","identity":{"type":"SystemAssigned"},"tags":{"stage":"listed"},"properties":{"agentEndpoint":"https://example.invalid","actionConfiguration":{"identity":%q,"mode":"ReadOnly","accessLevel":"Low"},"knowledgeGraphConfiguration":{"identity":%q,"managedResources":[]},"defaultModel":{"name":"server-selected","provider":"MicrosoftFoundry"}}}]`,
						id.ID(), id.AgentName, testIdentityID, testIdentityID+"-resources")
					if calls == 1 {
						body += fmt.Sprintf(`,"nextLink":"https://management.azure.com%s?api-version=2026-10-01&$skiptoken=next"`, path)
					}
					return sreAgentHTTPResponse(req, http.StatusOK, body+"}"), nil
				})
				server, request := sreAgentListTestServer(t, ctx, md)
				request.IncludeResource = includeResource
				stream, err := server.ListResource(ctx, request)
				if err != nil {
					t.Fatal(err)
				}
				wrapped := sdk.WrappedResource(SreAgentResource{})
				results := 0
				for result := range stream.Results {
					results++
					if len(result.Diagnostics) != 0 || result.Identity == nil || result.Identity.IdentityData == nil {
						t.Fatalf("incomplete list result: %#v", result)
					}
					name := fmt.Sprintf("agent-%d", results)
					if result.DisplayName != name {
						t.Fatalf("display name = %q, want %q", result.DisplayName, name)
					}
					identity, err := result.Identity.IdentityData.Unmarshal(wrapped.ProtoIdentitySchema(ctx)().ValueType())
					if err != nil {
						t.Fatal(err)
					}
					for key, want := range map[string]string{"subscription_id": subscriptionID, "resource_group_name": "example", "name": name} {
						assertSreAgentListValue(t, identity, tftypes.NewAttributePath().WithAttributeName(key), tftypes.NewValue(tftypes.String, want))
					}
					if includeResource {
						if result.Resource == nil {
							t.Fatal("include_resource must return resource state")
						}
						state, err := result.Resource.Unmarshal(wrapped.ProtoSchema(ctx)().ValueType())
						if err != nil {
							t.Fatal(err)
						}
						for key, want := range map[string]string{
							"id": agents.NewAgentID(subscriptionID, "example", name).ID(), "name": name,
							"resource_group_name": "example", "location": "eastus", "endpoint": "https://example.invalid",
						} {
							assertSreAgentListValue(t, state, tftypes.NewAttributePath().WithAttributeName(key), tftypes.NewValue(tftypes.String, want))
						}
						for _, field := range []struct{ block, key, want string }{
							{"identity", "type", "SystemAssigned"},
							{"action_configuration", "identity_id", testIdentityID},
							{"action_configuration", "mode", "ReadOnly"},
							{"resources_configuration", "identity_id", testIdentityID + "-resources"},
							{"default_model", "name", "server-selected"},
						} {
							path := tftypes.NewAttributePath().WithAttributeName(field.block).WithElementKeyInt(0).WithAttributeName(field.key)
							assertSreAgentListValue(t, state, path, tftypes.NewValue(tftypes.String, field.want))
						}
						assertSreAgentListValue(t, state, tftypes.NewAttributePath().WithAttributeName("tags").WithElementKeyString("stage"), tftypes.NewValue(tftypes.String, "listed"))
						scope := tftypes.NewAttributePath().WithAttributeName("resources_configuration").WithElementKeyInt(0).WithAttributeName("resource_ids")
						assertSreAgentListValue(t, state, scope, tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{}))
					}
					if stopEarly {
						break
					}
				}
				wantResults := 2
				if stopEarly {
					wantResults = 1
				}
				if results != wantResults || calls != 2 {
					t.Fatalf("got %d results and %d requests; want %d results and two list pages without per-agent GETs", results, calls, wantResults)
				}
			})
		}
	}
}

func TestSreAgentListInvalidID(t *testing.T) {
	for _, idField := range []string{`"id":"/invalid",`, `"id":"",`, ""} {
		t.Run(idField, func(t *testing.T) {
			calls := 0
			md, _, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
				calls++
				return sreAgentHTTPResponse(req, http.StatusOK, `{"value":[{`+idField+`"name":"invalid","location":"eastus"},{"name":"must-not-be-emitted"}]}`), nil
			})
			server, request := sreAgentListTestServer(t, ctx, md)
			stream, err := server.ListResource(ctx, request)
			if err != nil {
				t.Fatal(err)
			}
			results := 0
			for result := range stream.Results {
				results++
				if len(result.Diagnostics) != 1 || result.Diagnostics[0].Severity != tfprotov5.DiagnosticSeverityError ||
					!strings.Contains(result.Diagnostics[0].Summary, "parsing SRE Agent ID") {
					t.Fatalf("expected a surfaced ID parsing error, got %#v", result.Diagnostics)
				}
			}
			if calls != 1 || results != 1 {
				t.Fatalf("invalid ID must stop result emission: calls=%d results=%d", calls, results)
			}
		})
	}
}

// Use the framework's protocol server to convert the real SDK schemas, rather than maintaining test copies.
type sreAgentListTestProvider struct {
	metadata sdk.ResourceMetadata
}

func (sreAgentListTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "azurerm"
}

func (sreAgentListTestProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {
}

func (sreAgentListTestProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (sreAgentListTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func (sreAgentListTestProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p sreAgentListTestProvider) ListResources(context.Context) []func() list.ListResource {
	return []func() list.ListResource{func() list.ListResource {
		return &sdk.FrameworkListResourceWrapper{ResourceMetadata: p.metadata, FrameworkListWrappedResource: SreAgentListResource{}}
	}}
}

func sreAgentListTestServer(t *testing.T, ctx context.Context, md sdk.ResourceMetaData) (tfprotov5.ListResourceServer, *tfprotov5.ListResourceRequest) {
	t.Helper()
	wrapper := sdk.FrameworkListResourceWrapper{FrameworkListWrappedResource: SreAgentListResource{}}
	var schema list.ListResourceSchemaResponse
	wrapper.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &schema)
	configType := schema.Schema.Type().TerraformType(ctx)
	config, err := tfprotov5.NewDynamicValue(configType, tftypes.NewValue(configType, map[string]tftypes.Value{
		"subscription_id": tftypes.NewValue(tftypes.String, nil), "resource_group_name": tftypes.NewValue(tftypes.String, nil),
	}))
	if err != nil {
		t.Fatal(err)
	}
	server, ok := providerserver.NewProtocol5(sreAgentListTestProvider{metadata: sdk.ResourceMetadata{Client: md.Client}})().(tfprotov5.ListResourceServer)
	if !ok {
		t.Fatal("framework protocol server does not implement ListResource")
	}
	return server, &tfprotov5.ListResourceRequest{TypeName: "azurerm_sre_agent", Config: &config, Limit: 100}
}

func assertSreAgentListValue(t *testing.T, state tftypes.Value, path *tftypes.AttributePath, want tftypes.Value) {
	t.Helper()
	value, remaining, err := tftypes.WalkAttributePath(state, path)
	if err != nil {
		t.Fatalf("reading %s (%s remaining): %v", path, remaining, err)
	}
	got, ok := value.(tftypes.Value)
	if !ok || !got.Equal(want) {
		t.Fatalf("%s = %v, want %v", path, value, want)
	}
}

func TestSreAgentListScopesPaginationAndErrors(t *testing.T) {
	for _, scope := range []string{"provider", "subscription", "resource_group"} {
		for _, status := range []int{http.StatusOK, http.StatusForbidden} {
			t.Run(fmt.Sprintf("%s/%d", scope, status), func(t *testing.T) {
				subscriptionID := "00000000-0000-0000-0000-000000000000"
				if scope != "provider" {
					subscriptionID = "11111111-1111-1111-1111-111111111111"
				}
				path := "/subscriptions/" + subscriptionID
				if scope == "resource_group" {
					path += "/resourceGroups/example"
				}
				path += "/providers/Microsoft.App/agents"
				calls := 0
				md, _, ctx := sreAgentTestMetadata(t, func(req *http.Request) (*http.Response, error) {
					calls++
					if req.Method != http.MethodGet || req.URL.Path != path {
						return nil, fmt.Errorf("unexpected list request %s %s", req.Method, req.URL)
					}
					if status != http.StatusOK {
						return sreAgentHTTPResponse(req, status, `{"error":{"code":"AuthorizationFailed","message":"denied"}}`), nil
					}
					body := `{"value":[]}`
					if calls == 1 {
						body = fmt.Sprintf(`{"value":[],"nextLink":"https://management.azure.com%s?api-version=2026-10-01&$skiptoken=next"}`, path)
					} else if calls != 2 || req.URL.Query().Get("$skiptoken") != "next" {
						return nil, fmt.Errorf("unexpected page %d: %s", calls, req.URL)
					}
					return sreAgentHTTPResponse(req, http.StatusOK, body), nil
				})
				wrapper := sdk.FrameworkListResourceWrapper{FrameworkListWrappedResource: SreAgentListResource{}}
				var configSchema list.ListResourceSchemaResponse
				wrapper.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &configSchema)
				values := map[string]tftypes.Value{
					"subscription_id":     tftypes.NewValue(tftypes.String, nil),
					"resource_group_name": tftypes.NewValue(tftypes.String, nil),
				}
				if scope != "provider" {
					values["subscription_id"] = tftypes.NewValue(tftypes.String, subscriptionID)
				}
				if scope == "resource_group" {
					values["resource_group_name"] = tftypes.NewValue(tftypes.String, "example")
				}
				request := list.ListRequest{Config: tfsdk.Config{
					Raw: tftypes.NewValue(configSchema.Schema.Type().TerraformType(ctx), values), Schema: configSchema.Schema,
				}}
				var stream list.ListResultsStream
				SreAgentListResource{}.List(ctx, request, &stream, sdk.ResourceMetadata{Client: md.Client})
				results := 0
				for result := range stream.Results {
					results++
					if status == http.StatusOK || !result.Diagnostics.HasError() {
						t.Fatalf("unexpected list result: %#v", result)
					}
				}
				if status == http.StatusOK && (calls != 2 || results != 0) {
					t.Fatalf("expected two empty pages, got requests=%d results=%d", calls, results)
				}
				if status != http.StatusOK && (calls != 1 || results != 1) {
					t.Fatalf("expected one surfaced list error, got requests=%d results=%d", calls, results)
				}
			})
		}
	}
}
