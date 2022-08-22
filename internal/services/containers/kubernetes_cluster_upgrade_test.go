package containers_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKubernetesCluster_upgradeAutoScaleMinCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeAutoScaleMinCountConfig(data, olderKubernetesVersion, 4, 8),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeAutoScaleMinCountConfig(data, olderKubernetesVersion, 5, 8),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_upgradeControlPlane(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeControlPlaneConfig(data, olderKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeControlPlaneConfig(data, currentKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// the control plane should have been upgraded but the default node pool shouldn't have been
				// TODO: confirm if we can roll the default node pool if the value is unset in the config
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTogether(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesVersion, olderKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesVersion, currentKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(currentKubernetesVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTwoPhase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesVersion, olderKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesVersion, olderKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesVersion, currentKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(currentKubernetesVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_upgradeNodePoolBeforeControlPlaneFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesVersion, olderKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
			),
		},
		data.ImportStep(),
		{
			Config:      r.upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesVersion, currentKubernetesVersion),
			ExpectError: regexp.MustCompile("Node Pools cannot use a version of Kubernetes that is not supported on the Control Plane."),
		},
	})
}

func TestAccKubernetesCluster_upgradeCustomNodePoolAfterControlPlane(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// all on the older version
			Config: r.upgradeVersionsConfig(data, olderKubernetesVersion, olderKubernetesVersion, olderKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
				acceptance.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesVersion),
			),
		},
		data.ImportStep(),
		{
			// upgrade the control plane
			Config: r.upgradeVersionsConfig(data, currentKubernetesVersion, olderKubernetesVersion, olderKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
				acceptance.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesVersion),
			),
		},
		data.ImportStep(),
		{
			// upgrade the node pool
			Config: r.upgradeVersionsConfig(data, currentKubernetesVersion, olderKubernetesVersion, currentKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
				acceptance.TestCheckResourceAttr(nodePoolName, "orchestrator_version", currentKubernetesVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_upgradeCustomNodePoolBeforeControlPlaneFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// all on the older version
			Config: r.upgradeVersionsConfig(data, olderKubernetesVersion, olderKubernetesVersion, olderKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesVersion),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersion),
				acceptance.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesVersion),
			),
		},
		data.ImportStep(),
		{
			// upgrade the node pool
			Config:      r.upgradeVersionsConfig(data, olderKubernetesVersion, olderKubernetesVersion, currentKubernetesVersion),
			ExpectError: regexp.MustCompile("Node Pools cannot use a version of Kubernetes that is not supported on the Control Plane."),
		},
	})
}

func TestAccKubernetesCluster_upgradeControlPlaneAndAllPoolsTogetherVersionAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// all on the older version
			Config: r.upgradeVersionsConfig(data, olderKubernetesVersionAlias, olderKubernetesVersionAlias, olderKubernetesVersionAlias),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesVersionAlias),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(olderKubernetesVersionAlias),
				check.That(nodePoolName).Key("orchestrator_version").HasValue(olderKubernetesVersionAlias),
			),
		},
		data.ImportStep(),
		{
			// upgrade control plane, default and custom node pools
			Config: r.upgradeVersionsConfig(data, currentKubernetesVersionAlias, currentKubernetesVersionAlias, currentKubernetesVersionAlias),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesVersionAlias),
				check.That(data.ResourceName).Key("default_node_pool.0.orchestrator_version").HasValue(currentKubernetesVersionAlias),
				check.That(nodePoolName).Key("orchestrator_version").HasValue(currentKubernetesVersionAlias),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_upgradeSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeSettingsConfig(data, "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.upgrade_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("default_node_pool.0.upgrade_settings.0.max_surge").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeSettingsConfig(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.upgrade_settings.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeSettingsConfig(data, "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.upgrade_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("default_node_pool.0.upgrade_settings.0.max_surge").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (KubernetesClusterResource) upgradeControlPlaneConfig(data acceptance.TestData, controlPlaneVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesClusterResource) upgradeControlPlaneDefaultNodePoolConfig(data acceptance.TestData, controlPlaneVersion, defaultNodePoolVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name                 = "default"
    node_count           = 1
    vm_size              = "Standard_DS2_v2"
    orchestrator_version = %q
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, defaultNodePoolVersion)
}

func (r KubernetesClusterResource) upgradeVersionsConfig(data acceptance.TestData, controlPlaneVersion, defaultNodePoolVersion, customNodePoolVersion string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  orchestrator_version  = %q
}
`, r.upgradeControlPlaneDefaultNodePoolConfig(data, controlPlaneVersion, defaultNodePoolVersion), customNodePoolVersion)
}

func (KubernetesClusterResource) upgradeAutoScaleMinCountConfig(data acceptance.TestData, controlPlaneVersion string, minCount int, maxCount int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  default_node_pool {
    name                = "default"
    vm_size             = "Standard_DS2_v2"
    enable_auto_scaling = true
    min_count           = %d
    max_count           = %d
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, minCount, maxCount)
}
