// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"encoding/json"
	"maps"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesClusterAPIAccessProfileEmptyPlan(t *testing.T) {
	resource := &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
		"api_server_access_profile": resourceKubernetesCluster().Schema["api_server_access_profile"],
	}}
	profiles := map[string]*managedclusters.ManagedClusterAPIServerAccessProfile{
		"nil profile":   nil,
		"nil fields":    {},
		"empty ranges":  {AuthorizedIPRanges: pointer.To([]string{})},
		"API defaults":  {AuthorizedIPRanges: pointer.To([]string{}), EnableVnetIntegration: pointer.To(false), SubnetId: pointer.To("")},
		"private flags": {EnablePrivateCluster: pointer.To(true), EnablePrivateClusterPublicFQDN: pointer.To(true), DisableRunCommand: pointer.To(true)},
	}

	tests := []struct {
		name   string
		config map[string]interface{}
	}{
		{name: "removed block", config: map[string]interface{}{}},
		{name: "null block", config: map[string]interface{}{"api_server_access_profile": nil}},
		{name: "empty ranges", config: map[string]interface{}{"api_server_access_profile": []interface{}{map[string]interface{}{"authorized_ip_ranges": []interface{}{}}}}},
		{name: "null ranges", config: map[string]interface{}{"api_server_access_profile": []interface{}{map[string]interface{}{"authorized_ip_ranges": nil}}}},
		{name: "omitted ranges", config: map[string]interface{}{"api_server_access_profile": []interface{}{map[string]interface{}{}}}},
	}
	for name, profile := range profiles {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					data := schema.TestResourceDataRaw(t, resource.Schema, test.config)
					data.SetId("existing-cluster")
					configured := len(data.Get("api_server_access_profile").([]interface{})) > 0
					if err := data.Set("api_server_access_profile", flattenKubernetesClusterAPIAccessProfile(profile, configured)); err != nil {
						t.Fatal(err)
					}
					diff, err := resource.SimpleDiff(context.Background(), data.State(), terraform.NewResourceConfigRaw(test.config), nil)
					if err != nil {
						t.Fatal(err)
					}
					if !diff.Empty() {
						t.Fatalf("empty API response produced a follow-up diff: %#v", diff.Attributes)
					}
				})
			}
		})
	}
}

func TestKubernetesClusterAPIAccessProfileExpand(t *testing.T) {
	tests := []struct {
		name   string
		config map[string]interface{}
		want   map[string]interface{}
	}{
		{
			name:   "omitted block",
			config: map[string]interface{}{},
		},
		{
			name: "empty block",
			config: map[string]interface{}{
				"api_server_access_profile": []interface{}{map[string]interface{}{}},
			},
			want: map[string]interface{}{"enableVnetIntegration": false},
		},
		{
			name: "positive ranges",
			config: map[string]interface{}{
				"api_server_access_profile": []interface{}{map[string]interface{}{"authorized_ip_ranges": []interface{}{"8.8.8.8/32"}}},
			},
			want: map[string]interface{}{
				"authorizedIPRanges":    []interface{}{"8.8.8.8/32"},
				"enableVnetIntegration": false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, resourceKubernetesCluster().Schema, test.config)
			payload, err := json.Marshal(expandKubernetesClusterAPIAccessProfile(data))
			if err != nil {
				t.Fatal(err)
			}
			var actual map[string]interface{}
			if err := json.Unmarshal(payload, &actual); err != nil {
				t.Fatal(err)
			}
			expected := map[string]interface{}{
				"enablePrivateCluster":           false,
				"enablePrivateClusterPublicFQDN": false,
				"disableRunCommand":              false,
			}
			maps.Copy(expected, test.want)
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("payload = %s, want %#v", payload, expected)
			}
		})
	}
}

func TestKubernetesClusterAPIAccessProfileUpdate(t *testing.T) {
	const subnetID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Network/virtualNetworks/test/subnets/api"
	positive := []interface{}{map[string]interface{}{"authorized_ip_ranges": []interface{}{"8.8.8.8/32"}}}
	withVnet := []interface{}{map[string]interface{}{
		"authorized_ip_ranges":                []interface{}{"8.8.8.8/32"},
		"virtual_network_integration_enabled": true,
		"subnet_id":                           subnetID,
	}}
	tests := []struct {
		name             string
		oldProfile       []interface{}
		newProfile       []interface{}
		nullBlock        bool
		toggleRunCommand bool
		privateCluster   bool
		wantRanges       []interface{}
		wantVnet         bool
		wantSubnet       string
	}{
		{name: "removed block", oldProfile: positive, wantRanges: []interface{}{}},
		{name: "null block", oldProfile: positive, nullBlock: true, wantRanges: []interface{}{}},
		{name: "empty ranges", oldProfile: positive, newProfile: []interface{}{map[string]interface{}{"authorized_ip_ranges": []interface{}{}}}, wantRanges: []interface{}{}},
		{name: "null ranges", oldProfile: positive, newProfile: []interface{}{map[string]interface{}{"authorized_ip_ranges": nil}}, wantRanges: []interface{}{}},
		{name: "omitted ranges", oldProfile: positive, newProfile: []interface{}{map[string]interface{}{}}, wantRanges: []interface{}{}},
		{name: "positive change", oldProfile: positive, newProfile: []interface{}{map[string]interface{}{"authorized_ip_ranges": []interface{}{"1.1.1.1/32"}}}, wantRanges: []interface{}{"1.1.1.1/32"}},
		{name: "ranges retained during run command update", oldProfile: positive, newProfile: positive, toggleRunCommand: true, wantRanges: []interface{}{"8.8.8.8/32"}},
		{name: "absent ranges during run command update", toggleRunCommand: true},
		{name: "private flags during run command update", toggleRunCommand: true, privateCluster: true},
		{
			name: "VNet retained with omitted ranges", oldProfile: withVnet,
			newProfile: []interface{}{map[string]interface{}{"virtual_network_integration_enabled": true, "subnet_id": subnetID}},
			wantRanges: []interface{}{}, wantVnet: true, wantSubnet: subnetID,
		},
		{
			name: "VNet retained with empty ranges", oldProfile: withVnet,
			newProfile: []interface{}{map[string]interface{}{"authorized_ip_ranges": []interface{}{}, "virtual_network_integration_enabled": true, "subnet_id": subnetID}},
			wantRanges: []interface{}{}, wantVnet: true, wantSubnet: subnetID,
		},
		{
			name: "VNet retained with null ranges", oldProfile: withVnet,
			newProfile: []interface{}{map[string]interface{}{"authorized_ip_ranges": nil, "virtual_network_integration_enabled": true, "subnet_id": subnetID}},
			wantRanges: []interface{}{}, wantVnet: true, wantSubnet: subnetID,
		},
	}
	productionSchema := resourceKubernetesCluster().Schema
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resource := &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{}}
			for _, key := range []string{"api_server_access_profile", "private_cluster_enabled", "private_cluster_public_fqdn_enabled", "run_command_enabled"} {
				resource.Schema[key] = productionSchema[key]
			}
			oldConfig := map[string]interface{}{
				"run_command_enabled":                 true,
				"private_cluster_enabled":             test.privateCluster,
				"private_cluster_public_fqdn_enabled": test.privateCluster,
			}
			newConfig := map[string]interface{}{
				"run_command_enabled":                 !test.toggleRunCommand,
				"private_cluster_enabled":             test.privateCluster,
				"private_cluster_public_fqdn_enabled": test.privateCluster,
			}
			if test.oldProfile != nil {
				oldConfig["api_server_access_profile"] = test.oldProfile
			}
			if test.newProfile != nil {
				newConfig["api_server_access_profile"] = test.newProfile
			} else if test.nullBlock {
				newConfig["api_server_access_profile"] = nil
			}
			data := schema.TestResourceDataRaw(t, resource.Schema, oldConfig)
			data.SetId("existing-cluster")
			diff, err := resource.SimpleDiff(context.Background(), data.State(), terraform.NewResourceConfigRaw(newConfig), nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff.Empty() || diff.RequiresNew() {
				t.Fatalf("expected an in-place update, got %#v", diff)
			}
			var profile *managedclusters.ManagedClusterAPIServerAccessProfile
			resource.UpdateContext = func(_ context.Context, d *pluginsdk.ResourceData, _ interface{}) diag.Diagnostics {
				profile = expandKubernetesClusterAPIAccessProfileUpdate(d)
				return nil
			}
			if _, diagnostics := resource.Apply(context.Background(), data.State(), diff, nil); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			if profile == nil {
				t.Fatal("update did not produce an API access profile")
			}
			payload, err := json.Marshal(profile)
			if err != nil {
				t.Fatal(err)
			}
			var actual map[string]interface{}
			if err := json.Unmarshal(payload, &actual); err != nil {
				t.Fatal(err)
			}
			expected := map[string]interface{}{
				"enablePrivateCluster":           test.privateCluster,
				"enablePrivateClusterPublicFQDN": test.privateCluster,
				"disableRunCommand":              test.toggleRunCommand,
			}
			if test.wantRanges != nil {
				expected["authorizedIPRanges"] = test.wantRanges
			}
			if test.newProfile != nil {
				expected["enableVnetIntegration"] = test.wantVnet
			}
			if test.wantSubnet != "" {
				expected["subnetId"] = test.wantSubnet
			}
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("update payload = %s, want %#v", payload, expected)
			}
		})
	}
}

func TestKubernetesClusterAPIAccessProfileFlatten(t *testing.T) {
	tests := []struct {
		name    string
		profile managedclusters.ManagedClusterAPIServerAccessProfile
		want    map[string]interface{}
	}{
		{
			name:    "positive ranges",
			profile: managedclusters.ManagedClusterAPIServerAccessProfile{AuthorizedIPRanges: pointer.To([]string{"8.8.8.8/32"})},
			want:    map[string]interface{}{"authorized_ip_ranges": []interface{}{"8.8.8.8/32"}, "virtual_network_integration_enabled": false, "subnet_id": ""},
		},
		{
			name:    "VNet integration",
			profile: managedclusters.ManagedClusterAPIServerAccessProfile{EnableVnetIntegration: pointer.To(true)},
			want:    map[string]interface{}{"authorized_ip_ranges": []interface{}{}, "virtual_network_integration_enabled": true, "subnet_id": ""},
		},
		{
			name:    "subnet only",
			profile: managedclusters.ManagedClusterAPIServerAccessProfile{SubnetId: pointer.To("subnet")},
			want:    map[string]interface{}{"authorized_ip_ranges": []interface{}{}, "virtual_network_integration_enabled": false, "subnet_id": "subnet"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := flattenKubernetesClusterAPIAccessProfile(&test.profile, false)
			if !reflect.DeepEqual(actual, []interface{}{test.want}) {
				t.Fatalf("profile = %#v, want %#v", actual, test.want)
			}
		})
	}
}
