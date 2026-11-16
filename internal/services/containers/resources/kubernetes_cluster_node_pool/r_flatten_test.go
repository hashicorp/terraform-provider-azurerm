// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/agentpools"
)

func TestFlattenAgentPoolNetworkProfileDranet(t *testing.T) {
	tests := []struct {
		name             string
		input            *agentpools.AgentPoolNetworkProfile
		hasExistingBlock bool
		expectedCount    int
		expectedMode     string
	}{
		{name: "omitted profile"},
		{name: "explicit block with omitted API profile", hasExistingBlock: true, expectedCount: 1, expectedMode: "Unmanaged"},
		{name: "API default without configuration", input: &agentpools.AgentPoolNetworkProfile{Dranet: &agentpools.DRANETProfile{Mode: pointer.To(agentpools.DRANETModeUnmanaged)}}},
		{name: "explicit unmanaged block", input: &agentpools.AgentPoolNetworkProfile{Dranet: &agentpools.DRANETProfile{Mode: pointer.To(agentpools.DRANETModeUnmanaged)}}, hasExistingBlock: true, expectedCount: 1, expectedMode: "Unmanaged"},
		{name: "managed profile without configuration", input: &agentpools.AgentPoolNetworkProfile{Dranet: &agentpools.DRANETProfile{Mode: pointer.To(agentpools.DRANETModeManaged)}}, expectedCount: 1, expectedMode: "Managed"},
		{name: "other network settings", input: &agentpools.AgentPoolNetworkProfile{ApplicationSecurityGroups: pointer.To([]string{"test"})}, expectedCount: 1, expectedMode: "Unmanaged"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := flattenAgentPoolNetworkProfile(test.input, test.hasExistingBlock)
			if len(actual) != test.expectedCount {
				t.Fatalf("expected %d blocks, got %#v", test.expectedCount, actual)
			}
			if len(actual) > 0 && actual[0].(map[string]any)["dranet_mode"] != test.expectedMode {
				t.Fatalf("expected DRANET mode %q, got %#v", test.expectedMode, actual)
			}
		})
	}
}
