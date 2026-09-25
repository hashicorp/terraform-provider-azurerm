// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestMapFleetToKubernetesFleetManagerDataSourceModel(t *testing.T) {
	state := KubernetesFleetManagerDataSourceModel{
		HubProfile: []FleetManagerHubProfile{{DnsPrefix: "stale"}},
	}

	if err := mapFleetToKubernetesFleetManagerDataSourceModel(fleets.Fleet{
		Location: "westus2",
		Tags: &map[string]string{
			"environment": "terraform-acctests",
			"some_key":    "some-value",
		},
		Properties: &fleets.FleetProperties{
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
		},
	}, &state); err != nil {
		t.Fatal(err)
	}

	expectedHubProfile := []FleetManagerHubProfile{{
		AgentProfile: []FleetManagerHubAgentProfile{{
			SubnetId:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/subnet",
			VirtualMachineSize: "Standard_DS2_v2",
		}},
		ApiServerAccessProfile: []FleetManagerHubAPIServerAccessProfile{{
			EnablePrivateCluster: true,
		}},
		DnsPrefix:         "fleet-test",
		Fqdn:              "fleet.example",
		KubernetesVersion: "1.32.0",
		PortalFqdn:        "portal.example",
	}}

	if state.Location != "westus2" {
		t.Fatalf("expected normalized location, got %q", state.Location)
	}
	if !reflect.DeepEqual(state.HubProfile, expectedHubProfile) {
		t.Fatalf("expected hub profile %#v, got %#v", expectedHubProfile, state.HubProfile)
	}
	if !reflect.DeepEqual(state.Tags, map[string]interface{}{
		"environment": "terraform-acctests",
		"some_key":    "some-value",
	}) {
		t.Fatalf("expected tags to round-trip, got %#v", state.Tags)
	}

	if err := mapFleetToKubernetesFleetManagerDataSourceModel(fleets.Fleet{}, &state); err != nil {
		t.Fatal(err)
	}
	if len(state.HubProfile) != 0 {
		t.Fatalf("expected missing API hub profile to clear state, got %#v", state.HubProfile)
	}
}

func TestKubernetesFleetManagerDataSourceSchemaInternalValidate(t *testing.T) {
	wrapper := sdk.NewDataSourceWrapper(KubernetesFleetManagerDataSource{})
	dataSource, err := wrapper.DataSource()
	if err != nil {
		t.Fatalf("building data source schema: %+v", err)
	}

	if err := dataSource.InternalValidate(nil, false); err != nil {
		t.Fatalf("validating data source schema: %+v", err)
	}
}
