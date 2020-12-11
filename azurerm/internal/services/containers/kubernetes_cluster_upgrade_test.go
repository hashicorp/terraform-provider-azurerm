package containers_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var kubernetesUpgradeTests = map[string]func(t *testing.T){
	"UpgradeAutoScaleMinCount":                      testAccAzureRMKubernetesCluster_upgradeAutoScaleMinCount,
	"upgradeControlPlane":                           testAccAzureRMKubernetesCluster_upgradeControlPlane,
	"upgradeControlPlaneAndDefaultNodePoolTogether": testAccAzureRMKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTogether,
	"upgradeControlPlaneAndDefaultNodePoolTwoPhase": testAccAzureRMKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTwoPhase,
	"upgradeNodePoolBeforeControlPlaneFails":        testAccAzureRMKubernetesCluster_upgradeNodePoolBeforeControlPlaneFails,
	"upgradeCustomNodePoolAfterControlPlane":        testAccAzureRMKubernetesCluster_upgradeCustomNodePoolAfterControlPlane,
	"upgradeCustomNodePoolBeforeControlPlaneFails":  testAccAzureRMKubernetesCluster_upgradeCustomNodePoolBeforeControlPlaneFails,
}

func TestAccAzureRMKubernetesCluster_upgradeAutoScaleMinCount(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgradeAutoScaleMinCount(t)
}

func testAccAzureRMKubernetesCluster_upgradeAutoScaleMinCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_upgradeAutoScaleMinCountConfig(data, olderKubernetesVersion, 3, 8),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_upgradeAutoScaleMinCountConfig(data, olderKubernetesVersion, 4, 8),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeControlPlane(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgradeControlPlane(t)
}

func testAccAzureRMKubernetesCluster_upgradeControlPlane(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_upgradeControlPlaneConfig(data, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_upgradeControlPlaneConfig(data, currentKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					// the control plane should have been upgraded but the default node pool shouldn't have been
					// TODO: confirm if we can roll the default node pool if the value is unset in the config
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", currentKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTogether(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTogether(t)
}

func testAccAzureRMKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTogether(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesVersion, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesVersion, currentKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", currentKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", currentKubernetesVersion),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTwoPhase(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTwoPhase(t)
}

func testAccAzureRMKubernetesCluster_upgradeControlPlaneAndDefaultNodePoolTwoPhase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesVersion, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesVersion, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", currentKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data, currentKubernetesVersion, currentKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", currentKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", currentKubernetesVersion),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeNodePoolBeforeControlPlaneFails(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgradeNodePoolBeforeControlPlaneFails(t)
}

func testAccAzureRMKubernetesCluster_upgradeNodePoolBeforeControlPlaneFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesVersion, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
			{
				Config:      testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data, olderKubernetesVersion, currentKubernetesVersion),
				ExpectError: regexp.MustCompile("Node Pools cannot use a version of Kubernetes that is not supported on the Control Plane."),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeCustomNodePoolAfterControlPlane(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgradeCustomNodePoolAfterControlPlane(t)
}

func testAccAzureRMKubernetesCluster_upgradeCustomNodePoolAfterControlPlane(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				// all on the older version
				Config: testAccAzureRMKubernetesCluster_upgradeVersionsConfig(data, olderKubernetesVersion, olderKubernetesVersion, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
			{
				// upgrade the control plane
				Config: testAccAzureRMKubernetesCluster_upgradeVersionsConfig(data, currentKubernetesVersion, olderKubernetesVersion, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", currentKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
			{
				// upgrade the node pool
				Config: testAccAzureRMKubernetesCluster_upgradeVersionsConfig(data, currentKubernetesVersion, olderKubernetesVersion, currentKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", currentKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(nodePoolName, "orchestrator_version", currentKubernetesVersion),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeCustomNodePoolBeforeControlPlaneFails(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_upgradeCustomNodePoolBeforeControlPlaneFails(t)
}

func testAccAzureRMKubernetesCluster_upgradeCustomNodePoolBeforeControlPlaneFails(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	nodePoolName := "azurerm_kubernetes_cluster_node_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				// all on the older version
				Config: testAccAzureRMKubernetesCluster_upgradeVersionsConfig(data, olderKubernetesVersion, olderKubernetesVersion, olderKubernetesVersion),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kubernetes_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.orchestrator_version", olderKubernetesVersion),
					resource.TestCheckResourceAttr(nodePoolName, "orchestrator_version", olderKubernetesVersion),
				),
			},
			data.ImportStep(),
			{
				// upgrade the node pool
				Config:      testAccAzureRMKubernetesCluster_upgradeVersionsConfig(data, olderKubernetesVersion, olderKubernetesVersion, currentKubernetesVersion),
				ExpectError: regexp.MustCompile("Node Pools cannot use a version of Kubernetes that is not supported on the Control Plane."),
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_upgradeControlPlaneConfig(data acceptance.TestData, controlPlaneVersion string) string {
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
  kubernetes_version  = %q

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data acceptance.TestData, controlPlaneVersion, defaultNodePoolVersion string) string {
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
  kubernetes_version  = %q

  default_node_pool {
    name                 = "default"
    node_count           = 1
    vm_size              = "Standard_DS2_v2"
    orchestrator_version = %q
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, defaultNodePoolVersion)
}

func testAccAzureRMKubernetesCluster_upgradeVersionsConfig(data acceptance.TestData, controlPlaneVersion, defaultNodePoolVersion, customNodePoolVersion string) string {
	template := testAccAzureRMKubernetesCluster_upgradeControlPlaneDefaultNodePoolConfig(data, controlPlaneVersion, defaultNodePoolVersion)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  orchestrator_version  = %q
}
`, template, customNodePoolVersion)
}

func testAccAzureRMKubernetesCluster_upgradeAutoScaleMinCountConfig(data acceptance.TestData, controlPlaneVersion string, minCount int, maxCount int) string {
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
  kubernetes_version  = %q

  default_node_pool {
    name                = "default"
    vm_size             = "Standard_DS2_v2"
    enable_auto_scaling = true
    min_count           = %d
    max_count           = %d
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion, minCount, maxCount)
}
