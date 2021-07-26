package containers_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KubernetesClusterNodePoolResource struct {
}

var kubernetesNodePoolTests = map[string]func(t *testing.T){
	"autoScale":                      testAccKubernetesClusterNodePool_autoScale,
	"autoScaleUpdate":                testAccKubernetesClusterNodePool_autoScaleUpdate,
	"availabilityZones":              testAccKubernetesClusterNodePool_availabilityZones,
	"errorForAvailabilitySet":        testAccKubernetesClusterNodePool_errorForAvailabilitySet,
	"kubeletAndLinuxOSConfig":        testAccKubernetesClusterNodePool_kubeletAndLinuxOSConfig,
	"kubeletAndLinuxOSConfigPartial": testAccKubernetesClusterNodePool_kubeletAndLinuxOSConfigPartial,
	"other":                          testAccKubernetesClusterNodePool_other,
	"multiplePools":                  testAccKubernetesClusterNodePool_multiplePools,
	"manualScale":                    testAccKubernetesClusterNodePool_manualScale,
	"manualScaleMultiplePools":       testAccKubernetesClusterNodePool_manualScaleMultiplePools,
	"manualScaleMultiplePoolsUpdate": testAccKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate,
	"manualScaleUpdate":              testAccKubernetesClusterNodePool_manualScaleUpdate,
	"manualScaleVMSku":               testAccKubernetesClusterNodePool_manualScaleVMSku,
	"maxSize":                        testAccKubernetesClusterNodePool_maxSize,
	"nodeLabels":                     testAccKubernetesClusterNodePool_nodeLabels,
	"nodePublicIP":                   testAccKubernetesClusterNodePool_nodePublicIP,
	"nodeTaints":                     testAccKubernetesClusterNodePool_nodeTaints,
	"requiresImport":                 testAccKubernetesClusterNodePool_requiresImport,
	"spot":                           testAccKubernetesClusterNodePool_spot,
	"osDiskSizeGB":                   testAccKubernetesClusterNodePool_osDiskSizeGB,
	"proximityPlacementGroupId":      testAccKubernetesClusterNodePool_proximityPlacementGroupId,
	"osDiskType":                     testAccKubernetesClusterNodePool_osDiskType,
	"modeSystem":                     testAccKubernetesClusterNodePool_modeSystem,
	"modeUpdate":                     testAccKubernetesClusterNodePool_modeUpdate,
	"upgradeSettings":                testAccKubernetesClusterNodePool_upgradeSettings,
	"virtualNetworkAutomatic":        testAccKubernetesClusterNodePool_virtualNetworkAutomatic,
	"virtualNetworkManual":           testAccKubernetesClusterNodePool_virtualNetworkManual,
	"windows":                        testAccKubernetesClusterNodePool_windows,
	"windowsAndLinux":                testAccKubernetesClusterNodePool_windowsAndLinux,
	"zeroSize":                       testAccKubernetesClusterNodePool_zeroSize,
	"hostEncryption":                 testAccKubernetesClusterNodePool_hostEncryption,
}

func TestAccKubernetesClusterNodePool_autoScale(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_autoScale(t)
}

func testAccKubernetesClusterNodePool_autoScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// Enabled
			Config: r.autoScaleConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			// Disabled
			Config: r.manualScaleConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Staging"),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.autoScaleConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_autoScaleUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_autoScaleUpdate(t)
}

func testAccKubernetesClusterNodePool_autoScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoScaleNodeCountConfig(data, 1, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoScaleNodeCountConfig(data, 3, 5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoScaleNodeCountConfig(data, 0, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_availabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_availabilityZones(t)
}

func testAccKubernetesClusterNodePool_availabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.availabilityZonesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_errorForAvailabilitySet(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_errorForAvailabilitySet(t)
}

func testAccKubernetesClusterNodePool_errorForAvailabilitySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.availabilitySetConfig(data),
			ExpectError: regexp.MustCompile("must be a VirtualMachineScaleSet to attach multiple node pools"),
		},
	})
}

func TestAccKubernetesClusterNodePool_kubeletAndLinuxOSConfig(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_kubeletAndLinuxOSConfig(t)
}

func testAccKubernetesClusterNodePool_kubeletAndLinuxOSConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.kubeletAndLinuxOSConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_kubeletAndLinuxOSConfigPartial(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_kubeletAndLinuxOSConfigPartial(t)
}

func testAccKubernetesClusterNodePool_kubeletAndLinuxOSConfigPartial(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.kubeletAndLinuxOSConfigPartial(data),
			Check: acceptance.ComposeTestCheckFunc(

				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_other(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_other(t)
}

func testAccKubernetesClusterNodePool_other(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.other(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_multiplePools(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_multiplePools(t)
}

func testAccKubernetesClusterNodePool_multiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "autoscale")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiplePoolsConfig(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			ResourceName:      "azurerm_kubernetes_cluster_node_pool.manual",
			ImportState:       true,
			ImportStateVerify: true,
		},
	})
}

func TestAccKubernetesClusterNodePool_manualScale(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_manualScale(t)
}

func testAccKubernetesClusterNodePool_manualScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualScaleConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_manualScaleMultiplePools(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_manualScaleMultiplePools(t)
}

func testAccKubernetesClusterNodePool_manualScaleMultiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "first")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualScaleMultiplePoolsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_kubernetes_cluster_node_pool.second").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
			ImportState:       true,
			ImportStateVerify: true,
		},
	})
}

func TestAccKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t)
}

func testAccKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "first")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualScaleMultiplePoolsNodeCountConfig(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_kubernetes_cluster_node_pool.second").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
			ImportState:       true,
			ImportStateVerify: true,
		},
		{
			Config: r.manualScaleMultiplePoolsNodeCountConfig(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_kubernetes_cluster_node_pool.second").ExistsInAzure(r),
			),
		},

		data.ImportStep(),
		{
			ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
			ImportState:       true,
			ImportStateVerify: true,
		},
	})
}

func TestAccKubernetesClusterNodePool_manualScaleIgnoreChanges(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_manualScaleIgnoreChanges(t)
}

func testAccKubernetesClusterNodePool_manualScaleIgnoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualScaleIgnoreChangesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_count").HasValue("1"),
				data.CheckWithClient(r.scaleNodePool(2)),
			),
		},
		{
			Config: r.manualScaleIgnoreChangesUpdatedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_count").HasValue("2"),
			),
		},
	})
}

func TestAccKubernetesClusterNodePool_manualScaleUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_manualScaleUpdate(t)
}

func testAccKubernetesClusterNodePool_manualScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualScaleNodeCountConfig(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// up
			Config: r.manualScaleNodeCountConfig(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// and down
			Config: r.manualScaleNodeCountConfig(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_manualScaleVMSku(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_manualScaleVMSku(t)
}

func testAccKubernetesClusterNodePool_manualScaleVMSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualScaleVMSkuConfig(data, "Standard_F2s_v2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.manualScaleVMSkuConfig(data, "Standard_F4s_v2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_modeSystem(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_modeSystem(t)
}

func testAccKubernetesClusterNodePool_modeSystem(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.modeSystemConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_modeUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_modeUpdate(t)
}

func testAccKubernetesClusterNodePool_modeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.modeUserConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.modeSystemConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.modeUserConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_nodeLabels(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_nodeLabels(t)
}

func testAccKubernetesClusterNodePool_nodeLabels(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}
	labels1 := map[string]string{"key": "value"}
	labels2 := map[string]string{"key2": "value2"}
	labels3 := map[string]string{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeLabelsConfig(data, labels1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("node_labels.%").HasValue("1"),
				check.That(data.ResourceName).Key("node_labels.key").HasValue("value"),
			),
		},
		{
			Config: r.nodeLabelsConfig(data, labels2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("node_labels.%").HasValue("1"),
				check.That(data.ResourceName).Key("node_labels.key2").HasValue("value2"),
			),
		},
		{
			Config: r.nodeLabelsConfig(data, labels3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("node_labels.%").HasValue("0"),
			),
		},
	})
}

func TestAccKubernetesClusterNodePool_nodePublicIP(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_nodePublicIP(t)
}

func testAccKubernetesClusterNodePool_nodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodePublicIPConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_nodeTaints(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_nodeTaints(t)
}

func testAccKubernetesClusterNodePool_nodeTaints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeTaintsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_osDiskSizeGB(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_osDiskSizeGB(t)
}

func testAccKubernetesClusterNodePool_osDiskSizeGB(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osDiskSizeGBConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_proximityPlacementGroupId(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_proximityPlacementGroupId(t)
}

func testAccKubernetesClusterNodePool_proximityPlacementGroupId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.proximityPlacementGroupIdConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_osDiskType(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_osDiskType(t)
}

func testAccKubernetesClusterNodePool_osDiskType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osDiskTypeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_requiresImport(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_requiresImport(t)
}

func testAccKubernetesClusterNodePool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualScaleConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_cluster_node_pool"),
		},
	})
}

func TestAccKubernetesClusterNodePool_spot(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_spot(t)
}

func testAccKubernetesClusterNodePool_spot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.spotConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_upgradeSettings(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_upgradeSettings(t)
}

func testAccKubernetesClusterNodePool_upgradeSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeSettingsConfig(data, "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("upgrade_settings.0.max_surge").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeSettingsConfig(data, "4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("upgrade_settings.0.max_surge").HasValue("4"),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeSettingsConfig(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("upgrade_settings.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_virtualNetworkAutomatic(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_virtualNetworkAutomatic(t)
}

func testAccKubernetesClusterNodePool_virtualNetworkAutomatic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetworkAutomaticConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_virtualNetworkManual(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_virtualNetworkManual(t)
}

func testAccKubernetesClusterNodePool_virtualNetworkManual(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetworkManualConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_windows(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_windows(t)
}

func testAccKubernetesClusterNodePool_windows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.Os").HasValue("Windows"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_windowsAndLinux(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_windowsAndLinux(t)
}

func testAccKubernetesClusterNodePool_windowsAndLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsAndLinuxConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_kubernetes_cluster_node_pool.linux").ExistsInAzure(r),
				check.That("azurerm_kubernetes_cluster_node_pool.windows").ExistsInAzure(r),
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
	})
}

func TestAccKubernetesClusterNodePool_zeroSize(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_zeroSize(t)
}

func testAccKubernetesClusterNodePool_zeroSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zeroSizeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_hostEncryption(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_hostEncryption(t)
}

func testAccKubernetesClusterNodePool_hostEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hostEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_host_encryption").HasValue("true"),
			),
		},
	})
}

func TestAccKubernetesClusterNodePool_maxSize(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_maxSize(t)
}

func testAccKubernetesClusterNodePool_maxSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.maxSizeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_sameSize(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccKubernetesClusterNodePool_sameSize(t)
}

func testAccKubernetesClusterNodePool_sameSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sameSizeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t KubernetesClusterNodePoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NodePoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.AgentPoolsClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
	if err != nil {
		return nil, fmt.Errorf("reading Kubernetes Cluster Node Pool (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (KubernetesClusterNodePoolResource) scaleNodePool(nodeCount int) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		nodePoolName := state.Attributes["name"]
		kubernetesClusterId := state.Attributes["kubernetes_cluster_id"]
		parsedK8sId, err := parse.ClusterID(kubernetesClusterId)
		if err != nil {
			return fmt.Errorf("Error parsing kubernetes cluster id: %+v", err)
		}

		clusterName := parsedK8sId.ManagedClusterName
		resourceGroup := parsedK8sId.ResourceGroup

		nodePool, err := clients.Containers.AgentPoolsClient.Get(ctx, resourceGroup, clusterName, nodePoolName)
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

		future, err := clients.Containers.AgentPoolsClient.CreateOrUpdate(ctx, resourceGroup, clusterName, nodePoolName, nodePool)
		if err != nil {
			return fmt.Errorf("Bad: updating node pool %q: %+v", nodePoolName, err)
		}

		if err := future.WaitForCompletionRef(ctx, clients.Containers.AgentPoolsClient.Client); err != nil {
			return fmt.Errorf("Bad: waiting for update of node pool %q: %+v", nodePoolName, err)
		}

		return nil
	}
}

func (r KubernetesClusterNodePoolResource) autoScaleConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) autoScaleNodeCountConfig(data acceptance.TestData, min int, max int) string {
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
`, r.templateConfig(data), min, max)
}

func (KubernetesClusterNodePoolResource) availabilitySetConfig(data acceptance.TestData) string {
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

func (KubernetesClusterNodePoolResource) kubeletAndLinuxOSConfig(data acceptance.TestData) string {
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

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1

  kubelet_config {
    cpu_manager_policy        = "static"
    cpu_cfs_quota_enabled     = true
    cpu_cfs_quota_period      = "10ms"
    image_gc_high_threshold   = 90
    image_gc_low_threshold    = 70
    topology_manager_policy   = "best-effort"
    allowed_unsafe_sysctls    = ["kernel.msg*", "net.core.somaxconn"]
    container_log_max_size_mb = 100
    container_log_max_line    = 100000
    pod_max_pid               = 12345
  }

  linux_os_config {
    transparent_huge_page_enabled = "always"
    transparent_huge_page_defrag  = "always"
    swap_file_size_mb             = 300

    sysctl_config {
      fs_aio_max_nr                      = 65536
      fs_file_max                        = 100000
      fs_inotify_max_user_watches        = 1000000
      fs_nr_open                         = 1048576
      kernel_threads_max                 = 200000
      net_core_netdev_max_backlog        = 1800
      net_core_optmem_max                = 30000
      net_core_rmem_max                  = 300000
      net_core_rmem_default              = 300000
      net_core_somaxconn                 = 5000
      net_core_wmem_default              = 300000
      net_core_wmem_max                  = 300000
      net_ipv4_ip_local_port_range_min   = 32768
      net_ipv4_ip_local_port_range_max   = 60000
      net_ipv4_neigh_default_gc_thresh1  = 128
      net_ipv4_neigh_default_gc_thresh2  = 512
      net_ipv4_neigh_default_gc_thresh3  = 1024
      net_ipv4_tcp_fin_timeout           = 60
      net_ipv4_tcp_keepalive_probes      = 9
      net_ipv4_tcp_keepalive_time        = 6000
      net_ipv4_tcp_max_syn_backlog       = 2048
      net_ipv4_tcp_max_tw_buckets        = 100000
      net_ipv4_tcp_tw_reuse              = true
      net_ipv4_tcp_keepalive_intvl       = 70
      net_netfilter_nf_conntrack_buckets = 65536
      net_netfilter_nf_conntrack_max     = 200000
      vm_max_map_count                   = 65536
      vm_swappiness                      = 45
      vm_vfs_cache_pressure              = 80
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterNodePoolResource) kubeletAndLinuxOSConfigPartial(data acceptance.TestData) string {
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

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1

  kubelet_config {
    cpu_manager_policy    = "static"
    cpu_cfs_quota_enabled = true
    cpu_cfs_quota_period  = "10ms"
  }

  linux_os_config {
    transparent_huge_page_enabled = "always"

    sysctl_config {
      fs_aio_max_nr               = 65536
      fs_file_max                 = 100000
      fs_inotify_max_user_watches = 1000000
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterNodePoolResource) availabilityZonesConfig(data acceptance.TestData) string {
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

func (r KubernetesClusterNodePoolResource) manualScaleConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) manualScaleIgnoreChangesConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) manualScaleIgnoreChangesUpdatedConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) manualScaleMultiplePoolsConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) manualScaleMultiplePoolsNodeCountConfig(data acceptance.TestData, numberOfAgents int) string {
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
`, r.templateConfig(data), numberOfAgents, numberOfAgents)
}

func (r KubernetesClusterNodePoolResource) manualScaleNodeCountConfig(data acceptance.TestData, numberOfAgents int) string {
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
`, r.templateConfig(data), numberOfAgents)
}

func (r KubernetesClusterNodePoolResource) manualScaleVMSkuConfig(data acceptance.TestData, sku string) string {
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
`, r.templateConfig(data), sku)
}

func (r KubernetesClusterNodePoolResource) modeSystemConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) modeUserConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) multiplePoolsConfig(data acceptance.TestData, numberOfAgents int) string {
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
`, r.templateConfig(data), numberOfAgents)
}

func (r KubernetesClusterNodePoolResource) nodeLabelsConfig(data acceptance.TestData, labels map[string]string) string {
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
`, r.templateConfig(data), labelsStr)
}

func (r KubernetesClusterNodePoolResource) nodePublicIPConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpipprefix%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  prefix_length       = 31
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                     = "internal"
  kubernetes_cluster_id    = azurerm_kubernetes_cluster.test.id
  vm_size                  = "Standard_DS2_v2"
  node_count               = 1
  enable_node_public_ip    = true
  node_public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, r.templateConfig(data), data.RandomInteger)
}

func (r KubernetesClusterNodePoolResource) nodeTaintsConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "import" {
  name                  = azurerm_kubernetes_cluster_node_pool.test.name
  kubernetes_cluster_id = azurerm_kubernetes_cluster_node_pool.test.kubernetes_cluster_id
  vm_size               = azurerm_kubernetes_cluster_node_pool.test.vm_size
  node_count            = azurerm_kubernetes_cluster_node_pool.test.node_count
}
`, r.manualScaleConfig(data))
}

func (r KubernetesClusterNodePoolResource) osDiskSizeGBConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (KubernetesClusterNodePoolResource) proximityPlacementGroupIdConfig(data acceptance.TestData) string {
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

func (r KubernetesClusterNodePoolResource) osDiskTypeConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) spotConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) upgradeSettingsConfig(data acceptance.TestData, maxSurge string) string {
	template := r.templateConfig(data)
	if maxSurge != "" {
		maxSurge = fmt.Sprintf(`upgrade_settings {
    max_surge = %q
  }`, maxSurge)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 3
  %s
}
`, template, maxSurge)
}

func (r KubernetesClusterNodePoolResource) virtualNetworkAutomaticConfig(data acceptance.TestData) string {
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
`, r.templateVirtualNetworkConfig(data))
}

func (r KubernetesClusterNodePoolResource) virtualNetworkManualConfig(data acceptance.TestData) string {
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
`, r.templateVirtualNetworkConfig(data))
}

func (r KubernetesClusterNodePoolResource) windowsConfig(data acceptance.TestData) string {
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
`, r.templateWindowsConfig(data))
}

func (r KubernetesClusterNodePoolResource) windowsAndLinuxConfig(data acceptance.TestData) string {
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
`, r.templateWindowsConfig(data))
}

func (KubernetesClusterNodePoolResource) templateConfig(data acceptance.TestData) string {
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

func (KubernetesClusterNodePoolResource) templateVirtualNetworkConfig(data acceptance.TestData) string {
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

func (KubernetesClusterNodePoolResource) templateWindowsConfig(data acceptance.TestData) string {
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
    admin_password = "P@55W0rd1234!h@2h1C0rP"
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

func (r KubernetesClusterNodePoolResource) zeroSizeConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) hostEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

	%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                   = "internal"
  kubernetes_cluster_id  = azurerm_kubernetes_cluster.test.id
  vm_size                = "Standard_DS2_v2"
  enable_host_encryption = true
  node_count             = 1
}
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) maxSizeConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) sameSizeConfig(data acceptance.TestData) string {
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
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) other(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 3
  fips_enabled          = true
  kubelet_disk_type     = "OS"
}
`, r.templateConfig(data))
}
