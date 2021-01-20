package containers_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

var kubernetesScalingTests = map[string]func(t *testing.T){
	"addAgent":                         testAccKubernetesCluster_addAgent,
	"manualScaleIgnoreChanges":         testAccKubernetesCluster_manualScaleIgnoreChanges,
	"removeAgent":                      testAccKubernetesCluster_removeAgent,
	"autoScalingEnabledError":          testAccKubernetesCluster_autoScalingError,
	"autoScalingEnabledErrorMax":       testAccKubernetesCluster_autoScalingErrorMax,
	"autoScalingEnabledErrorMin":       testAccKubernetesCluster_autoScalingErrorMin,
	"autoScalingEnabledWithMaxCount":   testAccKubernetesCluster_autoScalingWithMaxCount,
	"autoScalingNodeCountUnset":        testAccKubernetesCluster_autoScalingNodeCountUnset,
	"autoScalingNoAvailabilityZones":   testAccKubernetesCluster_autoScalingNoAvailabilityZones,
	"autoScalingWithAvailabilityZones": testAccKubernetesCluster_autoScalingWithAvailabilityZones,
}

func TestAccKubernetesCluster_addAgent(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_addAgent(t)
}

func testAccKubernetesCluster_addAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.addAgentConfig(data, 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("1"),
			),
		},
		{
			Config: r.addAgentConfig(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
	})
}

func TestAccKubernetesCluster_manualScaleIgnoreChanges(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_manualScaleIgnoreChanges(t)
}

func testAccKubernetesCluster_manualScaleIgnoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.manualScaleIgnoreChangesConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("1"),
				kubernetesClusterUpdateNodePoolCount(data.ResourceName, 2),
			),
		},
		{
			Config: r.manualScaleIgnoreChangesConfigUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
	})
}

func TestAccKubernetesCluster_removeAgent(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_removeAgent(t)
}

func testAccKubernetesCluster_removeAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.addAgentConfig(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
		{
			Config: r.addAgentConfig(data, 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("1"),
			),
		},
	})
}

func TestAccKubernetesCluster_autoScalingError(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_autoScalingError(t)
}

func testAccKubernetesCluster_autoScalingError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoScalingEnabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
		{
			Config: r.autoScalingEnabledUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("cannot change `node_count` when `enable_auto_scaling` is set to `true`"),
		},
	})
}

func TestAccKubernetesCluster_autoScalingErrorMax(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_autoScalingErrorMax(t)
}

func testAccKubernetesCluster_autoScalingErrorMax(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoScalingEnabledUpdateMax(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("`node_count`\\(11\\) must be equal to or less than `max_count`\\(10\\) when `enable_auto_scaling` is set to `true`"),
		},
	})
}

func TestAccKubernetesCluster_autoScalingWithMaxCount(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_autoScalingWithMaxCount(t)
}

func testAccKubernetesCluster_autoScalingWithMaxCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoScalingWithMaxCountConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_autoScalingErrorMin(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_autoScalingErrorMin(t)
}

func testAccKubernetesCluster_autoScalingErrorMin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoScalingEnabledUpdateMin(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("`node_count`\\(1\\) must be equal to or greater than `min_count`\\(2\\) when `enable_auto_scaling` is set to `true`"),
		},
	})
}

func TestAccKubernetesCluster_autoScalingNodeCountUnset(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_autoScalingNodeCountUnset(t)
}

func testAccKubernetesCluster_autoScalingNodeCountUnset(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoscaleNodeCountUnsetConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.min_count").HasValue("2"),
				check.That(data.ResourceName).Key("default_node_pool.0.max_count").HasValue("4"),
				check.That(data.ResourceName).Key("default_node_pool.0.enable_auto_scaling").HasValue("true"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.max_graceful_termination_sec").HasValue("600"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.new_pod_scale_up_delay").HasValue("0s"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_add").HasValue("10m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_delete").HasValue("10s"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_failure").HasValue("3m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_unneeded").HasValue("10m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_unready").HasValue("20m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_utilization_threshold").HasValue("0.5"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scan_interval").HasValue("10s"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_autoScalingNoAvailabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_autoScalingNoAvailabilityZones(t)
}

func testAccKubernetesCluster_autoScalingNoAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoscaleNoAvailabilityZonesConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.type").HasValue("VirtualMachineScaleSets"),
				check.That(data.ResourceName).Key("default_node_pool.0.min_count").HasValue("1"),
				check.That(data.ResourceName).Key("default_node_pool.0.max_count").HasValue("2"),
				check.That(data.ResourceName).Key("default_node_pool.0.enable_auto_scaling").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_autoScalingWithAvailabilityZones(t)
}

func testAccKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoscaleWithAvailabilityZonesConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.type").HasValue("VirtualMachineScaleSets"),
				check.That(data.ResourceName).Key("default_node_pool.0.min_count").HasValue("1"),
				check.That(data.ResourceName).Key("default_node_pool.0.max_count").HasValue("2"),
				check.That(data.ResourceName).Key("default_node_pool.0.enable_auto_scaling").HasValue("true"),
				check.That(data.ResourceName).Key("default_node_pool.0.availability_zones.#").HasValue("2"),
				check.That(data.ResourceName).Key("default_node_pool.0.availability_zones.0").HasValue("1"),
				check.That(data.ResourceName).Key("default_node_pool.0.availability_zones.1").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_autoScalingProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesCluster_autoScalingProfile(t)
}

func testAccKubernetesCluster_autoScalingProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.autoScalingProfileConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.enable_auto_scaling").HasValue("true"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.max_graceful_termination_sec").HasValue("15"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.new_pod_scale_up_delay").HasValue("10s"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_add").HasValue("10m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_delete").HasValue("10s"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_failure").HasValue("15m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_unneeded").HasValue("15m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_unready").HasValue("15m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_utilization_threshold").HasValue("0.5"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scan_interval").HasValue("10s"),
			),
		},
		data.ImportStep(),
	})
}

func (KubernetesClusterResource) addAgentConfig(data acceptance.TestData, numberOfAgents int) string {
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

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, numberOfAgents)
}

func (KubernetesClusterResource) manualScaleIgnoreChangesConfig(data acceptance.TestData) string {
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

  lifecycle {
    ignore_changes = [
      default_node_pool.0.node_count
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) manualScaleIgnoreChangesConfigUpdated(data acceptance.TestData) string {
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

    tags = {
      Hello = "World"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = [
      default_node_pool.0.node_count
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) autoscaleNodeCountUnsetConfig(data acceptance.TestData) string {
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

func (KubernetesClusterResource) autoscaleNoAvailabilityZonesConfig(data acceptance.TestData) string {
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

func (KubernetesClusterResource) autoscaleWithAvailabilityZonesConfig(data acceptance.TestData) string {
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

func (KubernetesClusterResource) autoScalingProfileConfig(data acceptance.TestData) string {
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
    name                = "default"
    enable_auto_scaling = true
    min_count           = 2
    max_count           = 4
    vm_size             = "Standard_DS2_v2"
  }

  auto_scaler_profile {
    balance_similar_node_groups      = true
    max_graceful_termination_sec     = 15
    new_pod_scale_up_delay           = "10s"
    scan_interval                    = "10s"
    scale_down_delay_after_add       = "10m"
    scale_down_delay_after_delete    = "10s"
    scale_down_delay_after_failure   = "15m"
    scale_down_unneeded              = "15m"
    scale_down_unready               = "15m"
    scale_down_utilization_threshold = "0.5"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, currentKubernetesVersion)
}
