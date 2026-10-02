// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestKubernetesFleetManagerHubProfileValidation(t *testing.T) {
	wrapper := sdk.NewResourceWrapper(KubernetesFleetManagerResource{})
	resource, err := wrapper.Resource()
	if err != nil {
		t.Fatalf("building resource: %v", err)
	}
	if err := resource.InternalValidate(nil, true); err != nil {
		t.Fatalf("validating resource schema: %v", err)
	}

	const subnetID = "/subscriptions/00000000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/subnet"
	testCases := []struct {
		name          string
		hub           map[string]any
		errorContains string
	}{
		{name: "without hub"},
		{name: "default hub", hub: map[string]any{}},
		{
			name:          "empty agent profile",
			hub:           map[string]any{"agent_profile": []any{map[string]any{}}},
			errorContains: "hub_profile.0.agent_profile.0.subnet_id",
		},
		{
			name: "agent subnet only",
			hub:  map[string]any{"agent_profile": []any{map[string]any{"subnet_id": subnetID}}},
		},
		{
			name: "agent VM size only",
			hub:  map[string]any{"agent_profile": []any{map[string]any{"virtual_machine_size": "Standard_DS2_v2"}}},
		},
		{
			name: "both agent fields",
			hub:  map[string]any{"agent_profile": []any{map[string]any{"subnet_id": subnetID, "virtual_machine_size": "Standard_D2as_v7"}}},
		},
		{
			name:          "empty API server access profile",
			hub:           map[string]any{"api_server_access_profile": []any{map[string]any{}}},
			errorContains: "private_cluster_enabled",
		},
		{
			name: "explicit public cluster",
			hub:  map[string]any{"api_server_access_profile": []any{map[string]any{"private_cluster_enabled": false}}},
		},
		{
			name: "explicit private cluster",
			hub: map[string]any{
				"agent_profile":             []any{map[string]any{"subnet_id": subnetID}},
				"api_server_access_profile": []any{map[string]any{"private_cluster_enabled": true}},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := map[string]any{
				"name":                "test",
				"resource_group_name": "rg",
				"location":            "eastus",
			}
			if testCase.hub != nil {
				config["hub_profile"] = []any{testCase.hub}
			}
			diagnostics := resource.Validate(terraform.NewResourceConfigRaw(config))
			if testCase.errorContains == "" {
				if diagnostics.HasError() {
					t.Fatalf("unexpected configuration errors: %v", diagnostics)
				}
				return
			}
			if !diagnostics.HasError() {
				t.Fatal("expected empty profile to be rejected")
			}
			for _, diagnostic := range diagnostics {
				if strings.Contains(diagnostic.Detail, testCase.errorContains) {
					return
				}
			}
			t.Fatalf("expected error containing %q, got %v", testCase.errorContains, diagnostics)
		})
	}

	for _, testCase := range []struct {
		profile   string
		attribute string
		valueType cty.Type
	}{
		{profile: "agent_profile", attribute: "subnet_id", valueType: cty.String},
		{profile: "agent_profile", attribute: "virtual_machine_size", valueType: cty.String},
		{profile: "api_server_access_profile", attribute: "private_cluster_enabled", valueType: cty.Bool},
	} {
		t.Run("unknown "+testCase.attribute, func(t *testing.T) {
			config := cty.ObjectVal(map[string]cty.Value{
				"name":                cty.StringVal("test"),
				"resource_group_name": cty.StringVal("rg"),
				"location":            cty.StringVal("eastus"),
				"hub_profile": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
					testCase.profile: cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
						testCase.attribute: cty.UnknownVal(testCase.valueType),
					})}),
				})}),
			})
			schema := resource.CoreConfigSchema()
			value, err := schema.CoerceValue(config)
			if err != nil {
				t.Fatalf("coercing configuration: %v", err)
			}
			if diagnostics := resource.Validate(terraform.NewResourceConfigShimmed(value, schema)); diagnostics.HasError() {
				t.Fatalf("unexpected errors for an unknown value: %v", diagnostics)
			}
		})
	}
}

func TestKubernetesFleetManagerHubProfileDiff(t *testing.T) {
	const subnetID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/subnet"
	testCases := []struct {
		name           string
		privateCluster bool
		hub            map[string]any
		updateTags     bool
		replacementKey string
	}{
		{name: "omit public profiles", hub: map[string]any{}},
		{name: "omit private profiles", privateCluster: true, hub: map[string]any{}},
		{name: "update tags with omitted private profiles", privateCluster: true, hub: map[string]any{}, updateTags: true},
		{
			name: "unchanged explicit public cluster",
			hub: map[string]any{
				"api_server_access_profile": []any{map[string]any{"private_cluster_enabled": false}},
			},
		},
		{
			name:           "unchanged explicit private cluster",
			privateCluster: true,
			hub: map[string]any{
				"api_server_access_profile": []any{map[string]any{"private_cluster_enabled": true}},
			},
		},
		{
			name:           "explicitly disable private cluster",
			privateCluster: true,
			hub: map[string]any{
				"api_server_access_profile": []any{map[string]any{"private_cluster_enabled": false}},
			},
			replacementKey: "hub_profile.0.api_server_access_profile.0.private_cluster_enabled",
		},
		{
			name: "change agent subnet",
			hub: map[string]any{
				"agent_profile": []any{map[string]any{"subnet_id": subnetID + "-other"}},
			},
			replacementKey: "hub_profile.0.agent_profile.0.subnet_id",
		},
		{
			name: "change agent VM size",
			hub: map[string]any{
				"agent_profile": []any{map[string]any{"virtual_machine_size": "Standard_D4s_v3"}},
			},
			replacementKey: "hub_profile.0.agent_profile.0.virtual_machine_size",
		},
		{name: "remove entire hub", replacementKey: "hub_profile.#"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wrapper := sdk.NewResourceWrapper(KubernetesFleetManagerResource{})
			resource, err := wrapper.Resource()
			if err != nil {
				t.Fatalf("building resource: %v", err)
			}

			state := &terraform.InstanceState{
				ID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.ContainerService/fleets/test",
				Attributes: map[string]string{
					"name":                          "test",
					"resource_group_name":           "rg",
					"location":                      "eastus",
					"tags.%":                        "1",
					"tags.phase":                    "initial",
					"hub_profile.#":                 "1",
					"hub_profile.0.agent_profile.#": "1",
					"hub_profile.0.agent_profile.0.subnet_id":                           subnetID,
					"hub_profile.0.agent_profile.0.virtual_machine_size":                "Standard_DS2_v2",
					"hub_profile.0.api_server_access_profile.#":                         "1",
					"hub_profile.0.api_server_access_profile.0.private_cluster_enabled": strconv.FormatBool(testCase.privateCluster),
					"hub_profile.0.dns_prefix":                                          "",
					"hub_profile.0.fqdn":                                                "fleet.example",
					"hub_profile.0.kubernetes_version":                                  "1.32.0",
					"hub_profile.0.portal_fqdn":                                         "portal.example",
				},
			}
			if !testCase.privateCluster {
				state.Attributes["hub_profile.0.dns_prefix"] = "fleet-test"
			}
			phase := "initial"
			if testCase.updateTags {
				phase = "updated"
			}
			config := map[string]any{
				"name":                "test",
				"resource_group_name": "rg",
				"location":            "eastus",
				"tags":                map[string]any{"phase": phase},
			}
			if testCase.hub != nil {
				config["hub_profile"] = []any{testCase.hub}
			}
			resourceConfig := terraform.NewResourceConfigRaw(config)
			if diagnostics := resource.Validate(resourceConfig); diagnostics.HasError() {
				t.Fatalf("validating configuration: %v", diagnostics)
			}
			diff, err := resource.Diff(context.Background(), state, resourceConfig, &clients.Client{})
			if err != nil {
				t.Fatalf("building update diff: %v", err)
			}
			if diff.RequiresNew() != (testCase.replacementKey != "") {
				t.Fatalf("unexpected replacement decision: %#v", diff)
			}
			if testCase.replacementKey != "" {
				change := diff.Attributes[testCase.replacementKey]
				if change == nil || !change.RequiresNew {
					t.Fatalf("expected replacement at %q, got %#v", testCase.replacementKey, diff.Attributes)
				}
				return
			}
			if diff != nil {
				for key, change := range diff.Attributes {
					if strings.HasPrefix(key, "hub_profile.") {
						t.Fatalf("omitted hub settings changed at %q: %#v", key, change)
					}
				}
			}
			if testCase.updateTags {
				if diff == nil {
					t.Fatal("expected a tag update")
				}
				change := diff.Attributes["tags.phase"]
				if change == nil || change.Old != "initial" || change.New != "updated" {
					t.Fatalf("expected tag update, got %#v", diff.Attributes)
				}
			}
		})
	}
}

func TestKubernetesFleetManagerHubProfileMapping(t *testing.T) {
	r := KubernetesFleetManagerResource{}

	for _, testCase := range []struct {
		name    string
		profile *fleets.APIServerAccessProfile
	}{
		{name: "absent API server profile"},
		{name: "absent private cluster setting", profile: &fleets.APIServerAccessProfile{}},
		{name: "explicit public cluster", profile: &fleets.APIServerAccessProfile{EnablePrivateCluster: pointer.To(false)}},
		{name: "explicit private cluster", profile: &fleets.APIServerAccessProfile{EnablePrivateCluster: pointer.To(true)}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			state := flattenFleetManagerHubAPIServerAccessProfile(testCase.profile)
			request := expandFleetManagerHubAPIServerAccessProfile(state)
			if testCase.profile == nil || testCase.profile.EnablePrivateCluster == nil {
				if len(state) != 0 || request != nil {
					t.Fatal("expected absent private cluster mode to remain unspecified")
				}
			} else if !reflect.DeepEqual(request, testCase.profile) {
				t.Fatalf("expected explicit private cluster mode to round-trip, got %#v", request)
			}
		})
	}

	t.Run("create without hub", func(t *testing.T) {
		var payload fleets.Fleet
		r.mapKubernetesFleetManagerResourceSchemaToFleet(KubernetesFleetManagerResourceSchema{}, &payload)
		if payload.Properties == nil || payload.Properties.HubProfile != nil {
			t.Fatalf("expected initialized properties without a hub, got %#v", payload.Properties)
		}
	})

	t.Run("create with hub", func(t *testing.T) {
		config := KubernetesFleetManagerResourceSchema{
			HubProfile: []FleetManagerHubProfile{{
				AgentProfile: []FleetManagerHubAgentProfile{{
					SubnetId:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/subnet",
					VirtualMachineSize: "Standard_DS2_v2",
				}},
				ApiServerAccessProfile: []FleetManagerHubAPIServerAccessProfile{{
					PrivateClusterEnabled: true,
				}},
				DnsPrefix:         "fleet-test",
				Fqdn:              "read-only.example",
				KubernetesVersion: "1.32.0",
				PortalFqdn:        "read-only-portal.example",
			}},
		}
		var payload fleets.Fleet
		r.mapKubernetesFleetManagerResourceSchemaToFleet(config, &payload)
		expected := &fleets.FleetHubProfile{
			AgentProfile: &fleets.AgentProfile{
				SubnetId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/subnet"),
				VMSize:   pointer.To("Standard_DS2_v2"),
			},
			ApiServerAccessProfile: &fleets.APIServerAccessProfile{
				EnablePrivateCluster: pointer.To(true),
			},
			DnsPrefix: pointer.To("fleet-test"),
		}
		if !reflect.DeepEqual(payload.Properties.HubProfile, expected) {
			t.Fatalf("expected configured hub profile in the request, got %#v", payload.Properties.HubProfile)
		}
	})

	t.Run("update without hub preserves existing hub", func(t *testing.T) {
		existing := &fleets.FleetHubProfile{
			DnsPrefix: pointer.To("fleet-test"),
			Fqdn:      pointer.To("fleet.example"),
		}
		payload := fleets.Fleet{Properties: &fleets.FleetProperties{HubProfile: existing}}
		r.mapKubernetesFleetManagerResourceSchemaToFleet(KubernetesFleetManagerResourceSchema{}, &payload)
		if payload.Properties.HubProfile != existing {
			t.Fatal("expected an omitted hub profile to preserve the existing API payload")
		}
	})

	t.Run("read and import include every hub attribute", func(t *testing.T) {
		payload := fleets.Fleet{Properties: &fleets.FleetProperties{
			HubProfile: &fleets.FleetHubProfile{
				AgentProfile: &fleets.AgentProfile{
					SubnetId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg/providers/microsoft.network/virtualnetworks/vnet/subnets/subnet"),
					VMSize:   pointer.To("Standard_DS2_v2"),
				},
				ApiServerAccessProfile: &fleets.APIServerAccessProfile{
					EnablePrivateCluster: pointer.To(true),
				},
				DnsPrefix:         pointer.To("fleet-test"),
				Fqdn:              pointer.To("fleet.example"),
				KubernetesVersion: pointer.To("1.32.0"),
				PortalFqdn:        pointer.To("portal.example"),
			},
		}}
		var state KubernetesFleetManagerResourceSchema
		if err := r.mapFleetToKubernetesFleetManagerResourceSchema(payload, &state); err != nil {
			t.Fatal(err)
		}
		expected := []FleetManagerHubProfile{{
			AgentProfile: []FleetManagerHubAgentProfile{{
				SubnetId:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/subnet",
				VirtualMachineSize: "Standard_DS2_v2",
			}},
			ApiServerAccessProfile: []FleetManagerHubAPIServerAccessProfile{{
				PrivateClusterEnabled: true,
			}},
			DnsPrefix:         "fleet-test",
			Fqdn:              "fleet.example",
			KubernetesVersion: "1.32.0",
			PortalFqdn:        "portal.example",
		}}
		if !reflect.DeepEqual(state.HubProfile, expected) {
			t.Fatalf("expected hub profile %#v, got %#v", expected, state.HubProfile)
		}
	})

	for name, properties := range map[string]*fleets.FleetProperties{
		"missing properties":  nil,
		"missing hub profile": {},
	} {
		t.Run(name, func(t *testing.T) {
			state := KubernetesFleetManagerResourceSchema{
				HubProfile: []FleetManagerHubProfile{{DnsPrefix: "stale"}},
			}
			if err := r.mapFleetToKubernetesFleetManagerResourceSchema(fleets.Fleet{Properties: properties}, &state); err != nil {
				t.Fatal(err)
			}
			if len(state.HubProfile) != 0 {
				t.Fatalf("expected missing API hub profile to clear state, got %#v", state.HubProfile)
			}
		})
	}

	t.Run("invalid subnet ID is rejected", func(t *testing.T) {
		payload := fleets.Fleet{Properties: &fleets.FleetProperties{
			HubProfile: &fleets.FleetHubProfile{
				AgentProfile: &fleets.AgentProfile{SubnetId: pointer.To("invalid")},
			},
		}}
		var state KubernetesFleetManagerResourceSchema
		if err := r.mapFleetToKubernetesFleetManagerResourceSchema(payload, &state); err == nil {
			t.Fatal("expected invalid subnet ID to be rejected")
		}
	})
}
