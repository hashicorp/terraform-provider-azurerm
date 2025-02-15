// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = NodePoolId{}

func TestNodePoolIDFormatter(t *testing.T) {
	actual := NewNodePoolID("12345678-1234-9876-4563-123456789012", "resGroup1", "cluster1", "pool1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestNodePoolID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NodePoolId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing ManagedClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/",
			Error: true,
		},

		{
			// missing value for ManagedClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/",
			Error: true,
		},

		{
			// missing AgentPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/",
			Error: true,
		},

		{
			// missing value for AgentPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1",
			Expected: &NodePoolId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				ManagedClusterName: "cluster1",
				AgentPoolName:      "pool1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.CONTAINERSERVICE/MANAGEDCLUSTERS/CLUSTER1/AGENTPOOLS/POOL1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := NodePoolID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.ManagedClusterName != v.Expected.ManagedClusterName {
			t.Fatalf("Expected %q but got %q for ManagedClusterName", v.Expected.ManagedClusterName, actual.ManagedClusterName)
		}
		if actual.AgentPoolName != v.Expected.AgentPoolName {
			t.Fatalf("Expected %q but got %q for AgentPoolName", v.Expected.AgentPoolName, actual.AgentPoolName)
		}
	}
}
