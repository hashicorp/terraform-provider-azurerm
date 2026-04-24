// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKubernetesAutomaticCluster_upgradeControlPlane(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeControlPlaneConfig(data, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeControlPlaneConfig(data, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// the control plane should have been upgraded but the default node pool shouldn't have been
				// TODO: confirm if we can roll the default node pool if the value is unset in the config
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeControlPlaneAndDefaultNodePoolTogether(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesAutomaticVersion, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(currentKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeControlPlaneAndDefaultNodePoolTwoPhase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesAutomaticVersion, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesAutomaticVersion, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(currentKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeNodePoolBeforeControlPlaneFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			Config:      r.upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesAutomaticVersion, currentKubernetesAutomaticVersion),
			ExpectError: regexp.MustCompile(fmt.Sprintf("Node pool version %s and control plane version %s are incompatible.", currentKubernetesAutomaticVersion, olderKubernetesAutomaticVersion)),
		},
	})
}

func TestAccKubernetesAutomaticCluster_upgradeCustomNodePoolAfterControlPlane(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// all on the older version
			Config: r.upgradeVersionsConfig(data, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
				acceptance.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			// upgrade the control plane
			Config: r.upgradeVersionsConfig(data, currentKubernetesAutomaticVersion, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
				acceptance.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			// upgrade the node pool
			Config: r.upgradeVersionsConfig(data, currentKubernetesAutomaticVersion, olderKubernetesAutomaticVersion, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
				acceptance.TestCheckResourceAttr(nodePoolName, "orchestrator_version", currentKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeCustomNodePoolBeforeControlPlaneFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// all on the older version
			Config: r.upgradeVersionsConfig(data, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
				acceptance.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			// upgrade the node pool
			Config:      r.upgradeVersionsConfig(data, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion, currentKubernetesAutomaticVersion),
			ExpectError: regexp.MustCompile("Node Pools cannot use a version of Kubernetes that is not supported on the Control Plane."),
		},
	})
}

func TestAccKubernetesAutomaticCluster_upgradeControlPlaneAndAllPoolsTogetherVersionAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// all on the older version
			Config: r.upgradeVersionsConfig(data, olderKubernetesAutomaticVersionAlias, olderKubernetesAutomaticVersionAlias, olderKubernetesAutomaticVersionAlias),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersionAlias),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersionAlias),
				check.That(nodePoolName).Key("orchestrator_version").HasValue(olderKubernetesAutomaticVersionAlias),
			),
		},
		data.ImportStep(),
		{
			// upgrade control plane, default and custom node pools
			Config: r.upgradeVersionsConfig(data, currentKubernetesAutomaticVersionAlias, currentKubernetesAutomaticVersionAlias, currentKubernetesAutomaticVersionAlias),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesAutomaticVersionAlias),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(currentKubernetesAutomaticVersionAlias),
				check.That(nodePoolName).Key("orchestrator_version").HasValue(currentKubernetesAutomaticVersionAlias),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeControlPlaneAndAllPoolsTogetherSpot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// all on the older version
			Config: r.upgradeVersionsConfigSpot(data, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(nodePoolName).Key("orchestrator_version").HasValue(olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
		{
			// upgrade control plane, default and custom node pools
			Config: r.upgradeVersionsConfigSpot(data, currentKubernetesAutomaticVersion, currentKubernetesAutomaticVersion, currentKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(currentKubernetesAutomaticVersion),
				check.That(nodePoolName).Key("orchestrator_version").HasValue(currentKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeSettings(data, 10, 5, "Cordon"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeSettings(data, 15, 10, "Schedule"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r KubernetesAutomaticClusterResource) upgradeSettings(data acceptance.TestData, drainTimeout int, nodeSoak int, undrainableBehavior string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS3_v2"
    upgrade_settings {
      max_surge                     = "10%%"
      drain_timeout_in_minutes      = %d
      node_soak_duration_in_minutes = %d
      undrainable_node_behavior     = %q
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, drainTimeout, nodeSoak, undrainableBehavior)
}

func (KubernetesAutomaticClusterResource) upgradeControlPlaneConfig(data acceptance.TestData, controlPlaneVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS3_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesAutomaticClusterResource) upgradeControlPlaneDefaultNodePoolConfig(data acceptance.TestData, controlPlaneVersion, defaultNodePoolVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name                 = "default"
    node_count           = 1
    vm_size              = "Standard_DS3_v2"
    orchestrator_version = %q
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, defaultNodePoolVersion)
}

func (r KubernetesAutomaticClusterResource) upgradeVersionsConfig(data acceptance.TestData, controlPlaneVersion, defaultNodePoolVersion, customNodePoolVersion string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_automatic_cluster.test.id
  vm_size               = "Standard_DS3_v2"
  node_count            = 1
  orchestrator_version  = %q
  upgrade_settings {
    max_surge = "10%%"
  }
  depends_on = [azurerm_kubernetes_automatic_cluster.test]
}
`, r.upgradeControlPlaneDefaultNodePoolConfig(data, controlPlaneVersion, defaultNodePoolVersion), customNodePoolVersion)
}

func (r KubernetesAutomaticClusterResource) upgradeVersionsConfigSpot(data acceptance.TestData, controlPlaneVersion, defaultNodePoolVersion, customNodePoolVersion string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_automatic_cluster.test.id
  vm_size               = "Standard_DS3_v2"
  node_count            = 1
  orchestrator_version  = %q
  priority              = "Spot"
  eviction_policy       = "Delete"
  spot_max_price        = 0.5 # high, but this is a maximum (we pay less) so ensures this won't fail
  node_labels = {
    "kubernetes.azure.com/scalesetpriority" = "spot"
  }
  node_taints = [
    "kubernetes.azure.com/scalesetpriority=spot:NoSchedule"
  ]
  depends_on = [azurerm_kubernetes_automatic_cluster.test]
}
`, r.upgradeControlPlaneDefaultNodePoolConfig(data, controlPlaneVersion, defaultNodePoolVersion), customNodePoolVersion)
}
