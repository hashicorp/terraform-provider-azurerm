// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
)

func TestKubernetesClusterNodePoolUpgradeSettingsFixtures(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 12345,
		Locations:     acceptance.Regions{Primary: "westus2"},
	}
	r := KubernetesClusterNodePoolResource{}
	for _, test := range []struct {
		name   string
		config string
	}{
		{name: "kubelet_partial", config: r.kubeletAndLinuxOSConfigPartial(data)},
		{name: "requires_import", config: r.requiresImportConfig(data)},
		{name: "spot", config: r.spotConfig(data)},
	} {
		t.Run(test.name, func(t *testing.T) {
			file, diagnostics := hclsyntax.ParseConfig([]byte(test.config), "test.tf", hcl.InitialPos)
			if diagnostics.HasErrors() {
				t.Fatalf("parsing fixture: %s", diagnostics.Error())
			}

			nodePools := 0
			for _, resource := range file.Body.(*hclsyntax.Body).Blocks {
				if resource.Type != "resource" || len(resource.Labels) != 2 || resource.Labels[0] != "azurerm_kubernetes_cluster_node_pool" {
					continue
				}
				nodePools++
				upgradeSettings := 0
				for _, block := range resource.Body.Blocks {
					if block.Type == "upgrade_settings" {
						upgradeSettings++
					}
				}
				if upgradeSettings != 1 {
					t.Errorf("%s: expected one upgrade_settings block, got %d", resource.Labels[1], upgradeSettings)
				}
			}
			if nodePools == 0 {
				t.Fatal("fixture does not contain a Kubernetes cluster node pool")
			}
		})
	}
}

func TestKubernetesClusterUpgradeSettingsVersionedValidation(t *testing.T) {
	for _, version := range []struct {
		name     string
		beta     string
		required bool
	}{
		{name: "v5", beta: "false"},
		{name: "v6", beta: "true", required: true},
	} {
		t.Run(version.name, func(t *testing.T) {
			t.Setenv("ARM_SIXPOINTZERO_BETA", version.beta)

			for _, resourceType := range []string{"azurerm_kubernetes_cluster", "azurerm_kubernetes_cluster_node_pool"} {
				t.Run(resourceType, func(t *testing.T) {
					resource := containers.Registration{}.SupportedResources()[resourceType]
					for _, test := range []struct {
						name        string
						settings    []interface{}
						expectError bool
					}{
						{name: "omitted", expectError: version.required},
						{name: "configured", settings: []interface{}{map[string]interface{}{"max_surge": "10%"}}},
						{
							name: "duplicate_blocks",
							settings: []interface{}{
								map[string]interface{}{"max_surge": "10%"},
								map[string]interface{}{"max_surge": "20%"},
							},
							expectError: true,
						},
					} {
						t.Run(test.name, func(t *testing.T) {
							nodePool := map[string]interface{}{
								"name":       "default",
								"vm_size":    "Standard_D2s_v3",
								"node_count": 1,
							}

							config := nodePool
							if resourceType == "azurerm_kubernetes_cluster" {
								config = map[string]interface{}{
									"name":                "acctestaks",
									"location":            "westus2",
									"resource_group_name": "acctestRG",
									"dns_prefix":          "acctestaks",
									"default_node_pool":   []interface{}{nodePool},
									"identity":            []interface{}{map[string]interface{}{"type": "SystemAssigned"}},
									"node_provisioning_profile": []interface{}{
										map[string]interface{}{"mode": "Manual", "default_node_pools": "Auto"},
									},
								}
							} else {
								config["kubernetes_cluster_id"] = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.ContainerService/managedClusters/acctestaks"
							}

							if test.settings != nil {
								nodePool["upgrade_settings"] = test.settings
							}
							diagnostics := resource.Validate(terraform.NewResourceConfigRaw(config))
							if diagnostics.HasError() != test.expectError {
								t.Fatalf("expected validation error %t, got %v", test.expectError, diagnostics)
							}
							if test.expectError && !strings.Contains(fmt.Sprint(diagnostics), "upgrade_settings") {
								t.Fatalf("expected an upgrade_settings validation error, got %v", diagnostics)
							}
						})
					}
				})
			}
		})
	}
}
