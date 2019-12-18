package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRangesConfig(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config_raw", ""),
					resource.TestCheckResourceAttrSet(resourceName, "default_node_pool.0.max_pods"),
					resource.TestCheckResourceAttr(resourceName, "api_server_authorized_ip_ranges.#", "3"),
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

func testAccAzureRMKubernetesCluster_enablePodSecurityPolicy(t *testing.T) {
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
				Config: testAccAzureRMKubernetesCluster_enablePodSecurityPolicyConfig(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enable_pod_security_policy", "true"),
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

func testAccAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(ri, location, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config_raw", ""),
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

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantId := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(ri, location, clientId, clientSecret, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "role_based_access_control.0.azure_active_directory.0.client_app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_based_access_control.0.azure_active_directory.0.server_app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_based_access_control.0.azure_active_directory.0.server_app_secret"),
					resource.TestCheckResourceAttrSet(resourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_admin_config_raw"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"service_principal.0.client_secret",
					"role_based_access_control.0.azure_active_directory.0.server_app_secret",
				},
			},
			{
				// should be no changes since the default for Tenant ID comes from the Provider block
				Config:   testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(ri, location, clientId, clientSecret, tenantId),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"service_principal.0.client_secret",
					"role_based_access_control.0.azure_active_directory.0.server_app_secret",
				},
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRangesConfig(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "Standard"
  }

  api_server_authorized_ip_ranges = [
    "8.8.8.8/32",
    "8.8.4.4/32",
    "8.8.2.0/24",
  ]
}
`, rInt, location, rInt, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_enablePodSecurityPolicyConfig(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                       = "acctestaks%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  dns_prefix                 = "acctestaks%d"
  enable_pod_security_policy = true

  role_based_access_control {
    enabled = true
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(rInt int, location, clientId, clientSecret string) string {
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

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  role_based_access_control {
    enabled = true
  }
}
`, rInt, location, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(rInt int, location, clientId, clientSecret, tenantId string) string {
	return fmt.Sprintf(`
variable "tenant_id" {
  default = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  role_based_access_control {
    enabled = true

    azure_active_directory {
      server_app_id     = "%s"
      server_app_secret = "%s"
      client_app_id     = "%s"
      tenant_id         = var.tenant_id
    }
  }
}
`, tenantId, rInt, location, rInt, rInt, rInt, clientId, clientSecret, clientId, clientSecret, clientId)
}
