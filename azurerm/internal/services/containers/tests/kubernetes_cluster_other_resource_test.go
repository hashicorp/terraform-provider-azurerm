package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var kubernetesOtherTests = map[string]func(t *testing.T){
	"basicAvailabilitySet":     testAccAzureRMKubernetesCluster_basicAvailabilitySet,
	"basicVMSS":                testAccAzureRMKubernetesCluster_basicVMSS,
	"requiresImport":           testAccAzureRMKubernetesCluster_requiresImport,
	"podIdentityProfile":       testAccAzureRMKubernetesCluster_podIdentityProfile,
	"linuxProfile":             testAccAzureRMKubernetesCluster_linuxProfile,
	"nodeLabels":               testAccAzureRMKubernetesCluster_nodeLabels,
	"nodeResourceGroup":        testAccAzureRMKubernetesCluster_nodeResourceGroup,
	"paidSku":                  testAccAzureRMKubernetesCluster_paidSku,
	"upgradeConfig":            testAccAzureRMKubernetesCluster_upgrade,
	"tags":                     testAccAzureRMKubernetesCluster_tags,
	"windowsProfile":           testAccAzureRMKubernetesCluster_windowsProfile,
	"outboundTypeLoadBalancer": testAccAzureRMKubernetesCluster_outboundTypeLoadBalancer,
	"privateClusterOn":         testAccAzureRMKubernetesCluster_privateClusterOn,
	"privateClusterOff":        testAccAzureRMKubernetesCluster_privateClusterOff,
}

func TestAccAzureRMKubernetesCluster_basicAvailabilitySet(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_basicAvailabilitySet(t)
}

func TestAccAzureRMKubernetesCluster_sameSize(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_sameSizeVMSSConfig(t)
}

func testAccAzureRMKubernetesCluster_basicAvailabilitySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_basicAvailabilitySetConfig(data),
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
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKubernetesCluster_sameSizeVMSSConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_sameSize(data),
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
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_basicVMSS(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_basicVMSS(t)
}

func testAccAzureRMKubernetesCluster_basicVMSS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_basicVMSSConfig(data),
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
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_requiresImport(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_requiresImport(t)
}

func testAccAzureRMKubernetesCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_basicVMSSConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMKubernetesCluster_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_cluster"),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_podIdentityProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_podIdentityProfile(t)
}

func testAccAzureRMKubernetesCluster_podIdentityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	clientData := data.Client()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_podIdentityProfileConfig(data, clientData.TenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_podIdentityProfileConfigUpdate(data, clientData.TenantID),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.0.enabled", "true"),

					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.0.user_assigned_identities.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.0.user_assigned_identities.0.name", "my-identity"),
					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.0.user_assigned_identities.0.namespace", "default"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "pod_identity_profile.0.user_assigned_identities.0.identity.0.resource_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "pod_identity_profile.0.user_assigned_identities.0.identity.0.client_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "pod_identity_profile.0.user_assigned_identities.0.identity.0.object_id"),

					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.0.user_assigned_identity_exceptions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.0.user_assigned_identity_exceptions.0.name", "my-exception"),
					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.0.user_assigned_identity_exceptions.0.namespace", "default"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "pod_identity_profile.0.user_assigned_identity_exceptions.0.pod_labels.myPod"),
					resource.TestCheckResourceAttr(data.ResourceName, "pod_identity_profile.0.user_assigned_identity_exceptions.0.pod_labels.myPod", "isAwesome"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_linuxProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_linuxProfile(t)
}

func testAccAzureRMKubernetesCluster_linuxProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_linuxProfileConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "linux_profile.0.admin_username"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_nodeLabels(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_nodeLabels(t)
}

func testAccAzureRMKubernetesCluster_nodeLabels(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	labels1 := map[string]string{"key": "value"}
	labels2 := map[string]string{"key2": "value2"}
	labels3 := map[string]string{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_nodeLabelsConfig(data, labels1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_labels.key", "value"),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_nodeLabelsConfig(data, labels2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_labels.key2", "value2"),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_nodeLabelsConfig(data, labels3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckNoResourceAttr(data.ResourceName, "default_node_pool.0.node_labels"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_nodeResourceGroup(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_nodeResourceGroup(t)
}

func testAccAzureRMKubernetesCluster_nodeResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_nodeResourceGroupConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeSkuTier(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgradeSkuTier(t)
}

func testAccAzureRMKubernetesCluster_upgradeSkuTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_freeSkuConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_paidSkuConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_paidSku(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_paidSku(t)
}

func testAccAzureRMKubernetesCluster_paidSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_paidSkuConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgrade(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgrade(t)
}

func testAccAzureRMKubernetesCluster_upgrade(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_upgradeConfig(data, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", olderKubernetesVersion),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_upgradeConfig(data, currentKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", currentKubernetesVersion),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_tags(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_tags(t)
}

func testAccAzureRMKubernetesCluster_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_tagsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_tagsUpdatedConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_windowsProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_windowsProfile(t)
}

func testAccAzureRMKubernetesCluster_windowsProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_windowsProfileConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_node_pool.0.max_pods"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "linux_profile.0.admin_username"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "windows_profile.0.admin_username"),
				),
			},
			data.ImportStep(
				"windows_profile.0.admin_password",
			),
		},
	})
}

func TestAccAzureRMKubernetesCluster_diskEncryption(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_diskEncryption(t)
}

func testAccAzureRMKubernetesCluster_diskEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_diskEncryptionConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "disk_encryption_set_id"),
				),
			},
			data.ImportStep(
				"windows_profile.0.admin_password",
			),
		},
	})
}

func testAccAzureRMKubernetesCluster_basicAvailabilitySetConfig(data acceptance.TestData) string {
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
    type       = "AvailabilitySet"
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_sameSize(data acceptance.TestData) string {
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
    node_count          = 1
    enable_auto_scaling = true
    vm_size             = "Standard_DS2_v2"
    min_count           = 1
    max_count           = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_autoScalingEnabled(data acceptance.TestData) string {
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
    node_count          = 2
    vm_size             = "Standard_DS2_v2"
    enable_auto_scaling = true
    max_count           = 10
    min_count           = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_autoScalingEnabledUpdate(data acceptance.TestData) string {
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
    node_count          = 1
    vm_size             = "Standard_DS2_v2"
    enable_auto_scaling = true
    max_count           = 10
    min_count           = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_autoScalingEnabledUpdateMax(data acceptance.TestData) string {
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
    node_count          = 11
    vm_size             = "Standard_DS2_v2"
    enable_auto_scaling = true
    max_count           = 10
    min_count           = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_autoScalingWithMaxCountConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-AKS-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestAKS%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestAKS%d"

  default_node_pool {
    name                = "default"
    vm_size             = "Standard_DS2_v2"
    enable_auto_scaling = true
    min_count           = 1
    max_count           = 1000
    node_count          = 1
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_autoScalingEnabledUpdateMin(data acceptance.TestData) string {
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
    node_count          = 1
    vm_size             = "Standard_DS2_v2"
    enable_auto_scaling = true
    max_count           = 10
    min_count           = 2
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_basicVMSSConfig(data acceptance.TestData) string {
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

func testAccAzureRMKubernetesCluster_requiresImportConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesCluster_basicVMSSConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster" "import" {
  name                = azurerm_kubernetes_cluster.test.name
  location            = azurerm_kubernetes_cluster.test.location
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
  dns_prefix          = azurerm_kubernetes_cluster.test.dns_prefix

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template)
}

func testAccAzureRMKubernetesCluster_podIdentityProfileConfig(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctestaks%d"
}

resource "azurerm_role_assignment" "test" {
  depends_on = [
    azurerm_user_assigned_identity.test,
    azurerm_kubernetes_cluster.test
  ]

  scope                = azurerm_user_assigned_identity.test.id
  role_definition_name = "Managed Identity Operator"
  principal_id         = azurerm_kubernetes_cluster.test.identity[0].principal_id
}

resource "azurerm_kubernetes_cluster" "test" {
  depends_on = [azurerm_user_assigned_identity.test]

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

  # cannot use kubenet
  network_profile {
    network_plugin = "azure"
  }

  role_based_access_control {
    enabled = true

    azure_active_directory {
      tenant_id = "%s"
      managed   = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMKubernetesCluster_podIdentityProfileConfigUpdate(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  name = "acctestaks%d"
}

resource "azurerm_role_assignment" "test" {
  depends_on = [
    azurerm_user_assigned_identity.test,
    azurerm_kubernetes_cluster.test
  ]

  scope                = azurerm_user_assigned_identity.test.id
  role_definition_name = "Managed Identity Operator"
  principal_id         = azurerm_kubernetes_cluster.test.identity[0].principal_id
}

resource "azurerm_kubernetes_cluster" "test" {
  depends_on = [azurerm_user_assigned_identity.test]

  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  pod_identity_profile {
    enabled = true
    user_assigned_identities {
      name      = "my-identity"
      namespace = "default"
      identity {
        resource_id = azurerm_user_assigned_identity.test.id
        client_id   = azurerm_user_assigned_identity.test.client_id
        object_id   = azurerm_user_assigned_identity.test.principal_id
      }
    }

    user_assigned_identity_exceptions {
      name      = "my-exception"
      namespace = "default"
      pod_labels = {
        "myPod" = "isAwesome"
      }
    }
  }

  identity {
    type = "SystemAssigned"
  }

  # cannot use kubenet
  network_profile {
    network_plugin = "azure"
  }

  role_based_access_control {
    enabled = true

    azure_active_directory {
      tenant_id = "%s"
      managed   = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, tenantID)
}

func testAccAzureRMKubernetesCluster_linuxProfileConfig(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_nodeLabelsConfig(data acceptance.TestData, labels map[string]string) string {
	labelsSlice := make([]string, 0, len(labels))
	for k, v := range labels {
		labelsSlice = append(labelsSlice, fmt.Sprintf("      \"%s\" = \"%s\"", k, v))
	}
	labelsStr := strings.Join(labelsSlice, "\n")
	return fmt.Sprintf(`
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
    node_labels = {
%s
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, labelsStr)
}

func testAccAzureRMKubernetesCluster_nodeResourceGroupConfig(data acceptance.TestData) string {
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
  node_resource_group = "acctestRGAKS-%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_paidSkuConfig(data acceptance.TestData) string {
	// @tombuildsstuff (2020-05-29) - this is only supported in a handful of regions
	// 								  whilst in Preview - hard-coding for now
	location := "westus2" // TODO: data.Locations.Primary
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
  sku_tier            = "Paid"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, location, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_freeSkuConfig(data acceptance.TestData) string {
	// @tombuildsstuff (2020-05-29) - this is only supported in a handful of regions
	// 								  whilst in Preview - hard-coding for now
	location := "westus2" // TODO: data.Locations.Primary
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
`, data.RandomInteger, location, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_tagsConfig(data acceptance.TestData) string {
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

  tags = {
    dimension = "C-137"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_tagsUpdatedConfig(data acceptance.TestData) string {
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

  tags = {
    dimension = "D-99"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_upgradeConfig(data acceptance.TestData, version string) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, version, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_windowsProfileConfig(data acceptance.TestData) string {
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

  windows_profile {
    admin_username = "azureuser"
    admin_password = "P@55W0rd1234!"
  }

  # the default node pool /has/ to be Linux agents - Windows agents can be added via the node pools resource
  default_node_pool {
    name       = "np"
    node_count = 3
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin     = "azure"
    network_policy     = "azure"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_diskEncryptionConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkeyvault%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  enabled_for_disk_encryption = true
  soft_delete_enabled         = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "acctest" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "get",
    "create",
    "delete",
    "purge",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "destestkey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [azurerm_key_vault_access_policy.acctest]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption-perm" {
  key_vault_id = azurerm_key_vault.test.id

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id

  key_permissions = [
    "get",
    "wrapkey",
    "unwrapkey",
  ]
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_kubernetes_cluster" "test" {
  name                   = "acctestaks%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  dns_prefix             = "acctestaks%d"
  disk_encryption_set_id = azurerm_disk_encryption_set.test.id

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "np"
    node_count = 3
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin     = "azure"
    network_policy     = "azure"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }

  depends_on = [
    azurerm_key_vault_access_policy.disk-encryption-perm,
    azurerm_role_assignment.disk-encryption-read-keyvault
  ]

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
