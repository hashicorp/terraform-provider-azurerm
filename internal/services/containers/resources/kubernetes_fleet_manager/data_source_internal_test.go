// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_manager

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
)

func TestMapFleetToKubernetesFleetManagerDataSourceModel(t *testing.T) {
	state := KubernetesFleetManagerDataSourceModel{
		HubProfile: []FleetManagerHubProfile{{DnsPrefix: "stale"}},
	}

	if err := mapFleetToKubernetesFleetManagerDataSourceModel(fleets.Fleet{}, &state); err != nil {
		t.Fatal(err)
	}
	if len(state.HubProfile) != 0 {
		t.Fatalf("expected missing API hub profile to clear state, got %#v", state.HubProfile)
	}
}
