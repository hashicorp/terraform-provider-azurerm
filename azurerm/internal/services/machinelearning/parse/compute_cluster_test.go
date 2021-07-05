package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ComputeClusterId{}

func TestComputeClusterIDFormatter(t *testing.T) {
	actual := NewComputeClusterID("00000000-0000-0000-0000-000000000000", "resGroup1", "workspace1", "cluster1").ID()
	expected := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestComputeClusterID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ComputeClusterId
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},

		{
			// missing WorkspaceName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/",
			Error: true,
		},

		{
			// missing value for WorkspaceName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/",
			Error: true,
		},

		{
			// missing ComputeName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/",
			Error: true,
		},

		{
			// missing value for ComputeName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1",
			Expected: &ComputeClusterId{
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
				ResourceGroup:  "resGroup1",
				WorkspaceName:  "workspace1",
				ComputeName:    "cluster1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.MACHINELEARNINGSERVICES/WORKSPACES/WORKSPACE1/COMPUTES/CLUSTER1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ComputeClusterID(v.Input)
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
		if actual.ComputeName != v.Expected.ComputeName {
			t.Fatalf("Expected %q but got %q for ComputeName", v.Expected.ComputeName, actual.ComputeName)
		}
	}
}
