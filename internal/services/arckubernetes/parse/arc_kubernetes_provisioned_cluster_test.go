// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ArcKubernetesProvisionedClusterId{}

func TestArcKubernetesProvisionedClusterIDFormatter(t *testing.T) {
	actual := NewArcKubernetesProvisionedClusterID("12345678-1234-9876-4563-123456789012", "group1", "cluster1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestArcKubernetesProvisionedClusterID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ArcKubernetesProvisionedClusterId
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
			// missing ConnectedClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/",
			Error: true,
		},

		{
			// missing value for ConnectedClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/",
			Error: true,
		},

		{
			// missing ProvisionedClusterInstanceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.HybridContainerService/",
			Error: true,
		},

		{
			// missing value for ProvisionedClusterInstanceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.HybridContainerService/provisionedClusterInstances/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default",
			Expected: &ArcKubernetesProvisionedClusterId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                  "group1",
				ConnectedClusterName:           "cluster1",
				ProvisionedClusterInstanceName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.KUBERNETES/CONNECTEDCLUSTERS/CLUSTER1/PROVIDERS/MICROSOFT.HYBRIDCONTAINERSERVICE/PROVISIONEDCLUSTERINSTANCES/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ArcKubernetesProvisionedClusterID(v.Input)
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
		if actual.ConnectedClusterName != v.Expected.ConnectedClusterName {
			t.Fatalf("Expected %q but got %q for ConnectedClusterName", v.Expected.ConnectedClusterName, actual.ConnectedClusterName)
		}
		if actual.ProvisionedClusterInstanceName != v.Expected.ProvisionedClusterInstanceName {
			t.Fatalf("Expected %q but got %q for ProvisionedClusterInstanceName", v.Expected.ProvisionedClusterInstanceName, actual.ProvisionedClusterInstanceName)
		}
	}
}

func TestArcKubernetesProvisionedClusterIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ArcKubernetesProvisionedClusterId
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
			// missing ConnectedClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/",
			Error: true,
		},

		{
			// missing value for ConnectedClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/",
			Error: true,
		},

		{
			// missing ProvisionedClusterInstanceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.HybridContainerService/",
			Error: true,
		},

		{
			// missing value for ProvisionedClusterInstanceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.HybridContainerService/provisionedClusterInstances/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default",
			Expected: &ArcKubernetesProvisionedClusterId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                  "group1",
				ConnectedClusterName:           "cluster1",
				ProvisionedClusterInstanceName: "default",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedclusters/cluster1/providers/Microsoft.HybridContainerService/provisionedclusterinstances/default",
			Expected: &ArcKubernetesProvisionedClusterId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                  "group1",
				ConnectedClusterName:           "cluster1",
				ProvisionedClusterInstanceName: "default",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/CONNECTEDCLUSTERS/cluster1/providers/Microsoft.HybridContainerService/PROVISIONEDCLUSTERINSTANCES/default",
			Expected: &ArcKubernetesProvisionedClusterId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                  "group1",
				ConnectedClusterName:           "cluster1",
				ProvisionedClusterInstanceName: "default",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Kubernetes/CoNnEcTeDcLuStErS/cluster1/providers/Microsoft.HybridContainerService/PrOvIsIoNeDcLuStErInStAnCeS/default",
			Expected: &ArcKubernetesProvisionedClusterId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                  "group1",
				ConnectedClusterName:           "cluster1",
				ProvisionedClusterInstanceName: "default",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ArcKubernetesProvisionedClusterIDInsensitively(v.Input)
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
		if actual.ConnectedClusterName != v.Expected.ConnectedClusterName {
			t.Fatalf("Expected %q but got %q for ConnectedClusterName", v.Expected.ConnectedClusterName, actual.ConnectedClusterName)
		}
		if actual.ProvisionedClusterInstanceName != v.Expected.ProvisionedClusterInstanceName {
			t.Fatalf("Expected %q but got %q for ProvisionedClusterInstanceName", v.Expected.ProvisionedClusterInstanceName, actual.ProvisionedClusterInstanceName)
		}
	}
}
