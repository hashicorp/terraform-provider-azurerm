// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

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
