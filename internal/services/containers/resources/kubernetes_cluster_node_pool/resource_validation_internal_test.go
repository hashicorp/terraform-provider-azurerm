// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

func TestKubernetesClusterNodePoolNameValidation(t *testing.T) {
	for _, maxLength := range []int{12, 17, 39} {
		for _, poolName := range []string{
			strings.Repeat("a", 12),
			strings.Repeat("a", 13),
			"linuxworkerpool01",
			strings.Repeat("a", 18),
			strings.Repeat("a", 39),
		} {
			for _, rotation := range []bool{false, true} {
				testName := fmt.Sprintf("maximum_%d/length_%d/rotation_%t", maxLength, len(poolName), rotation)
				t.Run(testName, func(t *testing.T) {
					poolConfig := map[string]any{
						"name":       poolName,
						"vm_size":    "Standard_D2s_v3",
						"node_count": 1,
					}
					if rotation {
						poolConfig["name"] = "shortpool"
						poolConfig["temporary_name_for_rotation"] = poolName
					}

					config := poolConfig
					config["kubernetes_cluster_id"] = "/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/test/providers/Microsoft.ContainerService/managedClusters/test"

					userFeatures := features.Default()
					if maxLength != 12 {
						userFeatures.KubernetesCluster.NodePoolNameMaxLength = maxLength
					}
					client := &clients.Client{Features: userFeatures}
					resource := RegisterResource()
					_, err := resource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), client)
					wantError := len(poolName) > maxLength
					if (err != nil) != wantError {
						t.Fatalf("expected error=%t, got %v", wantError, err)
					}
					if err != nil && !strings.Contains(err.Error(), "characters in length") {
						t.Fatalf("unexpected validation error: %v", err)
					}
				})
			}
		}
	}
}
