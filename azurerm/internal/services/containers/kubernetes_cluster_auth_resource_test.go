package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

var kubernetesAuthTests = map[string]func(t *testing.T){
	"apiServerAuthorizedIPRanges": testAccKubernetesCluster_apiServerAuthorizedIPRanges,
	"managedClusterIdentity":      testAccKubernetesCluster_managedClusterIdentity,
	"userAssignedIdentity":        testAccKubernetesCluster_userAssignedIdentity,
	"roleBasedAccessControl":      testAccKubernetesCluster_roleBasedAccessControl,
	"AAD":                         testAccKubernetesCluster_roleBasedAccessControlAAD,
	"AADUpdateToManaged":          testAccKubernetesCluster_roleBasedAccessControlAADUpdateToManaged,
	"AADManaged":                  testAccKubernetesCluster_roleBasedAccessControlAADManaged,
	"AADManagedChange":            testAccKubernetesCluster_roleBasedAccessControlAADManagedChange,
	"servicePrincipal":            testAccKubernetesCluster_servicePrincipal,
}

func TestAccKubernetesCluster_apiServerAuthorizedIPRanges(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_apiServerAuthorizedIPRanges(t)
}

func testAccKubernetesCluster_apiServerAuthorizedIPRanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.apiServerAuthorizedIPRangesConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_config.0.client_key").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.client_certificate").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.cluster_ca_certificate").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.host").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.username").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.password").Exists(),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").HasValue(""),
				check.That(data.ResourceName).Key("default_node_pool.0.max_pods").Exists(),
				check.That(data.ResourceName).Key("api_server_authorized_ip_ranges.#").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_managedClusterIdentity(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_managedClusterIdentity(t)
}

func TestAccKubernetesCluster_userAssignedIdentity(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_userAssignedIdentity(t)
}

func testAccKubernetesCluster_managedClusterIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.managedClusterIdentityConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("kubelet_identity.0.client_id").Exists(),
				check.That(data.ResourceName).Key("kubelet_identity.0.object_id").Exists(),
				check.That(data.ResourceName).Key("kubelet_identity.0.user_assigned_identity_id").Exists(),
				check.That(data.ResourceName).Key("service_principal.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func testAccKubernetesCluster_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.userAssignedIdentityConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.user_assigned_identity_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_roleBasedAccessControl(t)
}

func testAccKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleBasedAccessControlConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_roleBasedAccessControlAAD(t)
}

func testAccKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	clientData := data.Client()
	auth := clientData.Default

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleBasedAccessControlAADConfig(data, auth.ClientID, auth.ClientSecret, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.client_app_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.server_app_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.server_app_secret").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
		{
			// should be no changes since the default for Tenant ID comes from the Provider block
			Config:   r.roleBasedAccessControlAADConfig(data, auth.ClientID, auth.ClientSecret, clientData.TenantID),
			PlanOnly: true,
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKubernetesCluster_roleBasedAccessControlAADUpdateToManaged(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_roleBasedAccessControlAADUpdateToManaged(t)
}

func testAccKubernetesCluster_roleBasedAccessControlAADUpdateToManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	clientData := data.Client()
	auth := clientData.Default

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleBasedAccessControlAADConfig(data, auth.ClientID, auth.ClientSecret, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.client_app_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.server_app_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.server_app_secret").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
		{
			Config: r.roleBasedAccessControlAADManagedConfig(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.managed").HasValue("true"),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
	})
}

func TestAccKubernetesCluster_roleBasedAccessControlAADManaged(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_roleBasedAccessControlAADManaged(t)
}

func testAccKubernetesCluster_roleBasedAccessControlAADManaged(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	clientData := data.Client()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleBasedAccessControlAADManagedConfig(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.managed").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.azure_rbac_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
		{
			// should be no changes since the default for Tenant ID comes from the Provider block
			Config:   r.roleBasedAccessControlAADManagedConfig(data, clientData.TenantID),
			PlanOnly: true,
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
	})
}

func TestAccKubernetesCluster_roleBasedAccessControlAADManagedChange(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_roleBasedAccessControlAADManagedChange(t)
}

func testAccKubernetesCluster_roleBasedAccessControlAADManagedChange(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	clientData := data.Client()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleBasedAccessControlAADManagedConfig(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.managed").Exists(),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("1"),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
		{
			Config: r.roleBasedAccessControlAADManagedConfigScale(data, clientData.TenantID),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
	})
}

func TestAccKubernetesCluster_roleBasedAccessControlAzure(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_roleBasedAccessControlAzure(t)
}

func testAccKubernetesCluster_roleBasedAccessControlAzure(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	clientData := data.Client()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleBasedAccessControlAzureConfig(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.managed").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.azure_rbac_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
		{
			// should be no changes since the default for Tenant ID comes from the Provider block
			Config:   r.roleBasedAccessControlAzureConfig(data, clientData.TenantID),
			PlanOnly: true,
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
	})
}

func TestAccKubernetesCluster_servicePrincipal(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_servicePrincipal(t)
}

func testAccKubernetesCluster_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}
	clientData := data.Client()

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.servicePrincipalConfig(data, clientData.Default.ClientID, clientData.Default.ClientSecret),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.%").HasValue("0"),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
		{
			Config: r.servicePrincipalConfig(data, clientData.Alternate.ClientID, clientData.Alternate.ClientSecret),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.%").HasValue("0"),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccKubernetesCluster_updateRoleBasedAccessControlAAD(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_updateRoleBaseAccessControlAAD(t)
}

func testAccKubernetesCluster_updateRoleBaseAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	clientData := data.Client()
	auth := clientData.Default
	altAlt := clientData.Alternate

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.roleBasedAccessControlAADConfig(data, auth.ClientID, auth.ClientSecret, clientData.TenantID),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.client_app_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.server_app_id").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.server_app_secret").Exists(),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
		{
			Config: r.updateRoleBasedAccessControlAADConfig(data, altAlt.ClientID, altAlt.ClientSecret, clientData.TenantID),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.client_app_id").HasValue(altAlt.ClientID),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.server_app_id").HasValue(altAlt.ClientID),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.server_app_secret").HasValue(altAlt.ClientSecret),
				check.That(data.ResourceName).Key("role_based_access_control.0.azure_active_directory.0.tenant_id").HasValue(clientData.TenantID),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
			),
		},
		data.ImportStep(
			"role_based_access_control.0.azure_active_directory.0.server_app_secret",
		),
	})
}

func (KubernetesClusterResource) apiServerAuthorizedIPRangesConfig(data acceptance.TestData) string {
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

func (KubernetesClusterResource) managedClusterIdentityConfig(data acceptance.TestData) string {
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

func (KubernetesClusterResource) userAssignedIdentityConfig(data acceptance.TestData) string {
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

func (KubernetesClusterResource) roleBasedAccessControlConfig(data acceptance.TestData) string {
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

func (KubernetesClusterResource) roleBasedAccessControlAADConfig(data acceptance.TestData, clientId, clientSecret, tenantId string) string {
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

func (KubernetesClusterResource) roleBasedAccessControlAADManagedConfig(data acceptance.TestData, tenantId string) string {
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

func (KubernetesClusterResource) roleBasedAccessControlAADManagedConfigScale(data acceptance.TestData, tenantId string) string {
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

func (KubernetesClusterResource) roleBasedAccessControlAzureConfig(data acceptance.TestData, tenantId string) string {
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
      tenant_id          = var.tenant_id
      managed            = true
      azure_rbac_enabled = true
    }
  }
}
`, tenantId, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) servicePrincipalConfig(data acceptance.TestData, clientId, clientSecret string) string {
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

func (KubernetesClusterResource) updateRoleBasedAccessControlAADConfig(data acceptance.TestData, altClientId, altClientSecret, tenantId string) string {
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
