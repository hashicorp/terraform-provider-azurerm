// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	networksdk "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
)

func TestIpGroupListResource_IDCasing(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	id := ipgroups.NewIPGroupID("12345678-1234-9876-4563-123456789012", "Example-RG", "Example-IPGroup")

	for scope, resourceGroup := range map[string]string{"subscription": "", "resource group": id.ResourceGroupName} {
		for _, testCase := range []struct {
			name string
			id   string
		}{
			{name: "canonical", id: id.ID()},
			{name: "lowercase resourcegroups", id: strings.Replace(id.ID(), "/resourceGroups/", "/resourcegroups/", 1)},
			{name: "mixed static segments", id: "/SUBSCRIPTIONS/" + id.SubscriptionId + "/RESOURCEGROUPS/" + id.ResourceGroupName + "/PROVIDERS/microsoft.network/IPGROUPS/" + id.IpGroupName},
		} {
			t.Run(scope+"/"+testCase.name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					expectedPath := "/subscriptions/" + id.SubscriptionId
					if resourceGroup != "" {
						expectedPath += "/resourceGroups/" + resourceGroup
					}
					expectedPath += "/providers/Microsoft.Network/ipGroups"
					if req.Method != http.MethodGet || req.URL.Path != expectedPath {
						t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
					}
					w.Header().Set("Content-Type", "application/json")
					if err := json.NewEncoder(w).Encode(map[string]any{
						"value": []map[string]any{{
							"id": testCase.id, "name": id.IpGroupName, "location": "centralus",
							"properties": map[string]any{"ipAddresses": []string{"10.0.0.0/24"}},
						}},
					}); err != nil {
						t.Errorf("encoding response: %v", err)
					}
				}))
				defer server.Close()

				client, err := ipgroups.NewIPGroupsClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
				if err != nil {
					t.Fatal(err)
				}
				client.Client.AuthorizeRequest = nil
				client.Client.DisableRetries = true

				wrapper := sdk.FrameworkListResourceWrapper{FrameworkListWrappedResource: IpGroupListResource{}}
				var configSchema list.ListResourceSchemaResponse
				wrapper.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &configSchema)
				var groupValue any
				if resourceGroup != "" {
					groupValue = resourceGroup
				}
				request := list.ListRequest{
					Config: tfsdk.Config{
						Schema: configSchema.Schema,
						Raw: tftypes.NewValue(configSchema.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
							"resource_group_name": tftypes.NewValue(tftypes.String, groupValue),
							"subscription_id":     tftypes.NewValue(tftypes.String, nil),
						}),
					},
					IncludeResource: true,
					ResourceIdentitySchema: schema.Schema{Attributes: map[string]schema.Attribute{
						"name":                schema.StringAttribute{Computed: true},
						"resource_group_name": schema.StringAttribute{Computed: true},
						"subscription_id":     schema.StringAttribute{Computed: true},
					}},
					ResourceSchema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"id":                  schema.StringAttribute{Computed: true},
							"name":                schema.StringAttribute{Computed: true},
							"resource_group_name": schema.StringAttribute{Computed: true},
							"location":            schema.StringAttribute{Computed: true},
							"cidrs":               schema.SetAttribute{Computed: true, ElementType: types.StringType},
							"firewall_ids":        schema.ListAttribute{Computed: true, ElementType: types.StringType},
							"firewall_policy_ids": schema.ListAttribute{Computed: true, ElementType: types.StringType},
							"tags":                schema.MapAttribute{Computed: true, ElementType: types.StringType},
						},
						Blocks: map[string]schema.Block{
							"timeouts": schema.SingleNestedBlock{Attributes: map[string]schema.Attribute{
								"create": schema.StringAttribute{Optional: true},
								"read":   schema.StringAttribute{Optional: true},
								"update": schema.StringAttribute{Optional: true},
								"delete": schema.StringAttribute{Optional: true},
							}},
						},
					},
				}
				var stream list.ListResultsStream
				IpGroupListResource{}.List(ctx, request, &stream, sdk.ResourceMetadata{
					SubscriptionId: id.SubscriptionId,
					Client:         &clients.Client{Network: &networkclient.Client{Client: &networksdk.Client{IPGroups: client}}},
				})
				count := 0
				for result := range stream.Results {
					count++
					if result.Diagnostics.HasError() {
						t.Fatalf("listing IP Groups: %v", result.Diagnostics)
					}
					if result.DisplayName != id.IpGroupName {
						t.Errorf("unexpected display name %q", result.DisplayName)
					}
					var actualID string
					if diags := result.Resource.GetAttribute(ctx, path.Root("id"), &actualID); diags.HasError() {
						t.Fatal(diags)
					}
					if actualID != id.ID() {
						t.Errorf("expected normalized ID %q, got %q", id.ID(), actualID)
					}
					for attribute, expected := range map[string]string{
						"name": id.IpGroupName, "resource_group_name": id.ResourceGroupName, "subscription_id": id.SubscriptionId,
					} {
						var actual string
						if diags := result.Identity.GetAttribute(ctx, path.Root(attribute), &actual); diags.HasError() {
							t.Fatal(diags)
						}
						if actual != expected {
							t.Errorf("identity %s: expected %q, got %q", attribute, expected, actual)
						}
					}
				}
				if count != 1 {
					t.Errorf("expected one result, got %d", count)
				}
			})
		}
	}
}
