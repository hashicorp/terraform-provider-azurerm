// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id":                    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1",
				"kubernetes_cluster_id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1"),
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
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
		},
		{
			name: "new id",
			input: map[string]interface{}{
				"id":                    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1",
				"kubernetes_cluster_id": "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1",
			},
			expected: utils.String("/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
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
