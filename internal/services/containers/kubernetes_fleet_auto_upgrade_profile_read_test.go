// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	containerclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
)

type fleetAutoUpgradeProfileTransport func(*http.Request) (*http.Response, error)

func (f fleetAutoUpgradeProfileTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestKubernetesFleetAutoUpgradeProfileRead(t *testing.T) {
	id := autoupgradeprofiles.NewAutoUpgradeProfileID("00000000-0000-0000-0000-000000000000", "ExampleRG", "Fleet", "profile")
	strategyID := fleetupdatestrategies.NewUpdateStrategyID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "Strategy").ID()
	tests := []struct {
		name          string
		body          string
		status        int
		wantError     bool
		errorContains string
		wantStrategy  string
		wantGone      bool
		update        bool
	}{
		{
			name:   "omitted optional properties",
			body:   `{"properties":{"channel":"Stable"}}`,
			status: http.StatusOK,
		},
		{
			name:   "null optional properties",
			body:   `{"properties":{"channel":"Stable","nodeImageSelection":null,"updateStrategyId":null,"disabled":null}}`,
			status: http.StatusOK,
		},
		{
			name: "mixed case API strategy ID",
			body: fmt.Sprintf(`{"properties":{"channel":"Stable","updateStrategyId":%q}}`,
				strings.NewReplacer("resourceGroups", "RESOURCEGROUPS", "Microsoft.ContainerService", "microsoft.containerservice", "updateStrategies", "UPDATESTRATEGIES").Replace(strategyID)),
			status:       http.StatusOK,
			wantStrategy: strategyID,
		},
		{
			name:          "wrong API strategy resource type",
			body:          fmt.Sprintf(`{"properties":{"channel":"Stable","updateStrategyId":%q}}`, id.ID()),
			status:        http.StatusOK,
			wantError:     true,
			errorContains: "parsing",
		},
		{
			name:      "missing properties",
			body:      `{}`,
			status:    http.StatusOK,
			wantError: true,
		},
		{
			name:      "null properties",
			body:      `{"properties":null}`,
			status:    http.StatusOK,
			wantError: true,
		},
		{
			name:      "update missing properties",
			body:      `{}`,
			status:    http.StatusOK,
			wantError: true,
			update:    true,
		},
		{
			name:      "update null properties",
			body:      `{"properties":null}`,
			status:    http.StatusOK,
			wantError: true,
			update:    true,
		},
		{
			name:     "not found",
			body:     `{"error":{"code":"ResourceNotFound","message":"not found"}}`,
			status:   http.StatusNotFound,
			wantGone: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := environments.NewApiEndpoint("fleet-test", "https://fleet.invalid", nil)
			client, err := autoupgradeprofiles.NewAutoUpgradeProfilesClientWithBaseURI(endpoint)
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			requests := 0
			client.Client.SetTransport(fleetAutoUpgradeProfileTransport(func(request *http.Request) (*http.Response, error) {
				requests++
				if request.Method != http.MethodGet {
					t.Fatalf("unexpected request method: %s", request.Method)
				}
				return &http.Response{
					StatusCode: test.status,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(test.body)),
					Request:    request,
				}, nil
			}))

			resource := containers.KubernetesFleetAutoUpgradeProfileResource{}
			metadata := sdk.NewResourceMetaData(&clients.Client{
				Containers: &containerclient.Client{FleetAutoUpgradeProfilesClient: client},
			}, resource)
			metadata.SetID(id)
			previous := containers.KubernetesFleetAutoUpgradeProfileResourceModel{
				Name:                     id.AutoUpgradeProfileName,
				KubernetesFleetManagerId: commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID(),
				Channel:                  "Rapid",
				NodeImageSelectionType:   "Latest",
				UpdateStrategyId:         fleetupdatestrategies.NewUpdateStrategyID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "strategy").ID(),
				Enabled:                  false,
			}
			if err := metadata.Encode(&previous); err != nil {
				t.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			if test.update {
				err = resource.Update().Func(ctx, metadata)
			} else {
				err = resource.Read().Func(ctx, metadata)
			}
			if requests != 1 {
				t.Fatalf("expected one GET request, got %d requests", requests)
			}
			if (err != nil) != test.wantError {
				t.Fatalf("want error = %t, got %v", test.wantError, err)
			}
			wantID := id.ID()
			if test.wantGone {
				wantID = ""
			}
			if got := metadata.ResourceData.Id(); got != wantID {
				t.Fatalf("unexpected state ID: want %q, got %q", wantID, got)
			}
			if test.wantError {
				errorContains := test.errorContains
				if errorContains == "" {
					errorContains = "properties"
				}
				if !strings.Contains(err.Error(), errorContains) {
					t.Fatalf("expected error containing %q, got %v", errorContains, err)
				}
				var state containers.KubernetesFleetAutoUpgradeProfileResourceModel
				if err := metadata.Decode(&state); err != nil {
					t.Fatal(err)
				}
				if state != previous {
					t.Fatalf("invalid response changed existing state: want %+v, got %+v", previous, state)
				}
				return
			}
			if test.wantGone {
				return
			}
			if got := metadata.ResourceData.Get("channel"); got != "Stable" {
				t.Fatalf("unexpected channel: %v", got)
			}
			if got := metadata.ResourceData.Get("node_image_selection_type"); got != "" {
				t.Fatalf("removed image selection was not cleared: %v", got)
			}
			if got := metadata.ResourceData.Get("update_strategy_id"); got != test.wantStrategy {
				t.Fatalf("unexpected update strategy: want %q, got %v", test.wantStrategy, got)
			}
			if got := metadata.ResourceData.Get("enabled"); got != true {
				t.Fatalf("omitted disabled property must leave the profile enabled: %v", got)
			}
		})
	}
}

func TestKubernetesFleetAutoUpgradeProfileValidation(t *testing.T) {
	id := autoupgradeprofiles.NewAutoUpgradeProfileID("00000000-0000-0000-0000-000000000000", "ExampleRG", "Fleet", "profile")
	fleetID := commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID()
	strategyID := fleetupdatestrategies.NewUpdateStrategyID(id.SubscriptionId, id.ResourceGroupName, id.FleetName, "Strategy").ID()
	wrapped := sdk.WrappedResource(containers.KubernetesFleetAutoUpgradeProfileResource{})
	configSchema := wrapped.CoreConfigSchema()

	for _, test := range []struct {
		name      string
		key       string
		value     cty.Value
		wantError bool
	}{
		{name: "omitted optional fields"},
		{name: "null selection", key: "node_image_selection_type", value: cty.NullVal(cty.String)},
		{name: "unknown name", key: "name", value: cty.UnknownVal(cty.String)},
		{name: "unknown channel", key: "channel", value: cty.UnknownVal(cty.String)},
		{name: "unknown enabled", key: "enabled", value: cty.UnknownVal(cty.Bool)},
		{name: "unknown parent", key: "kubernetes_fleet_manager_id", value: cty.UnknownVal(cty.String)},
		{name: "unknown strategy", key: "update_strategy_id", value: cty.UnknownVal(cty.String)},
		{name: "unknown selection", key: "node_image_selection_type", value: cty.UnknownVal(cty.String)},
		{name: "canonical strategy", key: "update_strategy_id", value: cty.StringVal(strategyID)},
		{name: "noncanonical config strategy", key: "update_strategy_id", value: cty.StringVal(strings.ReplaceAll(strategyID, "resourceGroups", "RESOURCEGROUPS")), wantError: true},
		{name: "wrong configured strategy type", key: "update_strategy_id", value: cty.StringVal(id.ID()), wantError: true},
		{name: "foreign parent", key: "kubernetes_fleet_manager_id", value: cty.StringVal(id.ID()), wantError: true},
		{name: "invalid channel", key: "channel", value: cty.StringVal("Unknown"), wantError: true},
		{name: "invalid selection", key: "node_image_selection_type", value: cty.StringVal("Unknown"), wantError: true},
		{name: "shortest name", key: "name", value: cty.StringVal("a")},
		{name: "empty name", key: "name", value: cty.StringVal(""), wantError: true},
		{name: "longest name", key: "name", value: cty.StringVal(strings.Repeat("a", 50))},
		{name: "name too long", key: "name", value: cty.StringVal(strings.Repeat("a", 51)), wantError: true},
		{name: "uppercase name", key: "name", value: cty.StringVal("Profile"), wantError: true},
		{name: "trailing hyphen", key: "name", value: cty.StringVal("profile-"), wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := map[string]cty.Value{
				"name":                        cty.StringVal(id.AutoUpgradeProfileName),
				"kubernetes_fleet_manager_id": cty.StringVal(fleetID),
				"channel":                     cty.StringVal("Stable"),
			}
			if test.key != "" {
				config[test.key] = test.value
			}
			value, err := configSchema.CoerceValue(cty.ObjectVal(config))
			if err != nil {
				t.Fatal(err)
			}
			diags := wrapped.Validate(terraform.NewResourceConfigShimmed(value, configSchema))
			if diags.HasError() != test.wantError {
				t.Fatalf("want validation error = %t, got %v", test.wantError, diags)
			}
		})
	}
}

func TestKubernetesFleetAutoUpgradeProfileImport(t *testing.T) {
	id := autoupgradeprofiles.NewAutoUpgradeProfileID("00000000-0000-0000-0000-000000000000", "ExampleRG", "Fleet", "profile")
	resource := containers.KubernetesFleetAutoUpgradeProfileResource{}
	wrapped := sdk.WrappedResource(resource)
	for _, test := range []struct {
		name      string
		id        string
		wantError bool
	}{
		{name: "canonical", id: id.ID()},
		{name: "noncanonical static casing", id: strings.ReplaceAll(id.ID(), "resourceGroups", "RESOURCEGROUPS"), wantError: true},
		{name: "wrong child type", id: strings.ReplaceAll(id.ID(), "autoUpgradeProfiles", "updateStrategies"), wantError: true},
		{name: "wrong parent type", id: strings.ReplaceAll(id.ID(), "fleets", "managedClusters"), wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			metadata := sdk.NewResourceMetaData(nil, resource)
			metadata.ResourceData.SetId(test.id)
			states, err := wrapped.Importer.StateContext(context.Background(), metadata.ResourceData, nil)
			if (err != nil) != test.wantError {
				t.Fatalf("want import error = %t, got %v", test.wantError, err)
			}
			if !test.wantError && (len(states) != 1 || states[0].Id() != id.ID()) {
				t.Fatalf("import changed the canonical resource ID: %+v", states)
			}
		})
	}
}
