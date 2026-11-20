// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
)

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
