package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMKubernetesCluster_addAgent(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	initConfig := testAccAzureRMKubernetesCluster_basic(ri, clientId, clientSecret, testLocation())
	addAgentConfig := testAccAzureRMKubernetesCluster_addAgent(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: initConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
				),
			},
			{
				Config: addAgentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.count", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	config := testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZones(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.min_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.max_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.enable_auto_scaling", "true"),
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

func TestAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	config := testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZones(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.min_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.max_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.enable_auto_scaling", "true"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.availability_zones.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.availability_zones.0", "1"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.availability_zones.1", "2"),
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

func TestAccAzureRMKubernetesCluster_multipleAgents(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	config := testAccAzureRMKubernetesCluster_multipleAgents(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.0.name", "pool1"),
					resource.TestCheckResourceAttr(resourceName, "agent_pool_profile.1.name", "pool2"),
				),
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_addAgent(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks%d"

  agent_pool_profile {
    name    = "default"
    count   = "2"
    vm_size = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZones(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks%d"

  agent_pool_profile {
    name                = "pool1"
    min_count           = "1"
    max_count           = "2"
    enable_auto_scaling = "true"
    type                = "VirtualMachineScaleSets"
    vm_size             = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZones(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  agent_pool_profile {
    name                = "pool1"
    min_count           = "1"
    max_count           = "2"
    enable_auto_scaling = "true"
    type                = "VirtualMachineScaleSets"
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

func testAccAzureRMKubernetesCluster_multipleAgents(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks%d"

  agent_pool_profile {
    name    = "pool1"
    count   = "1"
    vm_size = "Standard_DS2_v2"
  }

  agent_pool_profile {
    name    = "pool2"
    count   = "1"
    vm_size = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}
