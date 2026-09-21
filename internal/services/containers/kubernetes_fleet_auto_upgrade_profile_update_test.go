// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

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

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	azclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	containerclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

func TestKubernetesFleetAutoUpgradeProfileUpdate(t *testing.T) {
	id := autoupgradeprofiles.NewAutoUpgradeProfileID("00000000-0000-0000-0000-000000000000", "example", "fleet", "profile")
	fleetID := commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID()
	strategyID := fleetupdatestrategies.NewUpdateStrategyID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "strategy").ID()

	tests := []struct {
		name         string
		beforeImage  string
		afterImage   string
		beforeEnable bool
		afterEnable  bool
	}{
		{name: "remove selection", beforeImage: "Latest"},
		{name: "restore selection", afterImage: "Latest"},
		{name: "change selection", beforeImage: "Latest", afterImage: "Consistent"},
		{name: "reenable preserving selection", beforeImage: "Latest", afterImage: "Latest", afterEnable: true},
		{name: "disable preserving selection", beforeImage: "Latest", afterImage: "Latest", beforeEnable: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := autoupgradeprofiles.NewAutoUpgradeProfilesClientWithBaseURI(environments.NewApiEndpoint("fleet-test", "https://fleet.invalid", nil))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			props := map[string]interface{}{
				"channel":           "Stable",
				"disabled":          !test.beforeEnable,
				"updateStrategyId":  strategyID,
				"provisioningState": "Succeeded",
			}
			if test.beforeImage != "" {
				props["nodeImageSelection"] = map[string]interface{}{"type": test.beforeImage}
			}
			requests := []string{}
			client.Client.SetTransport(fleetAutoUpgradeProfileTransport(func(request *http.Request) (*http.Response, error) {
				requests = append(requests, request.Method)
				if request.URL.Path != id.ID() || request.URL.Query().Get("api-version") != "2025-03-01" {
					t.Fatalf("unexpected request URL: %s", request.URL)
				}
				if request.Method == http.MethodPut {
					var payload struct {
						Properties map[string]interface{} `json:"properties"`
					}
					if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
						t.Fatal(err)
					}
					props = payload.Properties
					if props["updateStrategyId"] != strategyID || props["channel"] != "Stable" || props["disabled"] != !test.afterEnable {
						t.Fatalf("update lost unchanged fields or enabled inversion: %+v", props)
					}
					image, present := props["nodeImageSelection"]
					if test.afterImage == "" {
						if present {
							t.Fatalf("selection removal must omit the SDK omitempty field, got %#v", image)
						}
					} else if !present || image.(map[string]interface{})["type"] != test.afterImage {
						t.Fatalf("unexpected selection: %#v", image)
					}
				} else if request.Method != http.MethodGet {
					t.Fatalf("unexpected method: %s", request.Method)
				}
				body, err := json.Marshal(map[string]interface{}{"id": id.ID(), "properties": props})
				if err != nil {
					t.Fatal(err)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}, azclient.SkipPollingDelayHeader: []string{"true"}},
					Body:       io.NopCloser(strings.NewReader(string(body))),
					Request:    request,
				}, nil
			}))

			resource := containers.KubernetesFleetAutoUpgradeProfileResource{}
			metadata := sdk.NewResourceMetaData(&clients.Client{
				Containers: &containerclient.Client{FleetAutoUpgradeProfilesClient: client},
			}, resource)
			metadata.SetID(id)
			previous := containers.KubernetesFleetAutoUpgradeProfileResourceModel{
				Name: id.AutoUpgradeProfileName, KubernetesFleetManagerId: fleetID, Channel: "Stable",
				UpdateStrategyId: strategyID, NodeImageSelectionType: test.beforeImage, Enabled: test.beforeEnable,
			}
			if err := metadata.Encode(&previous); err != nil {
				t.Fatal(err)
			}
			config := map[string]interface{}{
				"name": id.AutoUpgradeProfileName, "kubernetes_fleet_manager_id": fleetID,
				"channel": "Stable", "update_strategy_id": strategyID,
			}
			if !test.afterEnable {
				config["enabled"] = false
			}
			if test.afterImage != "" {
				config["node_image_selection_type"] = test.afterImage
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			wrapped := sdk.WrappedResource(resource)
			diff, err := wrapped.SimpleDiff(ctx, metadata.ResourceData.State(), terraform.NewResourceConfigRaw(config), metadata.Client)
			if err != nil {
				t.Fatal(err)
			}
			if diff.Empty() || diff.RequiresNew() {
				t.Fatalf("expected an in-place update, got %+v", diff)
			}
			state, diags := wrapped.Apply(ctx, metadata.ResourceData.State(), diff, metadata.Client)
			if diags.HasError() {
				t.Fatalf("applying update: %+v", diags)
			}
			if got := strings.Join(requests, ","); got != "GET,PUT,GET,GET" {
				t.Fatalf("expected GET,PUT,poll GET,read GET, got %s", got)
			}
			for key, want := range map[string]string{
				"node_image_selection_type": test.afterImage, "enabled": fmt.Sprint(test.afterEnable),
				"update_strategy_id": strategyID, "channel": "Stable",
			} {
				if got := state.Attributes[key]; got != want {
					t.Errorf("%s: want %q, got %q", key, want, got)
				}
			}
			next, err := wrapped.SimpleDiff(ctx, state, terraform.NewResourceConfigRaw(config), metadata.Client)
			if err != nil {
				t.Fatal(err)
			}
			if len(next.Attributes) != 0 || next.Destroy || next.RequiresNew() || !maps.Equal(next.Identity, state.Identity) || state.ID != id.ID() {
				t.Fatalf("expected unchanged attributes and identity after read, got %+v", next)
			}
		})
	}
}

func TestKubernetesFleetAutoUpgradeProfileStrategyDiff(t *testing.T) {
	resource := containers.KubernetesFleetAutoUpgradeProfileResource{}
	wrapped := sdk.WrappedResource(resource)
	id := autoupgradeprofiles.NewAutoUpgradeProfileID("00000000-0000-0000-0000-000000000000", "example", "fleet", "profile")
	fleetID := commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID()
	strategyID := fleetupdatestrategies.NewUpdateStrategyID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "strategy").ID()
	otherID := fleetupdatestrategies.NewUpdateStrategyID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "other").ID()
	for _, test := range []struct {
		name   string
		before string
		after  string
	}{
		{name: "remove", before: strategyID},
		{name: "add", after: strategyID},
		{name: "replace", before: strategyID, after: otherID},
	} {
		t.Run(test.name, func(t *testing.T) {
			metadata := sdk.NewResourceMetaData(nil, resource)
			metadata.SetID(id)
			if err := metadata.Encode(&containers.KubernetesFleetAutoUpgradeProfileResourceModel{
				Name: id.AutoUpgradeProfileName, KubernetesFleetManagerId: fleetID, Channel: "Stable",
				UpdateStrategyId: test.before, Enabled: true,
			}); err != nil {
				t.Fatal(err)
			}
			config := map[string]interface{}{
				"name": id.AutoUpgradeProfileName, "kubernetes_fleet_manager_id": fleetID, "channel": "Stable",
			}
			if test.after != "" {
				config["update_strategy_id"] = test.after
			}
			diff, err := wrapped.SimpleDiff(context.Background(), metadata.ResourceData.State(), terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			if !diff.RequiresNew() || diff.Attributes["update_strategy_id"] == nil || !diff.Attributes["update_strategy_id"].RequiresNew {
				t.Fatalf("the create-only strategy reference must require replacement: %+v", diff)
			}
		})
	}
}
