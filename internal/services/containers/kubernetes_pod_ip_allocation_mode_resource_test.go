// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKubernetesClusterNodePool_podIPAllocationMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.podIPAllocationModeConfig(data, "StaticBlock"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("pod_ip_allocation_mode").HasValue("StaticBlock"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_defaultNodePoolPodIPAllocationMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultNodePoolPodIPAllocationMode(data, "StaticBlock"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.pod_ip_allocation_mode").HasValue("StaticBlock"),
			),
		},
		data.ImportStep(),
	})
}

func (KubernetesClusterNodePoolResource) podIPAllocationModeConfig(data acceptance.TestData, allocationMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "defaultpod" {
  name                 = "acctestdefaultpodsubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]

  delegation {
    name = "aks-delegation"
    service_delegation {
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_subnet" "testpod" {
  name                 = "acctesttestpodsubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.3.0.0/16"]

  delegation {
    name = "aks-delegation"
    service_delegation {
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_D2s_v3"
    vnet_subnet_id = azurerm_subnet.test.id
    pod_subnet_id  = azurerm_subnet.defaultpod.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  node_provisioning_profile {
    mode               = "Manual"
    default_node_pools = "Auto"
  }

  network_profile {
    network_plugin = "azure"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                   = "internal"
  kubernetes_cluster_id  = azurerm_kubernetes_cluster.test.id
  vm_size                = "Standard_D2s_v3"
  node_count             = 1
  vnet_subnet_id         = azurerm_subnet.test.id
  pod_subnet_id          = azurerm_subnet.testpod.id
  pod_ip_allocation_mode = "%[3]s"
  upgrade_settings {
    max_surge = "10%%"
  }
}
`, data.RandomInteger, data.Locations.Primary, allocationMode)
}

func (KubernetesClusterResource) defaultNodePoolPodIPAllocationMode(data acceptance.TestData, mode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/8"]
}

resource "azurerm_subnet" "node" {
  name                 = "node-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "pod" {
  name                 = "pod-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]

  delegation {
    name = "aks-delegation"
    service_delegation {
      name = "Microsoft.ContainerService/managedClusters"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name                   = "default"
    node_count             = 1
    vm_size                = "Standard_D2s_v3"
    vnet_subnet_id         = azurerm_subnet.node.id
    pod_subnet_id          = azurerm_subnet.pod.id
    pod_ip_allocation_mode = "%[3]s"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  node_provisioning_profile {
    mode               = "Manual"
    default_node_pools = "Auto"
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "azure"
  }
}
`, data.RandomInteger, data.Locations.Primary, mode)
}
