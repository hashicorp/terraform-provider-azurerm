// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestKubernetesClusterNodePoolV0ToV1_id(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]interface{}
		expected *string
	}{
		{
			name: "missing id",
			input: map[string]interface{}{
				"id":                    "",
				"kubernetes_cluster_id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1",
			},
			expected: nil,
		},
		{
			name: "old id",
			input: map[string]interface{}{
				"id":                    "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1",
				"kubernetes_cluster_id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id":                    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1",
				"kubernetes_cluster_id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1"),
		},
		{
			name: "mixed-case resource segments",
			input: map[string]interface{}{
				"id":                    "/SUBSCRIPTIONS/12345678-1234-5678-1234-123456789012/RESOURCEGROUPS/Group1/PROVIDERS/microsoft.containerservice/MANAGEDCLUSTERS/Cluster1/AGENTPOOLS/Pool1",
				"kubernetes_cluster_id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/Group1/providers/Microsoft.ContainerService/managedClusters/Cluster1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/Group1/providers/Microsoft.ContainerService/managedClusters/Cluster1/agentPools/Pool1"),
		},
	}
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			result, err := KubernetesClusterNodePoolV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, nil)
			if err != nil && test.expected == nil {
				return
			} else {
				if err == nil && test.expected == nil {
					t.Fatalf("Expected an error but didn't get one")
				} else if err != nil && test.expected != nil {
					t.Fatalf("Expected no error but got: %+v", err)
				}
			}

			actualId := result["id"].(string)
			if *test.expected != actualId {
				t.Fatalf("expected %q but got %q!", *test.expected, actualId)
			}
		})
	}
}

func TestKubernetesClusterNodePoolV0ToV1_kubernetes_cluster_id(t *testing.T) {
	testData := []struct {
		name     string
		input    map[string]interface{}
		expected *string
	}{
		{
			name: "missing id",
			input: map[string]interface{}{
				"id":                    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1",
				"kubernetes_cluster_id": "",
			},
			expected: nil,
		},
		{
			name: "old id",
			input: map[string]interface{}{
				"id":                    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1",
				"kubernetes_cluster_id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id":                    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1",
				"kubernetes_cluster_id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
		},
		{
			name: "mixed-case resource segments",
			input: map[string]interface{}{
				"id":                    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/Group1/providers/Microsoft.ContainerService/managedClusters/Cluster1/agentPools/Pool1",
				"kubernetes_cluster_id": "/SUBSCRIPTIONS/12345678-1234-5678-1234-123456789012/RESOURCEGROUPS/Group1/PROVIDERS/microsoft.containerservice/MANAGEDCLUSTERS/Cluster1",
			},
			expected: pointer.To("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/Group1/providers/Microsoft.ContainerService/managedClusters/Cluster1"),
		},
	}
	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			result, err := KubernetesClusterNodePoolV0ToV1{}.UpgradeFunc()(context.TODO(), test.input, nil)
			if err != nil && test.expected == nil {
				return
			} else {
				if err == nil && test.expected == nil {
					t.Fatalf("Expected an error but didn't get one")
				} else if err != nil && test.expected != nil {
					t.Fatalf("Expected no error but got: %+v", err)
				}
			}

			actualId := result["kubernetes_cluster_id"].(string)
			if *test.expected != actualId {
				t.Fatalf("expected %q but got %q!", *test.expected, actualId)
			}
		})
	}
}
