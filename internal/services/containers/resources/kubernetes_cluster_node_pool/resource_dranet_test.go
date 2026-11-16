// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
)

func TestAccKubernetesClusterNodePool_dranet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}
	meta := helpers.GetAKSTestMetadata(t, data.Locations.Primary)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dranet(data, meta.CurrentKubernetesVersionAlias, true, `dranet_mode = "Managed"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_network_profile.0.dranet_mode").HasValue("Managed"),
			),
		},
		data.ImportStep(),

		{
			Config: r.dranet(data, meta.CurrentKubernetesVersionAlias, true, `dranet_mode = "Unmanaged"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_network_profile.0.dranet_mode").HasValue("Unmanaged"),
			),
		},
		data.ImportBlockWithIDStep(true),

		{
			Config: r.dranet(data, meta.CurrentKubernetesVersionAlias, true, `dranet_mode = "Managed"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_network_profile.0.dranet_mode").HasValue("Managed"),
			),
		},
		data.ImportStep(),

		{
			Config: r.dranet(data, meta.CurrentKubernetesVersionAlias, true, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_network_profile.0.dranet_mode").HasValue("Unmanaged"),
			),
		},
		data.ImportBlockWithIDStep(true),

		{
			Config: r.dranet(data, meta.CurrentKubernetesVersionAlias, true, `dranet_mode = "Managed"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_network_profile.0.dranet_mode").HasValue("Managed"),
			),
		},
		data.ImportStep(),

		{
			Config: r.dranet(data, meta.CurrentKubernetesVersionAlias, false, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_network_profile.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (KubernetesClusterNodePoolResource) dranet(data acceptance.TestData, kubernetesVersion string, includeNetworkProfile bool, settings string) string {
	networkProfile := ""
	if includeNetworkProfile {
		networkProfile = fmt.Sprintf(`
node_network_profile {
%s
}
`, settings)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dranet-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"
  kubernetes_version  = "%[4]s"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v4"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  node_provisioning_profile {
    mode = "Manual"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "dranet"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v4"
  node_count            = 1
  %[3]s
  upgrade_settings {
    max_surge = "10%%"
  }
}
`, data.RandomInteger, data.Locations.Primary, networkProfile, kubernetesVersion)
}
