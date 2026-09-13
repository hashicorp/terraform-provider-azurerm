// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/zclconf/go-cty/cty"
)

func TestKubernetesClusterOSSKUConfigurations(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 12345,
		Locations:     acceptance.Regions{Primary: "eastus"},
	}
	cluster := KubernetesClusterResource{}
	nodePool := KubernetesClusterNodePoolResource{}
	cases := []struct {
		name    string
		configs []string
	}{
		{
			name: "TestAccKubernetesClusterNodePool_osSkuUbuntu",
			configs: []string{
				nodePool.osSku(data, "Ubuntu"),
				nodePool.osSku(data, "Ubuntu2204"),
			},
		},
		{
			name: "TestAccKubernetesClusterNodePool_osSkuAzureLinux",
			configs: []string{
				nodePool.osSku(data, "AzureLinux"),
				nodePool.osSku(data, "AzureLinux3"),
			},
		},
		{
			name: "TestAccKubernetesClusterNodePool_osSkuMigration",
			configs: []string{
				nodePool.osSku(data, "Ubuntu"),
				nodePool.osSku(data, "AzureLinux"),
				nodePool.osSku(data, "Ubuntu"),
			},
		},
		{
			name:    "TestAccKubernetesClusterNodePool_osSkuAzureContainerLinux",
			configs: []string{nodePool.osSkuAzureContainerLinux(data)},
		},
		{
			name:    "TestAccKubernetesCluster_osSku",
			configs: []string{cluster.osSku(data, "AzureLinux")},
		},
		{
			name: "TestAccKubernetesCluster_osSkuUpdate",
			configs: []string{
				cluster.osSku(data, "Ubuntu"),
				cluster.osSku(data, "AzureLinux"),
				cluster.osSku(data, "Ubuntu2204"),
				cluster.osSku(data, "AzureLinux3"),
			},
		},
		{
			name:    "TestAccKubernetesCluster_osSkuAzureContainerLinux",
			configs: []string{cluster.osSkuAzureContainerLinux(data)},
		},
	}

	resource := containers.Registration{}.SupportedResources()["azurerm_kubernetes_cluster"]
	profile := resource.Schema["node_provisioning_profile"]
	if !profile.Required || profile.MaxItems != 1 {
		t.Fatal("expected node_provisioning_profile to be a required singleton block")
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			for i, config := range test.configs {
				t.Run(fmt.Sprintf("config_%d", i+1), func(t *testing.T) {
					file, diags := hclsyntax.ParseConfig([]byte(config), "config.tf", hcl.InitialPos)
					if diags.HasErrors() {
						t.Fatalf("parsing fixture: %s", diags)
					}
					resources, _, diags := file.Body.PartialContent(&hcl.BodySchema{
						Blocks: []hcl.BlockHeaderSchema{{Type: "resource", LabelNames: []string{"type", "name"}}},
					})
					if diags.HasErrors() {
						t.Fatalf("reading resources: %s", diags)
					}
					clusters := 0
					for _, block := range resources.Blocks {
						if block.Labels[0] == "azurerm_kubernetes_cluster_node_pool" {
							content, _, diags := block.Body.PartialContent(&hcl.BodySchema{
								Attributes: []hcl.AttributeSchema{{Name: "node_count", Required: true}},
							})
							if diags.HasErrors() {
								t.Fatalf("reading node_count: %s", diags)
							}
							nodeCount, diags := content.Attributes["node_count"].Expr.Value(nil)
							if diags.HasErrors() {
								t.Fatalf("decoding node_count: %s", diags)
							}
							if !nodeCount.RawEquals(cty.NumberIntVal(1)) {
								t.Errorf("OS SKU fixtures must provision one node to exercise its image, got %s", nodeCount.GoString())
							}
						}
						if block.Labels[0] != "azurerm_kubernetes_cluster" {
							continue
						}
						clusters++
						content, _, diags := block.Body.PartialContent(&hcl.BodySchema{
							Blocks: []hcl.BlockHeaderSchema{{Type: "node_provisioning_profile"}},
						})
						if diags.HasErrors() {
							t.Fatalf("reading node_provisioning_profile: %s", diags)
						}
						if got := len(content.Blocks); got != profile.MaxItems {
							t.Errorf("node_provisioning_profile block count = %d; required singleton schema allows %d", got, profile.MaxItems)
						}
					}
					if clusters != 1 {
						t.Errorf("expected one cluster resource, got %d", clusters)
					}
				})
			}
		})
	}
}

func TestKubernetesClusterAzureContainerLinuxSchema(t *testing.T) {
	resources := containers.Registration{}.SupportedResources()
	cluster := resources["azurerm_kubernetes_cluster"]
	nodePool := resources["azurerm_kubernetes_cluster_node_pool"]
	defaultNodePool := cluster.Schema["default_node_pool"].Elem.(*pluginsdk.Resource)

	for name, schema := range map[string]*pluginsdk.Schema{
		"default_node_pool": defaultNodePool.Schema["os_sku"],
		"node_pool":         nodePool.Schema["os_sku"],
	} {
		t.Run(name, func(t *testing.T) {
			_, errors := schema.ValidateFunc("AzureContainerLinux", "os_sku")
			if len(errors) != 0 {
				t.Fatalf("AzureContainerLinux rejected: %v", errors)
			}
		})
	}
}
