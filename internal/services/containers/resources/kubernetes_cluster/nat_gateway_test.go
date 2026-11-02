// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesClusterNatGatewaySkuDiff(t *testing.T) {
	resource := &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"network_profile": {
				Type: pluginsdk.TypeList, Required: true, MaxItems: 1,
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"outbound_type": {Type: pluginsdk.TypeString, Required: true},
					"nat_gateway_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						// NOTE: O+C because the test reproduces Azure-populated omitted profiles.
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
							"sku": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								// NOTE: O+C because the test distinguishes an explicit SKU from a retained Azure value.
								Computed: true,
							},
						}},
					},
				}},
			},
		},
		CustomizeDiff: customizeKubernetesClusterNatGatewaySku,
	}
	config := func(outbound, sku cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{
			"network_profile": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"outbound_type":       outbound,
				"nat_gateway_profile": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"sku": sku})}),
			})}),
		})
	}
	managed := cty.StringVal("managedNATGateway")
	v2 := cty.StringVal("StandardV2")
	tests := []struct {
		name           string
		outbound       cty.Value
		sku            cty.Value
		existing       bool
		unknownProfile bool
		wantErr        bool
	}{
		{name: "managed SKU", outbound: managed, sku: v2},
		{name: "invalid outbound type", outbound: cty.StringVal("userAssignedNATGateway"), sku: v2, wantErr: true},
		{name: "unknown outbound type", outbound: cty.UnknownVal(cty.String), sku: v2},
		{name: "unknown SKU", outbound: cty.StringVal("loadBalancer"), sku: cty.UnknownVal(cty.String)},
		{name: "unknown profile", outbound: managed, sku: v2, unknownProfile: true},
		{name: "null SKU", outbound: cty.StringVal("loadBalancer"), sku: cty.NullVal(cty.String)},
		{name: "downgrade", outbound: managed, sku: cty.StringVal("Standard"), existing: true, wantErr: true},
		{name: "outbound-only change with explicit SKU", outbound: cty.StringVal("userAssignedNATGateway"), sku: v2, existing: true, wantErr: true},
		{name: "outbound change with computed SKU", outbound: cty.StringVal("userAssignedNATGateway"), sku: cty.NullVal(cty.String), existing: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := &terraform.InstanceState{}
			if test.existing {
				initialConfig := config(managed, v2)
				state.RawConfig = initialConfig
				diff, err := resource.Diff(context.Background(), state, terraform.NewResourceConfigShimmed(initialConfig, resource.CoreConfigSchema()), nil)
				if err != nil {
					t.Fatal(err)
				}
				state.Attributes, err = diff.Apply(nil, resource.CoreConfigSchema())
				if err != nil {
					t.Fatal(err)
				}
				state.ID = "test"
			}
			state.RawConfig = config(test.outbound, test.sku)
			if test.unknownProfile {
				state.RawConfig = cty.ObjectVal(map[string]cty.Value{
					"network_profile": cty.UnknownVal(resource.CoreConfigSchema().ImpliedType().AttributeType("network_profile")),
				})
			}
			_, err := resource.Diff(context.Background(), state, terraform.NewResourceConfigShimmed(state.RawConfig, resource.CoreConfigSchema()), nil)
			if (err != nil) != test.wantErr {
				t.Fatalf("error = %v, want error %t", err, test.wantErr)
			}
		})
	}
}

func TestKubernetesClusterNatGatewaySkuOmittedAfterStandardV2(t *testing.T) {
	networkProfile := RegisterResource().Schema["network_profile"].Elem.(*pluginsdk.Resource)
	natGatewayProfile := networkProfile.Schema["nat_gateway_profile"].Elem.(*pluginsdk.Resource)
	resource := &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
		"sku": natGatewayProfile.Schema["sku"],
	}}
	config := func(sku string) *terraform.ResourceConfig {
		profile := map[string]any{}
		if sku != "" {
			profile["sku"] = sku
		}
		return terraform.NewResourceConfigRaw(profile)
	}

	initial, err := resource.Diff(context.Background(), nil, config("StandardV2"), nil)
	if err != nil {
		t.Fatalf("building initial diff: %v", err)
	}
	attributes, err := initial.Apply(nil, resource.CoreConfigSchema())
	if err != nil {
		t.Fatalf("building initial state: %v", err)
	}
	state := &terraform.InstanceState{ID: "test", Attributes: attributes}

	diff, err := resource.Diff(context.Background(), state, config(""), nil)
	if err != nil {
		t.Fatalf("building omission diff: %v", err)
	}
	if diff != nil && diff.Attributes["sku"] != nil {
		t.Fatalf("omitting the SKU must retain StandardV2: %#v", diff.Attributes["sku"])
	}
}

func TestExpandKubernetesClusterNatGatewaySku(t *testing.T) {
	tests := map[string]struct {
		profile []any
		want    *managedclusters.ManagedClusterNATGatewaySku
	}{
		"omitted profile": {profile: []any{}},
		"omitted sku":     {profile: []any{map[string]any{"idle_timeout_in_minutes": 4, "managed_outbound_ip_count": 1}}},
		"standard":        {profile: []any{map[string]any{"sku": "Standard"}}, want: pointer.To(managedclusters.ManagedClusterNATGatewaySkuStandard)},
		"standard v2":     {profile: []any{map[string]any{"sku": "StandardV2"}}, want: pointer.To(managedclusters.ManagedClusterNATGatewaySkuStandardVTwo)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			profile, err := expandKubernetesClusterNetworkProfile([]any{map[string]any{
				"network_plugin":        "kubenet",
				"network_mode":          "",
				"network_policy":        "",
				"load_balancer_profile": []any{},
				"load_balancer_sku":     "standard",
				"nat_gateway_profile":   test.profile,
				"outbound_type":         "managedNATGateway",
				"ip_versions":           []any{},
				"network_data_plane":    "",
				"network_plugin_mode":   "",
			}}, nil)
			if err != nil {
				t.Fatalf("expanding network profile: %+v", err)
			}
			if profile.NatGatewayProfile == nil {
				if test.want == nil {
					return
				}
				t.Fatalf("expected NAT Gateway SKU %q, got no NAT Gateway profile", *test.want)
			}
			if (profile.NatGatewayProfile.Sku == nil) != (test.want == nil) || (test.want != nil && *profile.NatGatewayProfile.Sku != *test.want) {
				if test.want == nil {
					t.Fatalf("expected no NAT Gateway SKU, got %+v", profile.NatGatewayProfile.Sku)
				}
				t.Fatalf("expected NAT Gateway SKU %q, got %+v", *test.want, profile.NatGatewayProfile.Sku)
			}
		})
	}
}

func TestKubernetesClusterNatGatewaySkuValidation(t *testing.T) {
	schema := RegisterResource().Schema["network_profile"].Elem.(*pluginsdk.Resource).Schema["nat_gateway_profile"].Elem.(*pluginsdk.Resource).Schema["sku"]
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{name: "Standard exact casing", value: "Standard", valid: true},
		{name: "StandardV2 exact casing", value: "StandardV2", valid: true},
		{name: "lowercase", value: "standardv2"},
		{name: "unknown enum", value: "StandardV3"},
		{name: "overlong value", value: strings.Repeat("x", 256)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, errors := schema.ValidateFunc(test.value, "sku")
			if got := len(errors) == 0; got != test.valid {
				t.Fatalf("validation success = %t, want %t (errors: %v)", got, test.valid, errors)
			}
		})
	}
}

func TestKubernetesClusterNatGatewayManagedOutboundIPCountValidation(t *testing.T) {
	networkProfile := RegisterResource().Schema["network_profile"].Elem.(*pluginsdk.Resource)
	natGatewayProfile := networkProfile.Schema["nat_gateway_profile"].Elem.(*pluginsdk.Resource)
	count := natGatewayProfile.Schema["managed_outbound_ip_count"]
	tests := []struct {
		name  string
		value int
		valid bool
	}{
		{name: "below minimum", value: 0},
		{name: "minimum", value: 1, valid: true},
		{name: "maximum", value: 16, valid: true},
		{name: "above maximum", value: 17},
		{name: "previous schema maximum", value: 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, errors := count.ValidateFunc(test.value, "managed_outbound_ip_count")
			if got := len(errors) == 0; got != test.valid {
				t.Fatalf("validation success = %t, want %t (errors: %v)", got, test.valid, errors)
			}
		})
	}
}

func TestKubernetesClusterNatGatewaySkuChange(t *testing.T) {
	tests := []struct {
		name    string
		old     string
		new     string
		wantErr bool
	}{
		{name: "standard to v2 is allowed", old: "Standard", new: "StandardV2"},
		{name: "v2 to standard is rejected", old: "StandardV2", new: "Standard", wantErr: true},
		{name: "omitted sku retains existing value", old: "StandardV2", new: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateKubernetesClusterNatGatewaySkuChange(test.old, test.new)
			if (err != nil) != test.wantErr {
				t.Fatalf("error = %v, want error %t", err, test.wantErr)
			}
		})
	}
}
