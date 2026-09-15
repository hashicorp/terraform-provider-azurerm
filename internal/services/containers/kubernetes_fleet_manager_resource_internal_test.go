// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
)

func TestKubernetesFleetManagerHubProfileMapping(t *testing.T) {
	r := KubernetesFleetManagerResource{}

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
				DnsPrefix:         "fleet-test",
				Fqdn:              "read-only.example",
				KubernetesVersion: "1.32.0",
				PortalFqdn:        "read-only-portal.example",
			}},
		}
		var payload fleets.Fleet
		r.mapKubernetesFleetManagerResourceSchemaToFleet(config, &payload)
		expected := &fleets.FleetHubProfile{DnsPrefix: pointer.To("fleet-test")}
		if !reflect.DeepEqual(payload.Properties.HubProfile, expected) {
			t.Fatalf("expected only configured DNS prefix in the request, got %#v", payload.Properties.HubProfile)
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
				DnsPrefix:         pointer.To("fleet-test"),
				Fqdn:              pointer.To("fleet.example"),
				KubernetesVersion: pointer.To("1.32.0"),
				PortalFqdn:        pointer.To("portal.example"),
			},
		}}
		var state KubernetesFleetManagerResourceSchema
		r.mapFleetToKubernetesFleetManagerResourceSchema(payload, &state)
		expected := []FleetManagerHubProfile{{
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
			r.mapFleetToKubernetesFleetManagerResourceSchema(fleets.Fleet{Properties: properties}, &state)
			if len(state.HubProfile) != 0 {
				t.Fatalf("expected missing API hub profile to clear state, got %#v", state.HubProfile)
			}
		})
	}
}
