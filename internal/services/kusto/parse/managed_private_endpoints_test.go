// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ManagedPrivateEndpointsId{}

func TestManagedPrivateEndpointsIDFormatter(t *testing.T) {
	actual := NewManagedPrivateEndpointsID("12345678-1234-9876-4563-123456789012", "resGroup1", "cluster1", "endpoint1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/managedPrivateEndpoints/endpoint1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestManagedPrivateEndpointsID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagedPrivateEndpointsId
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
			// missing ClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/",
			Error: true,
		},

		{
			// missing value for ClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/",
			Error: true,
		},

		{
			// missing ManagedPrivateEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/",
			Error: true,
		},

		{
			// missing value for ManagedPrivateEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/managedPrivateEndpoints/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/managedPrivateEndpoints/endpoint1",
			Expected: &ManagedPrivateEndpointsId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				ClusterName:                "cluster1",
				ManagedPrivateEndpointName: "endpoint1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.KUSTO/CLUSTERS/CLUSTER1/MANAGEDPRIVATEENDPOINTS/ENDPOINT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ManagedPrivateEndpointsID(v.Input)
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
		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for ClusterName", v.Expected.ClusterName, actual.ClusterName)
		}
		if actual.ManagedPrivateEndpointName != v.Expected.ManagedPrivateEndpointName {
			t.Fatalf("Expected %q but got %q for ManagedPrivateEndpointName", v.Expected.ManagedPrivateEndpointName, actual.ManagedPrivateEndpointName)
		}
	}
}

func TestManagedPrivateEndpointsIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagedPrivateEndpointsId
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
			// missing ClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/",
			Error: true,
		},

		{
			// missing value for ClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/",
			Error: true,
		},

		{
			// missing ManagedPrivateEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/",
			Error: true,
		},

		{
			// missing value for ManagedPrivateEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/managedPrivateEndpoints/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/managedPrivateEndpoints/endpoint1",
			Expected: &ManagedPrivateEndpointsId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				ClusterName:                "cluster1",
				ManagedPrivateEndpointName: "endpoint1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/managedprivateendpoints/endpoint1",
			Expected: &ManagedPrivateEndpointsId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				ClusterName:                "cluster1",
				ManagedPrivateEndpointName: "endpoint1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/CLUSTERS/cluster1/MANAGEDPRIVATEENDPOINTS/endpoint1",
			Expected: &ManagedPrivateEndpointsId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				ClusterName:                "cluster1",
				ManagedPrivateEndpointName: "endpoint1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/ClUsTeRs/cluster1/MaNaGeDpRiVaTeEnDpOiNtS/endpoint1",
			Expected: &ManagedPrivateEndpointsId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				ClusterName:                "cluster1",
				ManagedPrivateEndpointName: "endpoint1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ManagedPrivateEndpointsIDInsensitively(v.Input)
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
		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for ClusterName", v.Expected.ClusterName, actual.ClusterName)
		}
		if actual.ManagedPrivateEndpointName != v.Expected.ManagedPrivateEndpointName {
			t.Fatalf("Expected %q but got %q for ManagedPrivateEndpointName", v.Expected.ManagedPrivateEndpointName, actual.ManagedPrivateEndpointName)
		}
	}
}
