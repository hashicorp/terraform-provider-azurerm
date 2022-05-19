package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ClusterSnapshotId{}

func TestClusterSnapshotIDFormatter(t *testing.T) {
	actual := NewClusterSnapshotID("12345678-1234-9876-4563-123456789012", "resGroup1", "managedclustersnapshot1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedclustersnapshots/managedclustersnapshot1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestClusterSnapshotID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ClusterSnapshotId
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
			// missing ManagedclustersnapshotName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/",
			Error: true,
		},

		{
			// missing value for ManagedclustersnapshotName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedclustersnapshots/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedclustersnapshots/managedclustersnapshot1",
			Expected: &ClusterSnapshotId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				ManagedclustersnapshotName: "managedclustersnapshot1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.CONTAINERSERVICE/MANAGEDCLUSTERSNAPSHOTS/MANAGEDCLUSTERSNAPSHOT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ClusterSnapshotID(v.Input)
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
		if actual.ManagedclustersnapshotName != v.Expected.ManagedclustersnapshotName {
			t.Fatalf("Expected %q but got %q for ManagedclustersnapshotName", v.Expected.ManagedclustersnapshotName, actual.ManagedclustersnapshotName)
		}
	}
}
