package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMKubernetesCluster_addAgent(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_addAgent(t)
}

func testAccAzureRMKubernetesCluster_addAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_count", "1"),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(data, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_count", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_removeAgent(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_removeAgent(t)
}

func testAccAzureRMKubernetesCluster_removeAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_count", "2"),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(data, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_count", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_autoScalingNodeCountUnset(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_autoScalingNodeCountUnset(t)
}

func testAccAzureRMKubernetesCluster_autoScalingNodeCountUnset(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleNodeCountUnsetConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.min_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.max_count", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.enable_auto_scaling", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(t)
}

func testAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.min_count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.max_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.enable_auto_scaling", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(t)
}

func testAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.min_count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.max_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.enable_auto_scaling", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.availability_zones.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.availability_zones.0", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.availability_zones.1", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKubernetesCluster_addAgentConfig(data acceptance.TestData, numberOfAgents int) string {
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

  default_node_pool {
    name       = "default"
    node_count = %d
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, numberOfAgents)
}

func testAccAzureRMKubernetesCluster_autoscaleNodeCountUnsetConfig(data acceptance.TestData) string {
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

  default_node_pool {
    name                = "default"
    enable_auto_scaling = true
    min_count           = 2
    max_count           = 4
    vm_size             = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(data acceptance.TestData) string {
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

  default_node_pool {
    name                = "pool1"
    min_count           = 1
    max_count           = 2
    enable_auto_scaling = true
    vm_size             = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(data acceptance.TestData) string {
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
  kubernetes_version  = "%s"

  default_node_pool {
    name                = "pool1"
    min_count           = 1
    max_count           = 2
    enable_auto_scaling = true
    vm_size             = "Standard_DS2_v2"
    availability_zones  = ["1", "2"]
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, olderKubernetesVersion)
}
