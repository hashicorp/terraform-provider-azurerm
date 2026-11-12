// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/ehnamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func geoReplicationTestInput(replicas ...map[string]any) []any {
	set := pluginsdk.NewSet(resourceGeoReplicationReplicaHash, nil)
	for _, r := range replicas {
		set.Add(r)
	}
	return []any{
		map[string]any{
			"maximum_replication_lag_in_seconds": 300,
			"replica":                            set,
		},
	}
}

func TestExpandEventHubNamespaceGeoReplication(t *testing.T) {
	if got := expandEventHubNamespaceGeoReplication([]any{}); got != nil {
		t.Fatalf("expected nil for an empty block, got %+v", got)
	}

	got := expandEventHubNamespaceGeoReplication(geoReplicationTestInput(
		map[string]any{"location": "East US 2", "role": "Primary", "dedicated_cluster_id": ""},
		map[string]any{"location": "centralus", "role": "Secondary", "dedicated_cluster_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.EventHub/clusters/c1"},
	))
	if got == nil || got.Locations == nil {
		t.Fatalf("expected geo replication properties, got nil")
	}
	if pointer.From(got.MaxReplicationLagDurationInSeconds) != 300 {
		t.Fatalf("expected max replication lag 300, got %d", pointer.From(got.MaxReplicationLagDurationInSeconds))
	}

	byLocation := map[string]ehnamespaces.NamespaceReplicaLocation{}
	for _, l := range *got.Locations {
		byLocation[pointer.From(l.LocationName)] = l
	}
	if len(byLocation) != 2 {
		t.Fatalf("expected 2 replicas, got %d", len(byLocation))
	}
	if primary, ok := byLocation["eastus2"]; !ok || pointer.From(primary.RoleType) != ehnamespaces.GeoDRRoleTypePrimary || primary.ClusterArmId != nil {
		t.Fatalf("unexpected primary replica: %+v", primary)
	}
	if secondary, ok := byLocation["centralus"]; !ok || pointer.From(secondary.RoleType) != ehnamespaces.GeoDRRoleTypeSecondary || secondary.ClusterArmId == nil {
		t.Fatalf("unexpected secondary replica: %+v", secondary)
	}
}

func TestFlattenEventHubNamespaceGeoReplication(t *testing.T) {
	cases := map[string]struct {
		input    *ehnamespaces.GeoDataReplicationProperties
		expected []any
	}{
		"nil": {
			input:    nil,
			expected: []any{},
		},
		"primary only": {
			input: &ehnamespaces.GeoDataReplicationProperties{
				Locations: &[]ehnamespaces.NamespaceReplicaLocation{
					{LocationName: pointer.To("eastus2"), RoleType: pointer.To(ehnamespaces.GeoDRRoleTypePrimary)},
				},
			},
			expected: []any{},
		},
		"primary and secondary": {
			input: &ehnamespaces.GeoDataReplicationProperties{
				MaxReplicationLagDurationInSeconds: pointer.To(int64(600)),
				Locations: &[]ehnamespaces.NamespaceReplicaLocation{
					{LocationName: pointer.To("East US 2"), RoleType: pointer.To(ehnamespaces.GeoDRRoleTypePrimary)},
					{LocationName: pointer.To("centralus"), RoleType: pointer.To(ehnamespaces.GeoDRRoleTypeSecondary), ReplicaState: pointer.To("Ready")},
				},
			},
			expected: []any{
				map[string]any{
					"maximum_replication_lag_in_seconds": 600,
					"replica": []any{
						map[string]any{"location": "eastus2", "role": "Primary", "dedicated_cluster_id": ""},
						map[string]any{"location": "centralus", "role": "Secondary", "dedicated_cluster_id": ""},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := flattenEventHubNamespaceGeoReplication(tc.input); !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("expected %+v, got %+v", tc.expected, got)
			}
		})
	}
}

func TestEventHubNamespaceGeoReplicationPrimaryOnly(t *testing.T) {
	previous := geoReplicationTestInput(
		map[string]any{"location": "centralus", "role": "Primary", "dedicated_cluster_id": ""},
		map[string]any{"location": "eastus2", "role": "Secondary", "dedicated_cluster_id": ""},
	)

	for name, tc := range map[string]struct {
		previous []any
		expected string
	}{
		"uses the previous primary":            {previous: previous, expected: "centralus"},
		"falls back to the namespace location": {previous: []any{}, expected: "eastus2"},
	} {
		t.Run(name, func(t *testing.T) {
			got := eventHubNamespaceGeoReplicationPrimaryOnly(tc.previous, "eastus2")
			if got == nil || got.Locations == nil || len(*got.Locations) != 1 {
				t.Fatalf("expected a single primary location, got %+v", got)
			}
			only := (*got.Locations)[0]
			if pointer.From(only.LocationName) != tc.expected || pointer.From(only.RoleType) != ehnamespaces.GeoDRRoleTypePrimary {
				t.Fatalf("expected primary %q, got %+v", tc.expected, only)
			}
		})
	}
}

func TestResourceGeoReplicationReplicaHash(t *testing.T) {
	a := resourceGeoReplicationReplicaHash(map[string]any{"location": "East US 2", "role": "Primary", "dedicated_cluster_id": ""})
	b := resourceGeoReplicationReplicaHash(map[string]any{"location": "eastus2", "role": "primary", "dedicated_cluster_id": ""})
	if a != b {
		t.Fatalf("expected location and role to hash insensitively")
	}

	c := resourceGeoReplicationReplicaHash(map[string]any{"location": "eastus2", "role": "Secondary", "dedicated_cluster_id": ""})
	if a == c {
		t.Fatalf("expected different roles to hash differently")
	}
}
