// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"encoding/json"
	"maps"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestKubernetesClusterAPIAccessProfileExpand(t *testing.T) {
	tests := []struct {
		name   string
		config map[string]any
		want   map[string]any
	}{
		{
			name:   "omitted block",
			config: map[string]any{},
		},
		{
			name: "empty block",
			config: map[string]any{
				"api_server_access_profile": []any{map[string]any{}},
			},
			want: map[string]any{"enableVnetIntegration": false},
		},
		{
			name: "positive ranges",
			config: map[string]any{
				"api_server_access_profile": []any{map[string]any{"authorized_ip_ranges": []any{"8.8.8.8/32"}}},
			},
			want: map[string]any{
				"authorizedIPRanges":    []any{"8.8.8.8/32"},
				"enableVnetIntegration": false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, RegisterResource().Schema, test.config)
			payload, err := json.Marshal(expandKubernetesClusterAPIAccessProfile(data))
			if err != nil {
				t.Fatal(err)
			}
			var actual map[string]any
			if err := json.Unmarshal(payload, &actual); err != nil {
				t.Fatal(err)
			}
			expected := map[string]any{
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

func TestKubernetesClusterAPIAccessProfileFlatten(t *testing.T) {
	emptyBlock := map[string]any{"authorized_ip_ranges": []any{}, "virtual_network_integration_enabled": false, "subnet_id": ""}
	tests := []struct {
		name          string
		profile       *managedclusters.ManagedClusterAPIServerAccessProfile
		preserveEmpty bool
		want          []any
	}{
		{name: "nil profile", want: []any{}},
		{
			name: "API defaults without a configured block",
			profile: &managedclusters.ManagedClusterAPIServerAccessProfile{
				AuthorizedIPRanges:    pointer.To([]string{}),
				EnableVnetIntegration: pointer.To(false),
				SubnetId:              pointer.To(""),
			},
			want: []any{},
		},
		{name: "configured empty block", preserveEmpty: true, want: []any{emptyBlock}},
		{
			name:    "positive ranges",
			profile: &managedclusters.ManagedClusterAPIServerAccessProfile{AuthorizedIPRanges: pointer.To([]string{"8.8.8.8/32"})},
			want:    []any{map[string]any{"authorized_ip_ranges": []any{"8.8.8.8/32"}, "virtual_network_integration_enabled": false, "subnet_id": ""}},
		},
		{
			name:    "VNet integration",
			profile: &managedclusters.ManagedClusterAPIServerAccessProfile{EnableVnetIntegration: pointer.To(true)},
			want:    []any{map[string]any{"authorized_ip_ranges": []any{}, "virtual_network_integration_enabled": true, "subnet_id": ""}},
		},
		{
			name:    "subnet only",
			profile: &managedclusters.ManagedClusterAPIServerAccessProfile{SubnetId: pointer.To("subnet")},
			want:    []any{map[string]any{"authorized_ip_ranges": []any{}, "virtual_network_integration_enabled": false, "subnet_id": "subnet"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := flattenKubernetesClusterAPIAccessProfile(test.profile, test.preserveEmpty)
			if !reflect.DeepEqual(actual, test.want) {
				t.Fatalf("profile = %#v, want %#v", actual, test.want)
			}
		})
	}
}
