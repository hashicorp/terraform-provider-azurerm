// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
)

func TestFlattenClusterPoolNetworkProfileDranet(t *testing.T) {
	tests := []struct {
		name             string
		input            *managedclusters.AgentPoolNetworkProfile
		hasExistingBlock bool
		expectedCount    int
		expectedMode     string
	}{
		{name: "omitted profile"},
		{name: "explicit block with omitted API profile", hasExistingBlock: true, expectedCount: 1, expectedMode: "Unmanaged"},
		{name: "API default without configuration", input: &managedclusters.AgentPoolNetworkProfile{Dranet: &managedclusters.DRANETProfile{Mode: pointer.To(managedclusters.DRANETModeUnmanaged)}}},
		{name: "explicit unmanaged block", input: &managedclusters.AgentPoolNetworkProfile{Dranet: &managedclusters.DRANETProfile{Mode: pointer.To(managedclusters.DRANETModeUnmanaged)}}, hasExistingBlock: true, expectedCount: 1, expectedMode: "Unmanaged"},
		{name: "managed profile without configuration", input: &managedclusters.AgentPoolNetworkProfile{Dranet: &managedclusters.DRANETProfile{Mode: pointer.To(managedclusters.DRANETModeManaged)}}, expectedCount: 1, expectedMode: "Managed"},
		{name: "other network settings", input: &managedclusters.AgentPoolNetworkProfile{ApplicationSecurityGroups: pointer.To([]string{"test"})}, expectedCount: 1, expectedMode: "Unmanaged"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := FlattenClusterPoolNetworkProfile(test.input, test.hasExistingBlock)
			if len(actual) != test.expectedCount {
				t.Fatalf("expected %d blocks, got %#v", test.expectedCount, actual)
			}
			if len(actual) > 0 && actual[0].(map[string]any)["dranet_mode"] != test.expectedMode {
				t.Fatalf("expected DRANET mode %q, got %#v", test.expectedMode, actual)
			}
		})
	}
}
