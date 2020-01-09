package tests

import (
	"fmt"
	"os"
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
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(data, clientId, clientSecret, data.Locations.Primary, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_count", "1"),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(data, clientId, clientSecret, data.Locations.Primary, 2),
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
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(data, clientId, clientSecret, data.Locations.Primary, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_count", "2"),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(data, clientId, clientSecret, data.Locations.Primary, 1),
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
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleNodeCountUnsetConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.min_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.max_count", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.enable_auto_scaling", "true"),
				),
			},
			data.ImportStep(
				"service_principal.0.client_secret",
			),
		},
	})
}

func TestAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(t)
}

func testAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.min_count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.max_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.enable_auto_scaling", "true"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(t)
}

func testAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(data, clientId, clientSecret),
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
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func testAccAzureRMKubernetesCluster_addAgentConfig(data acceptance.TestData, clientId, clientSecret, location string, numberOfAgents int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, numberOfAgents, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_autoscaleNodeCountUnsetConfig(data acceptance.TestData, clientId, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, olderKubernetesVersion, clientId, clientSecret)
}
