package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMKubernetesCluster_basicAvailabilitySet(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_basicAvailabilitySet(t)
}

func testAccAzureRMKubernetesCluster_basicAvailabilitySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_basicAvailabilitySetConfig(data, clientId, clientSecret),
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
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Basic"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_basicVMSS(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_basicVMSS(t)
}

func testAccAzureRMKubernetesCluster_basicVMSS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_basicVMSSConfig(data, clientId, clientSecret),
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
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Basic"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_requiresImport(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_requiresImport(t)
}

func testAccAzureRMKubernetesCluster_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_basicVMSSConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMKubernetesCluster_requiresImportConfig(data, clientId, clientSecret),
				ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_cluster"),
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
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_linuxProfileConfig(data, clientId, clientSecret),
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
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_nodeTaints(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_nodeTaints(t)
}

func testAccAzureRMKubernetesCluster_nodeTaints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_nodeTaintsConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.node_taints.0", "key=value:PreferNoSchedule"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_nodeResourceGroup(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_nodeResourceGroup(t)
}

func testAccAzureRMKubernetesCluster_nodeResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_nodeResourceGroupConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgrade(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgrade(t)
}

func testAccAzureRMKubernetesCluster_upgrade(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_upgradeConfig(data, clientId, clientSecret, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", olderKubernetesVersion),
				),
			},
			{
				Config: testAccAzureRMKubernetesCluster_upgradeConfig(data, clientId, clientSecret, currentKubernetesVersion),
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
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_tagsConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
			{
				Config: testAccAzureRMKubernetesCluster_tagsUpdatedConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccAzureRMKubernetesCluster_windowsProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_windowsProfile(t)
}

func testAccAzureRMKubernetesCluster_windowsProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_windowsProfileConfig(data, clientId, clientSecret),
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
				"service_principal.0.client_secret",
			),
		},
	})
}

func testAccAzureRMKubernetesCluster_basicAvailabilitySetConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_basicVMSSConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_requiresImportConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesCluster_basicVMSSConfig(data, clientId, clientSecret)
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

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, template, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_linuxProfileConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_nodeTaintsConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
    vm_size    = "Standard_DS2_v2"
    node_taints = [
      "key=value:PreferNoSchedule"
    ]
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_nodeResourceGroupConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
  node_resource_group = "acctestRGAKS-%d"

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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_tagsConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  tags = {
    dimension = "C-137"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_tagsUpdatedConfig(data acceptance.TestData, clientId string, clientSecret string) string {
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
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  tags = {
    dimension = "D-99"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_upgradeConfig(data acceptance.TestData, clientId, clientSecret, version string) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, version, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_windowsProfileConfig(data acceptance.TestData, clientId, clientSecret string) string {
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

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  network_profile {
    network_plugin     = "azure"
    network_policy     = "azure"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}
