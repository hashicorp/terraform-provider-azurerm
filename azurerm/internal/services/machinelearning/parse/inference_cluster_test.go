package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

func TestInferenceClusterID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *InferenceClusterId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1",
			Error: true,
		},
		{
			Name:  "Missing Workspace Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/",
			Error: true,
		},
		{
			Name:  "Missing Machine Learning Workspace Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1",
			Error: true,
		},
		{
			Name:  "Machine Learning Inference Cluster ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1",
			Error: false,
			Expect: &InferenceClusterId{
				ResourceGroup:        "resGroup1",
				WorkspaceName:        "workspace1",
				InferenceClusterName: "cluster1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := InferenceClusterID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.WorkspaceName != v.Expect.WorkspaceName {
			t.Fatalf("Expected %q but got %q for Workspace Name", v.Expect.WorkspaceName, actual.WorkspaceName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.InferenceClusterName != v.Expect.InferenceClusterName {
			t.Fatalf("Expected %q but got %q for InferenceClusterName", v.Expect.InferenceClusterName, actual.InferenceClusterName)
		}
	}
}

var _ resourceid.Formatter = KubernetesClusterId{}

func TestKubernetesClusterIDFormatter(t *testing.T) {
	actual := NewClusterID("12345678-1234-9876-4563-123456789012", "resGroup1", "cluster1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestKubernetesClusterID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *KubernetesClusterId
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
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1",
			Expected: &KubernetesClusterId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				ManagedClusterName: "cluster1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.CONTAINERSERVICE/MANAGEDCLUSTERS/CLUSTER1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := KubernetesClusterID(v.Input)
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
	}
}
