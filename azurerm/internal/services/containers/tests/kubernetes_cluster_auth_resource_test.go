package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var kubernetesAuthTests = map[string]func(t *testing.T){
	"apiServerAuthorizedIPRanges": testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges,
	"managedClusterIdentity":      testAccAzureRMKubernetesCluster_managedClusterIdentity,
	"userAssignedIdentity":        testAccAzureRMKubernetesCluster_userAssignedIdentity,
	"roleBasedAccessControl":      testAccAzureRMKubernetesCluster_roleBasedAccessControl,
	"AAD":                         testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD,
	"AADUpdateToManaged":          testAccAzureRMKubernetesCluster_roleBasedAccessControlAADUpdateToManaged,
	"AADManaged":                  testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManaged,
	"AADManagedChange":            testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedChange,
	"servicePrincipal":            testAccAzureRMKubernetesCluster_servicePrincipal,
}

func TestAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges(t)
}

func testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRangesConfig(data),
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
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_managedClusterIdentity(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_managedClusterIdentity(t)
}

func TestAccAzureRMKubernetesCluster_userAssignedIdentity(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_userAssignedIdentity(t)
}

func testAccAzureRMKubernetesCluster_managedClusterIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_managedClusterIdentityConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kubelet_identity.0.client_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kubelet_identity.0.object_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kubelet_identity.0.user_assigned_identity_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "service_principal.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKubernetesCluster_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_userAssignedIdentityConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "UserAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.user_assigned_identity_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_roleBasedAccessControl(t)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config_raw", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(t)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientData := data.Client()
	auth := clientData.Default

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, auth.ClientID, auth.ClientSecret, ""),
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
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
			{
				// should be no changes since the default for Tenant ID comes from the Provider block
				Config:   testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, auth.ClientID, auth.ClientSecret, clientData.TenantID),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_roleBasedAccessControlAADUpdateToManaged(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_roleBasedAccessControlAADUpdateToManaged(t)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAADUpdateToManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientData := data.Client()
	auth := clientData.Default

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, auth.ClientID, auth.ClientSecret, ""),
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
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedConfig(data, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.managed", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_admin_config_raw"),
				),
			},
			data.ImportStep(
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
		},
	})
}

func TestAccAzureRMKubernetesCluster_roleBasedAccessControlAADManaged(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManaged(t)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientData := data.Client()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedConfig(data, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.managed"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_admin_config_raw"),
				),
			},
			data.ImportStep(
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
			{
				// should be no changes since the default for Tenant ID comes from the Provider block
				Config:   testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedConfig(data, clientData.TenantID),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
		},
	})
}

func TestAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedChange(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedChange(t)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedChange(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientData := data.Client()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedConfig(data, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.managed"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_admin_config_raw"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_count", "1"),
				),
			},
			data.ImportStep(
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedConfigScale(data, clientData.TenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_count", "2"),
				),
			},
			data.ImportStep(
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
		},
	})
}

func TestAccAzureRMKubernetesCluster_servicePrincipal(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_servicePrincipal(t)
}

func testAccAzureRMKubernetesCluster_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientData := data.Client()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_servicePrincipalConfig(data, clientData.Default.ClientID, clientData.Default.ClientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.%", "0"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
			{
				Config: testAccAzureRMKubernetesCluster_servicePrincipalConfig(data, clientData.Alternate.ClientID, clientData.Alternate.ClientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.%", "0"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_updateRoleBasedAccessControlAAD(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_updateRoleBaseAccessControlAAD(t)
}

func testAccAzureRMKubernetesCluster_updateRoleBaseAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	clientData := data.Client()
	auth := clientData.Default
	altAlt := clientData.Alternate

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, auth.ClientID, auth.ClientSecret, clientData.TenantID),
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
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
			{
				Config: testAccAzureRMKubernetesCluster_updateRoleBasedAccessControlAADConfig(data, altAlt.ClientID, altAlt.ClientSecret, clientData.TenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.client_app_id", altAlt.ClientID),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_id", altAlt.ClientID),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_secret", altAlt.ClientSecret),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id", clientData.TenantID),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_admin_config_raw"),
				),
			},
			data.ImportStep(
				"role_based_access_control.0.azure_active_directory.0.server_app_secret",
			),
		},
	})
}

func testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRangesConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
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

  identity {
    type = "SystemAssigned"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_managedClusterIdentityConfig(data acceptance.TestData) string {
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
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_userAssignedIdentityConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type                      = "UserAssigned"
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(data acceptance.TestData) string {
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

  identity {
    type = "SystemAssigned"
  }

  role_based_access_control {
    enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data acceptance.TestData, clientId, clientSecret, tenantId string) string {
	return fmt.Sprintf(`
variable "tenant_id" {
  default = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
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

  identity {
    type = "SystemAssigned"
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
`, tenantId, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret, clientId)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedConfig(data acceptance.TestData, tenantId string) string {
	return fmt.Sprintf(`
variable "tenant_id" {
  default = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
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

  identity {
    type = "SystemAssigned"
  }

  role_based_access_control {
    enabled = true

    azure_active_directory {
      tenant_id = var.tenant_id
      managed   = true
    }
  }
}
`, tenantId, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_roleBasedAccessControlAADManagedConfigScale(data acceptance.TestData, tenantId string) string {
	return fmt.Sprintf(`
variable "tenant_id" {
  default = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
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
    node_count = 2
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  role_based_access_control {
    enabled = true

    azure_active_directory {
      tenant_id = var.tenant_id
      managed   = true
    }
  }
}
`, tenantId, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_servicePrincipalConfig(data acceptance.TestData, clientId, clientSecret string) string {
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

func testAccAzureRMKubernetesCluster_updateRoleBasedAccessControlAADConfig(data acceptance.TestData, altClientId, altClientSecret, tenantId string) string {
	return fmt.Sprintf(`
variable "tenant_id" {
  default = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
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

  identity {
    type = "SystemAssigned"
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
`, tenantId, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, altClientId, altClientSecret, altClientId)
}
