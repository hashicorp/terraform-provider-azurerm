package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges(t)
}

func testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRangesConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config_raw", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_node_pool.0.max_pods"),
					resource.TestCheckResourceAttr(data.ResourceName, "api_server_authorized_ip_ranges.#", "3"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_enablePodSecurityPolicy(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_enablePodSecurityPolicy(t)
}

func testAccAzureRMKubernetesCluster_enablePodSecurityPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_enablePodSecurityPolicyConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_pod_security_policy", "true"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_managedClusterIdentityServicePrincipal(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_managedClusterIdentityServicePrincipal(t)
}

func testAccAzureRMKubernetesCluster_managedClusterIdentityServicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_managedClusterIdentityServicePrincipalConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_roleBasedAccessControl(t)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config_raw", ""),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(t)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantId := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, clientId, clientSecret, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.client_app_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_secret"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_admin_config_raw"),
				),
			},
			data.ImportStep(
				"service_principal.0.client_secret",
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
			{
				// should be no changes since the default for Tenant ID comes from the Provider block
				Config:   testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, clientId, clientSecret, tenantId),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"service_principal.0.client_secret",
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
		},
	})
}

func TestAccAzureRMKubernetesCluster_updateRoleBasedAccessControlAAD(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_updateRoleBaseAccessControlAAD(t)
}

func testAccAzureRMKubernetesCluster_updateRoleBaseAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantId := os.Getenv("ARM_TENANT_ID")
	// TODO: find or create a suitable replacement client_id to use to extend the test and set ARM_CLIENT_ID_ALT in the CI job
	updateClientId := os.Getenv("ARM_CLIENT_ID_ALT")
	updateClientSecret := os.Getenv("ARM_CLIENT_SECRET_ALT")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, clientId, clientSecret, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.client_app_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_secret"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_admin_config_raw"),
				),
			},
			data.ImportStep(
				"service_principal.0.client_secret",
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
			{
				Config: testAccAzureRMKubernetesCluster_updateRoleBasedAccessControlAADConfig(data, clientId, clientSecret, updateClientId, updateClientSecret, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.client_app_id", updateClientId),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_id", updateClientId),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_secret", updateClientSecret),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id", tenantId),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_admin_config_raw"),
				),
			},
			data.ImportStep(
				"service_principal.0.client_secret",
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
		},
	})
}

func testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRangesConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_enablePodSecurityPolicyConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_managedClusterIdentityServicePrincipalConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
    node_count = 1
    type       = "AvailabilitySet"
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(data acceptance.TestData, clientId, clientSecret string) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data acceptance.TestData, clientId, clientSecret, tenantId string) string {
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
`, tenantId, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret, clientId, clientSecret, clientId)
}
func testAccAzureRMKubernetesCluster_updateRoleBasedAccessControlAADConfig(data acceptance.TestData, clientId, clientSecret, altClientId, altClientSecret, tenantId string) string {
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
`, tenantId, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret, altClientId, altClientSecret, altClientId)
}
