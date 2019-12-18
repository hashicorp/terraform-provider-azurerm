package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func testAccAzureRMKubernetesCluster_addAgent(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(ri, clientId, clientSecret, location, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.node_count", "1"),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(ri, clientId, clientSecret, location, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.node_count", "2"),
				),
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_removeAgent(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(ri, clientId, clientSecret, location, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.node_count", "2"),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_addAgentConfig(ri, clientId, clientSecret, location, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.node_count", "1"),
				),
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_autoScalingNodeCountUnset(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleNodeCountUnsetConfig(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.min_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.max_count", "4"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.enable_auto_scaling", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.min_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.max_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.enable_auto_scaling", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.min_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.max_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.enable_auto_scaling", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.availability_zones.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.availability_zones.0", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_node_pool.0.availability_zones.1", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_principal.0.client_secret"},
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_addAgentConfig(rInt int, clientId, clientSecret, location string, numberOfAgents int) string {
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
`, rInt, location, rInt, rInt, numberOfAgents, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_autoscaleNodeCountUnsetConfig(rInt int, clientId, clientSecret, location string) string {
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
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(rInt int, clientId string, clientSecret string, location string) string {
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
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(rInt int, clientId string, clientSecret string, location string) string {
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
`, rInt, location, rInt, rInt, olderKubernetesVersion, clientId, clientSecret)
}
