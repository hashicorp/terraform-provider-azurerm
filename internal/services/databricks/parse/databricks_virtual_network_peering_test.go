package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = DatabricksVirtualNetworkPeeringId{}

func TestDatabricksVirtualNetworkPeeringIDFormatter(t *testing.T) {
	actual := NewDatabricksVirtualNetworkPeeringID("12345678-1234-9876-4563-123456789012", "resGroup1", "workspace1", "peer1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/virtualNetworkPeerings/peer1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDatabricksVirtualNetworkPeeringID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DatabricksVirtualNetworkPeeringId
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
			// missing WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/",
			Error: true,
		},

		{
			// missing value for WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/",
			Error: true,
		},

		{
			// missing VirtualNetworkPeeringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/",
			Error: true,
		},

		{
			// missing value for VirtualNetworkPeeringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/virtualNetworkPeerings/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/virtualNetworkPeerings/peer1",
			Expected: &DatabricksVirtualNetworkPeeringId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				WorkspaceName:             "workspace1",
				VirtualNetworkPeeringName: "peer1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DATABRICKS/WORKSPACES/WORKSPACE1/VIRTUALNETWORKPEERINGS/PEER1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DatabricksVirtualNetworkPeeringID(v.Input)
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
		if actual.WorkspaceName != v.Expected.WorkspaceName {
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.WorkspaceName, actual.WorkspaceName)
		}
		if actual.VirtualNetworkPeeringName != v.Expected.VirtualNetworkPeeringName {
			t.Fatalf("Expected %q but got %q for VirtualNetworkPeeringName", v.Expected.VirtualNetworkPeeringName, actual.VirtualNetworkPeeringName)
		}
	}
}

func TestDatabricksVirtualNetworkPeeringIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DatabricksVirtualNetworkPeeringId
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
			// missing WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/",
			Error: true,
		},

		{
			// missing value for WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/",
			Error: true,
		},

		{
			// missing VirtualNetworkPeeringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/",
			Error: true,
		},

		{
			// missing value for VirtualNetworkPeeringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/virtualNetworkPeerings/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/virtualNetworkPeerings/peer1",
			Expected: &DatabricksVirtualNetworkPeeringId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				WorkspaceName:             "workspace1",
				VirtualNetworkPeeringName: "peer1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/virtualnetworkpeerings/peer1",
			Expected: &DatabricksVirtualNetworkPeeringId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				WorkspaceName:             "workspace1",
				VirtualNetworkPeeringName: "peer1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/WORKSPACES/workspace1/VIRTUALNETWORKPEERINGS/peer1",
			Expected: &DatabricksVirtualNetworkPeeringId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				WorkspaceName:             "workspace1",
				VirtualNetworkPeeringName: "peer1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/WoRkSpAcEs/workspace1/ViRtUaLnEtWoRkPeErInGs/peer1",
			Expected: &DatabricksVirtualNetworkPeeringId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				WorkspaceName:             "workspace1",
				VirtualNetworkPeeringName: "peer1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DatabricksVirtualNetworkPeeringIDInsensitively(v.Input)
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
		if actual.WorkspaceName != v.Expected.WorkspaceName {
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.WorkspaceName, actual.WorkspaceName)
		}
		if actual.VirtualNetworkPeeringName != v.Expected.VirtualNetworkPeeringName {
			t.Fatalf("Expected %q but got %q for VirtualNetworkPeeringName", v.Expected.VirtualNetworkPeeringName, actual.VirtualNetworkPeeringName)
		}
	}
}
