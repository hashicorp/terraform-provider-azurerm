// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKubernetesCluster_defaultNodePoolArtifactStreaming(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultNodePoolArtifactStreaming(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.artifact_streaming_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultNodePoolArtifactStreaming(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.artifact_streaming_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (KubernetesClusterResource) defaultNodePoolArtifactStreaming(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name                       = "default"
    vm_size                    = "Standard_D2s_v3"
    node_count                 = 1
    os_sku                     = "AzureLinux"
    artifact_streaming_enabled = %[3]t

    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  node_provisioning_profile {
    mode               = "Manual"
    default_node_pools = "Auto"
  }
}
`, data.RandomInteger, data.Locations.Primary, enabled)
}
