package tests

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var kubernetesNodePoolTests = map[string]func(t *testing.T){
	"autoScale":                      testAccAzureRMKubernetesClusterNodePool_autoScale,
	"autoScaleUpdate":                testAccAzureRMKubernetesClusterNodePool_autoScaleUpdate,
	"availabilityZones":              testAccAzureRMKubernetesClusterNodePool_availabilityZones,
	"errorForAvailabilitySet":        testAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet,
	"multiplePools":                  testAccAzureRMKubernetesClusterNodePool_multiplePools,
	"manualScale":                    testAccAzureRMKubernetesClusterNodePool_manualScale,
	"manualScaleMultiplePools":       testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools,
	"manualScaleMultiplePoolsUpdate": testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate,
	"manualScaleUpdate":              testAccAzureRMKubernetesClusterNodePool_manualScaleUpdate,
	"manualScaleVMSku":               testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku,
	"maxSize":                        testAccAzureRMKubernetesClusterNodePool_maxSize,
	"nodeLabels":                     testAccAzureRMKubernetesClusterNodePool_nodeLabels,
	"nodePublicIP":                   testAccAzureRMKubernetesClusterNodePool_nodePublicIP,
	"nodeTaints":                     testAccAzureRMKubernetesClusterNodePool_nodeTaints,
	"requiresImport":                 testAccAzureRMKubernetesClusterNodePool_requiresImport,
	"spot":                           testAccAzureRMKubernetesClusterNodePool_spot,
	"osDiskSizeGB":                   testAccAzureRMKubernetesClusterNodePool_osDiskSizeGB,
	"proximityPlacementGroupId":      testAccAzureRMKubernetesClusterNodePool_proximityPlacementGroupId,
	"osDiskType":                     testAccAzureRMKubernetesClusterNodePool_osDiskType,
	"modeSystem":                     testAccAzureRMKubernetesClusterNodePool_modeSystem,
	"modeUpdate":                     testAccAzureRMKubernetesClusterNodePool_modeUpdate,
	"virtualNetworkAutomatic":        testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic,
	"virtualNetworkManual":           testAccAzureRMKubernetesClusterNodePool_virtualNetworkManual,
	"windows":                        testAccAzureRMKubernetesClusterNodePool_windows,
	"windowsAndLinux":                testAccAzureRMKubernetesClusterNodePool_windowsAndLinux,
	"zeroSize":                       testAccAzureRMKubernetesClusterNodePool_zeroSize,
}

func TestAccAzureRMKubernetesClusterNodePool_autoScale(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_autoScale(t)
}

func testAccAzureRMKubernetesClusterNodePool_autoScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
			{
				// Disabled
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Staging"),
				),
			},
			data.ImportStep(),
			{
				// Enabled
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_autoScaleUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_autoScaleUpdate(t)
}

func testAccAzureRMKubernetesClusterNodePool_autoScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCountConfig(data, 1, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCountConfig(data, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCountConfig(data, 0, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_availabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_availabilityZones(t)
}

func testAccAzureRMKubernetesClusterNodePool_availabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_availabilityZonesConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet(t)
}

func testAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMKubernetesClusterNodePool_availabilitySetConfig(data),
				ExpectError: regexp.MustCompile("must be a VirtualMachineScaleSet to attach multiple node pools"),
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_multiplePools(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_multiplePools(t)
}

func testAccAzureRMKubernetesClusterNodePool_multiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "autoscale")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_multiplePoolsConfig(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.manual"),
				),
			},
			data.ImportStep(),
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.manual",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScale(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScale(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "first")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.second"),
				),
			},
			data.ImportStep(),
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "first")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCountConfig(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.second"),
				),
			},
			data.ImportStep(),
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCountConfig(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.second"),
				),
			},

			data.ImportStep(),
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleIgnoreChanges(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleIgnoreChanges(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleIgnoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleIgnoreChangesConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "node_count", "1"),
					testCheckAzureRMKubernetesNodePoolScale(data.ResourceName, 2),
				),
			},
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleIgnoreChangesUpdatedConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "node_count", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleUpdate(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCountConfig(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// up
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCountConfig(data, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// and down
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCountConfig(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleVMSkuConfig(data, "Standard_F2s_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleVMSkuConfig(data, "Standard_F4s_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_modeSystem(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_modeSystem(t)
}

func testAccAzureRMKubernetesClusterNodePool_modeSystem(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_modeSystemConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_modeUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_modeUpdate(t)
}

func testAccAzureRMKubernetesClusterNodePool_modeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_modeUserConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesClusterNodePool_modeSystemConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesClusterNodePool_modeUserConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_nodeLabels(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_nodeLabels(t)
}

func testAccAzureRMKubernetesClusterNodePool_nodeLabels(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	labels1 := map[string]string{"key": "value"}
	labels2 := map[string]string{"key2": "value2"}
	labels3 := map[string]string{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodeLabelsConfig(data, labels1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolNodeLabels(data.ResourceName, labels1),
				),
			},
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodeLabelsConfig(data, labels2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolNodeLabels(data.ResourceName, labels2),
				),
			},
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodeLabelsConfig(data, labels3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolNodeLabels(data.ResourceName, labels3),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_nodePublicIP(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_nodePublicIP(t)
}

func testAccAzureRMKubernetesClusterNodePool_nodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodePublicIPConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_nodeTaints(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_nodeTaints(t)
}

func testAccAzureRMKubernetesClusterNodePool_nodeTaints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodeTaintsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(t)
}

func testAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_osDiskSizeGBConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_proximityPlacementGroupId(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_proximityPlacementGroupId(t)
}

func testAccAzureRMKubernetesClusterNodePool_proximityPlacementGroupId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_proximityPlacementGroupIdConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_osDiskType(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_osDiskType(t)
}

func testAccAzureRMKubernetesClusterNodePool_osDiskType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_osDiskTypeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_requiresImport(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_requiresImport(t)
}

func testAccAzureRMKubernetesClusterNodePool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMKubernetesClusterNodePool_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_cluster_node_pool"),
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_spot(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_spot(t)
}

func testAccAzureRMKubernetesClusterNodePool_spot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_spotConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(t)
}

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomaticConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(t)
}

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_virtualNetworkManualConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_windows(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_windows(t)
}

func testAccAzureRMKubernetesClusterNodePool_windows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_windowsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Os", "Windows"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_windowsAndLinux(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_windowsAndLinux(t)
}

func testAccAzureRMKubernetesClusterNodePool_windowsAndLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_windowsAndLinuxConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.linux"),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.windows"),
				),
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.linux",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.windows",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_zeroSize(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_zeroSize(t)
}

func testAccAzureRMKubernetesClusterNodePool_zeroSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_zeroSizeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_maxSize(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_maxSize(t)
}

func testAccAzureRMKubernetesClusterNodePool_maxSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_maxSizeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_sameSize(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_sameSize(t)
}

func testAccAzureRMKubernetesClusterNodePool_sameSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_sameSizeConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMKubernetesClusterNodePoolDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.AgentPoolsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kubernetes_cluster_node_pool" {
			continue
		}

		parsedK8sId, err := parse.KubernetesNodePoolID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing kubernetes node pool id: %+v", err)
		}

		resp, err := client.Get(ctx, parsedK8sId.ResourceGroup, parsedK8sId.ClusterName, parsedK8sId.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Kubernetes Cluster Node Pool still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMKubernetesNodePoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.AgentPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		kubernetesClusterId := rs.Primary.Attributes["kubernetes_cluster_id"]
		parsedK8sId, err := parse.ClusterID(kubernetesClusterId)
		if err != nil {
			return fmt.Errorf("Error parsing kubernetes cluster id: %+v", err)
		}

		agentPool, err := client.Get(ctx, parsedK8sId.ResourceGroup, parsedK8sId.ManagedClusterName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on kubernetesClustersClient: %+v", err)
		}

		if agentPool.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q) does not exist", name, parsedK8sId.ManagedClusterName, parsedK8sId.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMKubernetesNodePoolScale(resourceName string, nodeCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.AgentPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		nodePoolName := rs.Primary.Attributes["name"]
		kubernetesClusterId := rs.Primary.Attributes["kubernetes_cluster_id"]
		parsedK8sId, err := parse.ClusterID(kubernetesClusterId)
		if err != nil {
			return fmt.Errorf("Error parsing kubernetes cluster id: %+v", err)
		}

		clusterName := parsedK8sId.ManagedClusterName
		resourceGroup := parsedK8sId.ResourceGroup

		nodePool, err := client.Get(ctx, resourceGroup, clusterName, nodePoolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on agentPoolsClient: %+v", err)
		}

		if nodePool.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q) does not exist", nodePoolName, clusterName, resourceGroup)
		}

		if nodePool.ManagedClusterAgentPoolProfileProperties == nil {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q): `properties` was nil", nodePoolName, clusterName, resourceGroup)
		}

		nodePool.ManagedClusterAgentPoolProfileProperties.Count = utils.Int32(int32(nodeCount))

		future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, nodePoolName, nodePool)
		if err != nil {
			return fmt.Errorf("Bad: updating node pool %q: %+v", nodePoolName, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: waiting for update of node pool %q: %+v", nodePoolName, err)
		}

		return nil
	}
}

func testCheckAzureRMKubernetesNodePoolNodeLabels(resourceName string, expectedLabels map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.AgentPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		kubernetesClusterId := rs.Primary.Attributes["kubernetes_cluster_id"]
		parsedK8sId, err := parse.ClusterID(kubernetesClusterId)
		if err != nil {
			return fmt.Errorf("Error parsing kubernetes cluster id: %+v", err)
		}

		agent_pool, err := client.Get(ctx, parsedK8sId.ResourceGroup, parsedK8sId.ManagedClusterName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on kubernetesClustersClient: %+v", err)
		}

		if agent_pool.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q) does not exist", name, parsedK8sId.ManagedClusterName, parsedK8sId.ResourceGroup)
		}

		labels := make(map[string]string)
		for k, v := range agent_pool.NodeLabels {
			labels[k] = *v
		}
		if !reflect.DeepEqual(labels, expectedLabels) {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q) nodeLabels %v do not match expected %v", name, parsedK8sId.ManagedClusterName, parsedK8sId.ResourceGroup, labels, expectedLabels)
		}

		return nil
	}
}

func testAccAzureRMKubernetesClusterNodePool_autoScaleConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 1
  max_count             = 3
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCountConfig(data acceptance.TestData, min int, max int) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = %d
  max_count             = %d
}
`, template, min, max)
}

func testAccAzureRMKubernetesClusterNodePool_availabilitySetConfig(data acceptance.TestData) string {
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

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesClusterNodePool_availabilityZonesConfig(data acceptance.TestData) string {
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
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  availability_zones    = ["1"]
  vnet_subnet_id        = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1

  tags = {
    environment = "Staging"
  }
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleIgnoreChangesConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1

  lifecycle {
    ignore_changes = [
      node_count,
    ]
  }
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleIgnoreChangesUpdatedConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1

  tags = {
    Environment = "Staging"
  }

  lifecycle {
    ignore_changes = [
      node_count,
    ]
  }
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "first" {
  name                  = "first"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}

resource "azurerm_kubernetes_cluster_node_pool" "second" {
  name                  = "second"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_F2s_v2"
  node_count            = 1
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCountConfig(data acceptance.TestData, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "first" {
  name                  = "first"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = %d
}

resource "azurerm_kubernetes_cluster_node_pool" "second" {
  name                  = "second"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_F2s_v2"
  node_count            = %d
}
`, template, numberOfAgents, numberOfAgents)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCountConfig(data acceptance.TestData, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = %d
}
`, template, numberOfAgents)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleVMSkuConfig(data acceptance.TestData, sku string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "%s"
  node_count            = 1
}
`, template, sku)
}

func testAccAzureRMKubernetesClusterNodePool_modeSystemConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  mode                  = "System"
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_modeUserConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  mode                  = "User"
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_multiplePoolsConfig(data acceptance.TestData, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "autoscale" {
  name                  = "autoscale"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 1
  max_count             = 3
}

resource "azurerm_kubernetes_cluster_node_pool" "manual" {
  name                  = "manual"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_F2s_v2"
  node_count            = %d
}
`, template, numberOfAgents)
}

func testAccAzureRMKubernetesClusterNodePool_nodeLabelsConfig(data acceptance.TestData, labels map[string]string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	labelsSlice := make([]string, 0, len(labels))
	for k, v := range labels {
		labelsSlice = append(labelsSlice, fmt.Sprintf("    \"%s\" = \"%s\"", k, v))
	}
	labelsStr := strings.Join(labelsSlice, "\n")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  node_labels = {
%s
  }
}
`, template, labelsStr)
}

func testAccAzureRMKubernetesClusterNodePool_nodePublicIPConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  enable_node_public_ip = true
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_nodeTaintsConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  node_taints = [
    "key=value:NoSchedule"
  ]
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_requiresImportConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "import" {
  name                  = azurerm_kubernetes_cluster_node_pool.test.name
  kubernetes_cluster_id = azurerm_kubernetes_cluster_node_pool.test.kubernetes_cluster_id
  vm_size               = azurerm_kubernetes_cluster_node_pool.test.vm_size
  node_count            = azurerm_kubernetes_cluster_node_pool.test.node_count
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_osDiskSizeGBConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  os_disk_size_gb       = 100
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_proximityPlacementGroupIdConfig(data acceptance.TestData) string {
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
resource "azurerm_proximity_placement_group" "test" {
  name                = "acctestPPG-aks-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    environment = "Production"
  }
}
resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                         = "internal"
  kubernetes_cluster_id        = azurerm_kubernetes_cluster.test.id
  vm_size                      = "Standard_DS2_v2"
  node_count                   = 1
  proximity_placement_group_id = azurerm_proximity_placement_group.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesClusterNodePool_osDiskTypeConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS3_v2"
  node_count            = 1
  os_disk_size_gb       = 100
  os_disk_type          = "Ephemeral"
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_spotConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  priority              = "Spot"
  eviction_policy       = "Delete"
  spot_max_price        = 0.5 # high, but this is a maximum (we pay less) so ensures this won't fail
  node_labels = {
    "kubernetes.azure.com/scalesetpriority" = "spot"
  }
  node_taints = [
    "kubernetes.azure.com/scalesetpriority=spot:NoSchedule"
  ]
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomaticConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateVirtualNetworkConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 1
  max_count             = 3
  vnet_subnet_id        = azurerm_subnet.test.id
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkManualConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateVirtualNetworkConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  vnet_subnet_id        = azurerm_subnet.test.id
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_windowsConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateWindowsConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "windoz"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  os_type               = "Windows"

  tags = {
    Os = "Windows"
  }
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_windowsAndLinuxConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateWindowsConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "linux" {
  name                  = "linux"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}

resource "azurerm_kubernetes_cluster_node_pool" "windows" {
  name                  = "windoz"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  os_type               = "Windows"
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_templateConfig(data acceptance.TestData) string {
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
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesClusterNodePool_templateVirtualNetworkConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesClusterNodePool_templateWindowsConfig(data acceptance.TestData) string {
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
  }

  identity {
    type = "SystemAssigned"
  }

  windows_profile {
    admin_username = "azureuser"
    admin_password = "P@55W0rd1234!"
  }

  network_profile {
    network_plugin     = "azure"
    network_policy     = "azure"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesClusterNodePool_zeroSizeConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 0
  max_count             = 3
  node_count            = 0
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_maxSizeConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 1
  max_count             = 1000
  node_count            = 1
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_sameSizeConfig(data acceptance.TestData) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 1
  max_count             = 1
  node_count            = 1
}
`, template)
}
