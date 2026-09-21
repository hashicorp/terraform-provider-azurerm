// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	containerclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

func TestKubernetesFleetAutoUpgradeProfileList(t *testing.T) {
	id := autoupgradeprofiles.NewAutoUpgradeProfileID("00000000-0000-0000-0000-000000000000", "example", "fleet", "first")
	fleetID := commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID()
	secondID := autoupgradeprofiles.NewAutoUpgradeProfileID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "second")
	strategyID := fleetupdatestrategies.NewUpdateStrategyID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "strategy").ID()
	first := fmt.Sprintf(`{"id":%q,"name":"first","properties":{"channel":"Stable","disabled":true,"nodeImageSelection":{"type":"Latest"},"updateStrategyId":%q}}`, id.ID(), strategyID)
	second := fmt.Sprintf(`{"id":%q,"name":"second","properties":{"channel":"Rapid","nodeImageSelection":null,"updateStrategyId":null}}`, secondID.ID())
	tests := []struct {
		name      string
		body      string
		nextPage  string
		status    int
		wantCount int
		wantError bool
	}{
		{name: "multiple pages", body: fmt.Sprintf(`{"value":[%s],"nextLink":"https://fleet.invalid/next?api-version=2025-03-01"}`, first), nextPage: fmt.Sprintf(`{"value":[%s]}`, second), status: http.StatusOK, wantCount: 2},
		{name: "empty fleet", body: `{"value":[]}`, status: http.StatusOK},
		{name: "not found is an error", body: `{"error":{"code":"ResourceNotFound","message":"fleet not found"}}`, status: http.StatusNotFound, wantError: true},
		{name: "missing properties", body: fmt.Sprintf(`{"value":[{"id":%q,"name":"first"}]}`, id.ID()), status: http.StatusOK, wantError: true},
		{name: "null properties", body: fmt.Sprintf(`{"value":[{"id":%q,"name":"first","properties":null}]}`, id.ID()), status: http.StatusOK, wantError: true},
		{name: "invalid identity", body: `{"value":[{"id":"invalid","name":"first","properties":{"channel":"Stable"}}]}`, status: http.StatusOK, wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			client, err := autoupgradeprofiles.NewAutoUpgradeProfilesClientWithBaseURI(environments.NewApiEndpoint("fleet-test", "https://fleet.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			requests := 0
			client.Client.SetTransport(fleetAutoUpgradeProfileTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				body := test.body
				wantPath := fleetID + "/autoUpgradeProfiles"
				if requests == 2 && test.nextPage != "" {
					body = test.nextPage
					wantPath = "/next"
				}
				if request.Method != http.MethodGet || request.URL.Path != wantPath || request.URL.Query().Get("api-version") != "2025-03-01" {
					t.Fatalf("unexpected list request: %s %s", request.Method, request.URL)
				}
				return &http.Response{
					StatusCode: test.status,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    request,
				}, nil
			}))

			resource := containers.KubernetesFleetAutoUpgradeProfileListResource{}
			configSchema := list.ListResourceSchemaResponse{}
			resource.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &configSchema)
			resourceAttributes := map[string]schema.Attribute{"enabled": schema.BoolAttribute{Computed: true}}
			for _, name := range []string{"id", "name", "kubernetes_fleet_manager_id", "channel", "node_image_selection_type", "update_strategy_id"} {
				resourceAttributes[name] = schema.StringAttribute{Computed: true}
			}
			resourceAttributes["timeouts"] = schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"create": schema.StringAttribute{Optional: true},
					"read":   schema.StringAttribute{Optional: true},
					"update": schema.StringAttribute{Optional: true},
					"delete": schema.StringAttribute{Optional: true},
				},
			}
			identityAttributes := map[string]identityschema.Attribute{}
			for _, name := range []string{"name", "fleet_name", "resource_group_name", "subscription_id"} {
				identityAttributes[name] = identityschema.StringAttribute{RequiredForImport: true}
			}
			request := list.ListRequest{
				Config: tfsdk.Config{
					Schema: configSchema.Schema,
					Raw: tftypes.NewValue(configSchema.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
						"kubernetes_fleet_manager_id": tftypes.NewValue(tftypes.String, fleetID),
					}),
				},
				IncludeResource:        true,
				ResourceSchema:         schema.Schema{Attributes: resourceAttributes},
				ResourceIdentitySchema: identityschema.Schema{Attributes: identityAttributes},
			}
			stream := list.ListResultsStream{}
			resource.List(ctx, request, &stream, sdk.ResourceMetadata{
				Client: &clients.Client{Containers: &containerclient.Client{FleetAutoUpgradeProfilesClient: client}},
			})
			count, errors := 0, 0
			for result := range stream.Results {
				if result.Diagnostics.HasError() {
					errors++
					if !test.wantError {
						t.Fatalf("unexpected diagnostics: %s", result.Diagnostics)
					}
					continue
				}
				count++
				name, channel, image, strategy, enabled := "first", "Stable", "Latest", strategyID, false
				if count == 2 {
					name, channel, image, strategy, enabled = "second", "Rapid", "", "", true
				}
				if result.DisplayName != name {
					t.Errorf("display name: want %q, got %q", name, result.DisplayName)
				}
				for key, want := range map[string]string{
					"name": name, "fleet_name": id.FleetName, "resource_group_name": id.ResourceGroupName, "subscription_id": id.SubscriptionId,
				} {
					var got string
					if diags := result.Identity.GetAttribute(ctx, path.Root(key), &got); diags.HasError() || got != want {
						t.Errorf("identity %s: want %q, got %q, %s", key, want, got, diags)
					}
				}
				for key, want := range map[string]string{
					"name": name, "channel": channel, "kubernetes_fleet_manager_id": fleetID, "node_image_selection_type": image, "update_strategy_id": strategy,
				} {
					var got string
					if diags := result.Resource.GetAttribute(ctx, path.Root(key), &got); diags.HasError() || got != want {
						t.Errorf("resource %s: want %q, got %q, %s", key, want, got, diags)
					}
				}
				var gotEnabled bool
				if diags := result.Resource.GetAttribute(ctx, path.Root("enabled"), &gotEnabled); diags.HasError() || gotEnabled != enabled {
					t.Errorf("enabled: want %t, got %t, %s", enabled, gotEnabled, diags)
				}
			}
			if count != test.wantCount || (errors == 1) != test.wantError {
				t.Fatalf("got %d resources and %d errors; want %d resources, error=%t", count, errors, test.wantCount, test.wantError)
			}
			wantRequests := 1
			if test.nextPage != "" {
				wantRequests = 2
			}
			if requests != wantRequests {
				t.Fatalf("want %d requests, got %d", wantRequests, requests)
			}
		})
	}
}
