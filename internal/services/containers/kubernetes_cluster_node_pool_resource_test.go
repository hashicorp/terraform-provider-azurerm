// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/snapshots"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesClusterNodePoolResource struct{}

func TestAccKubernetesClusterNodePool_autoScale(t *testing.T) {
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

func TestAccKubernetesClusterNodePool_capacityReservationGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.capacityReservationGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_errorForAvailabilitySet(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("AvailabilitySet not supported as an option for default_node_pool in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.availabilitySetConfig(data),
			ExpectError: regexp.MustCompile("multiple node pools are only supported when the Default Node Pool uses a VMScaleSet"),
		},
	})
}

func TestAccKubernetesClusterNodePool_kubeletAndLinuxOSConfig(t *testing.T) {
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

func TestAccKubernetesClusterNodePool_nodeTaints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}
	taints1 := []string{"key=value:NoSchedule"}
	taints2 := []string{"key=value:NoSchedule", "key2=value2:NoSchedule"}
	taints3 := []string{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeTaintsConfig(data, taints1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("node_taints.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_taints.0").HasValue("key=value:NoSchedule"),
			),
		},
		{
			Config: r.nodeTaintsConfig(data, taints2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("node_taints.#").HasValue("2"),
				check.That(data.ResourceName).Key("node_taints.0").HasValue("key=value:NoSchedule"),
				check.That(data.ResourceName).Key("node_taints.1").HasValue("key2=value2:NoSchedule"),
			),
		},
		{
			Config: r.nodeTaintsConfig(data, taints3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("node_taints.#").HasValue("0"),
			),
		},
	})
}

func TestAccKubernetesClusterNodePool_nodeLabels(t *testing.T) {
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

func TestAccKubernetesClusterNodePool_podSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.podSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_osDiskSizeGB(t *testing.T) {
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
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeSettings(data, 35, 18),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.upgradeSettings(data, 5, 0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_virtualNetworkAutomatic(t *testing.T) {
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

func TestAccKubernetesClusterNodePool_virtualNetworkMultipleSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetworkMultipleSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_windows(t *testing.T) {
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

func TestAccKubernetesClusterNodePool_windows2019(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windows2019Config(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.Os").HasValue("Windows"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_windows2022(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windows2022Config(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.Os").HasValue("Windows"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_windowsAndLinux(t *testing.T) {
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
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hostEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_encryption_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccKubernetesClusterNodePool_maxSize(t *testing.T) {
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

func TestAccKubernetesClusterNodePool_ultraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ultraSSD(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ultraSSD(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_osSkuUbuntu(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osSku(data, "Ubuntu"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_osSkuAzureLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osSku(data, "AzureLinux"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_osSkuCBLMariner(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("CBLMariner is an invalid `os_sku` in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osSku(data, "CBLMariner"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_osSkuMariner(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("Mariner is an invalid `os_sku` in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osSku(data, "Mariner"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_osSkuMigration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.osSku(data, "Ubuntu"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.osSku(data, "AzureLinux"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.osSku(data, "Ubuntu"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_dedicatedHost(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dedicatedHost(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_turnOnEnableAutoScalingWithDefaultMaxMinCountSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodePool(data, false, 0, 0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.nodePool(data, true, 0, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_scaleDownMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scaleDownMode(data, "Delete"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scaleDownMode(data, "Deallocate"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_workloadRuntime(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	if !features.FourPointOhBeta() {
		data.ResourceTest(t, r, []acceptance.TestStep{
			{
				Config: r.workloadRuntime(data, "OCIContainer"),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
			{
				Config: r.workloadRuntime(data, "KataMshvVmIsolation"),
				Check: acceptance.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
				),
			},
			data.ImportStep(),
		})
		return
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.workloadRuntime(data, "OCIContainer"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_customCATrustEnabled(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("Skipping this test in 4.0 beta as it is not supported")
	}
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customCATrustEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customCATrustEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_windowsProfileOutboundNatEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsProfileOutboundNatEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_nodeIPTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeIPTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_networkProfileComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkProfileComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_networkProfileUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeIPTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkProfileComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.nodeIPTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesClusterNodePool_snapshotId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.snapshotSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					if _, ok := ctx.Deadline(); !ok {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
						defer cancel()
					}
					client := clients.Containers.SnapshotClient
					poolId, err := agentpools.ParseAgentPoolID(state.ID)
					if err != nil {
						return err
					}
					id := snapshots.NewSnapshotID(poolId.SubscriptionId, poolId.ResourceGroupName, data.RandomString)
					snapshot := snapshots.Snapshot{
						Location: data.Locations.Primary,
						Properties: &snapshots.SnapshotProperties{
							CreationData: &snapshots.CreationData{
								SourceResourceId: utils.String(poolId.ID()),
							},
						},
					}
					_, err = client.CreateOrUpdate(ctx, id, snapshot)
					if err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
					return nil
				}, "azurerm_kubernetes_cluster_node_pool.source"),
			),
		},
		{
			Config: r.snapshotRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.snapshotSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					if _, ok := ctx.Deadline(); !ok {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
						defer cancel()
					}
					client := clients.Containers.SnapshotClient
					poolId, err := agentpools.ParseAgentPoolID(state.ID)
					if err != nil {
						return err
					}
					id := snapshots.NewSnapshotID(poolId.SubscriptionId, poolId.ResourceGroupName, data.RandomString)
					_, err = client.Delete(ctx, id)
					if err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
					return nil
				}, "azurerm_kubernetes_cluster_node_pool.source"),
			),
		},
	})
}

func TestAccKubernetesClusterNodePool_gpuInstance(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gpuInstance(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t KubernetesClusterNodePoolResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := agentpools.ParseAgentPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.AgentPoolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Kubernetes Cluster Node Pool (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil && resp.Model.Id != nil), nil
}

func (KubernetesClusterNodePoolResource) scaleNodePool(nodeCount int) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 1*time.Hour)
			defer cancel()
		}

		nodePoolName := state.Attributes["name"]
		kubernetesClusterId := state.Attributes["kubernetes_cluster_id"]
		parsedK8sId, err := commonids.ParseKubernetesClusterID(kubernetesClusterId)
		if err != nil {
			return fmt.Errorf("parsing kubernetes cluster id: %+v", err)
		}
		parsedAgentPoolId := agentpools.NewAgentPoolID(parsedK8sId.SubscriptionId, parsedK8sId.ResourceGroupName, parsedK8sId.ManagedClusterName, nodePoolName)

		clusterName := parsedK8sId.ManagedClusterName
		resourceGroup := parsedK8sId.ResourceGroupName

		nodePool, err := clients.Containers.AgentPoolsClient.Get(ctx, parsedAgentPoolId)
		if err != nil {
			return fmt.Errorf("Bad: Get on agentPoolsClient: %+v", err)
		}

		if response.WasNotFound(nodePool.HttpResponse) {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q) does not exist", nodePoolName, clusterName, resourceGroup)
		}

		if nodePool.Model == nil || nodePool.Model.Properties == nil {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q): `properties` was nil", nodePoolName, clusterName, resourceGroup)
		}

		nodePool.Model.Properties.Count = utils.Int64(int64(nodeCount))

		err = clients.Containers.AgentPoolsClient.CreateOrUpdateThenPoll(ctx, parsedAgentPoolId, *nodePool.Model)
		if err != nil {
			return fmt.Errorf("Bad: updating node pool %q: %+v", nodePoolName, err)
		}

		return nil
	}
}

func TestAccKubernetesClusterNodePool_virtualNetworkOwnershipRaceCondition(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test1")
	r := KubernetesClusterNodePoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetworkOwnershipRaceCondition(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_kubernetes_cluster_node_pool.test2").ExistsInAzure(r),
				check.That("azurerm_kubernetes_cluster_node_pool.test3").ExistsInAzure(r),
			),
		},
	})
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
  auto_scaling_enabled  = true
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
  auto_scaling_enabled  = true
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
    upgrade_settings {
      max_surge = "10%%"
    }
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
    upgrade_settings {
      max_surge = "10%%"
    }
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
    upgrade_settings {
      max_surge = "10%%"
    }
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
  address_prefixes     = ["10.1.0.0/24"]
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
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "standard"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  vnet_subnet_id        = azurerm_subnet.test.id
  zones                 = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterNodePoolResource) capacityReservationGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_capacity_reservation" "test" {
  name                          = "acctest-ccr-%[1]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id

  sku {
    name     = "Standard_D2s_v3"
    capacity = 2
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_capacity_reservation_group.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
  role_definition_name = "Owner"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"
  default_node_pool {
    name                          = "default"
    node_count                    = 1
    vm_size                       = "Standard_D2s_v3"
    capacity_reservation_group_id = azurerm_capacity_reservation.test.capacity_reservation_group_id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  depends_on = [
    azurerm_capacity_reservation.test,
    azurerm_role_assignment.test
  ]
}
resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                          = "internal"
  kubernetes_cluster_id         = azurerm_kubernetes_cluster.test.id
  vm_size                       = "Standard_D2s_v3"
  node_count                    = 1
  capacity_reservation_group_id = azurerm_capacity_reservation.test.capacity_reservation_group_id
}
`, data.RandomInteger, data.Locations.Primary)
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
  auto_scaling_enabled  = true
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

func (r KubernetesClusterNodePoolResource) nodeTaintsConfig(data acceptance.TestData, taints []string) string {
	taintsSlice := make([]string, 0, len(taints))
	for _, v := range taints {
		taintsSlice = append(taintsSlice, fmt.Sprintf("\"%s\"", v))
	}
	taintsStr := strings.Join(taintsSlice, ",")
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
%s
  ]
}
`, r.templateConfig(data), taintsStr)
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
  node_public_ip_enabled   = true
  node_public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, r.templateConfig(data), data.RandomInteger)
}

func (r KubernetesClusterNodePoolResource) podSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_subnet" "nodesubnet" {
  name                 = "nodesubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.240.0.0/16"]
}
resource "azurerm_subnet" "podsubnet" {
  name                 = "podsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.241.0.0/16"]
  delegation {
    name = "aks-delegation"
    service_delegation {
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Microsoft.ContainerService/managedClusters"
    }
  }
}
resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  sku_tier            = "Standard"
  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_DS2_v2"
    pod_subnet_id  = azurerm_subnet.podsubnet.id
    vnet_subnet_id = azurerm_subnet.nodesubnet.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  network_profile {
    network_plugin = "azure"
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
  pod_subnet_id         = azurerm_subnet.podsubnet.id
  vnet_subnet_id        = azurerm_subnet.nodesubnet.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
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
    upgrade_settings {
      max_surge = "10%%"
    }
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

func (r KubernetesClusterNodePoolResource) upgradeSettings(data acceptance.TestData, drainTimeout int, nodeSoakDuration int) string {
	template := r.templateConfig(data)

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
  upgrade_settings {
    max_surge                     = "10%%"
    drain_timeout_in_minutes      = %d
    node_soak_duration_in_minutes = %d
  }
}
`, template, drainTimeout, nodeSoakDuration)
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
  auto_scaling_enabled  = true
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

func (r KubernetesClusterNodePoolResource) virtualNetworkMultipleSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  vnet_subnet_id        = azurerm_subnet.test2.id
}
`, r.templateVirtualNetworkConfig(data), data.RandomInteger)
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

func (r KubernetesClusterNodePoolResource) windows2019Config(data acceptance.TestData) string {
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
  os_sku                = "Windows2019"
  tags = {
    Os = "Windows"
  }
}
`, r.templateWindowsConfig(data))
}

func (r KubernetesClusterNodePoolResource) windows2022Config(data acceptance.TestData) string {
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
  os_sku                = "Windows2022"
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
    upgrade_settings {
      max_surge = "10%%"
    }
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
  address_prefixes     = ["10.1.0.0/24"]
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
    upgrade_settings {
      max_surge = "10%%"
    }
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
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  windows_profile {
    admin_username = "azureuser"
    admin_password = "P@55W0rd1234!h@2h1C0rP"
  }

  network_profile {
    network_plugin = "azure"
    network_policy = "azure"
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
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
  auto_scaling_enabled  = true
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
  name                    = "internal"
  kubernetes_cluster_id   = azurerm_kubernetes_cluster.test.id
  vm_size                 = "Standard_DS2_v2"
  host_encryption_enabled = true
  node_count              = 1
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
  auto_scaling_enabled  = true
  min_count             = 1
  max_count             = 399
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
  auto_scaling_enabled  = true
  min_count             = 1
  max_count             = 1
  node_count            = 1
}
`, r.templateConfig(data))
}

func (r KubernetesClusterNodePoolResource) other(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
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
  message_of_the_day    = "daily message"
}
`, r.templateConfig(data))
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
  fips_enabled          = true
  kubelet_disk_type     = "OS"
}
`, r.templateConfig(data))
}

func (KubernetesClusterNodePoolResource) ultraSSD(data acceptance.TestData, ultraSSDEnabled bool) string {
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
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}
resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
  ultra_ssd_enabled     = %t
  zones                 = ["1", "2", "3"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, ultraSSDEnabled)
}

func (KubernetesClusterNodePoolResource) osSku(data acceptance.TestData, osSku string) string {
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
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}
resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
  os_sku                = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, osSku)
}

func (KubernetesClusterNodePoolResource) dedicatedHost(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctestDHG-compute-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 3
  automatic_placement_enabled = true
}

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%[1]d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type3"
  platform_fault_domain   = 0
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  principal_id         = azurerm_user_assigned_identity.test.principal_id
  role_definition_name = "Contributor"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
  node_count            = 1
  host_group_id         = azurerm_dedicated_host_group.test.id

  depends_on = [
    azurerm_role_assignment.test,
    azurerm_dedicated_host.test
  ]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r KubernetesClusterNodePoolResource) nodePool(data acceptance.TestData, enableAutoScaling bool, minCount, maxCount int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  auto_scaling_enabled  = %t
  min_count             = %d
  max_count             = %d
}
`, r.templateConfig(data), enableAutoScaling, minCount, maxCount)
}

func (KubernetesClusterNodePoolResource) scaleDownMode(data acceptance.TestData, scaleDownMode string) string {
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
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}
resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
  scale_down_mode       = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, scaleDownMode)
}

func (KubernetesClusterNodePoolResource) workloadRuntime(data acceptance.TestData, workloadRuntime string) string {
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
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}
resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
  workload_runtime      = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, workloadRuntime)
}

func (KubernetesClusterNodePoolResource) customCATrustEnabled(data acceptance.TestData, enabled bool) string {
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
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}
resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                    = "internal"
  kubernetes_cluster_id   = azurerm_kubernetes_cluster.test.id
  vm_size                 = "Standard_D2s_v3"
  custom_ca_trust_enabled = "%t"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, enabled)
}

func (KubernetesClusterNodePoolResource) windowsProfileOutboundNatEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
  network_profile {
    network_plugin = "azure"
    outbound_type  = "managedNATGateway"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "user"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
  os_type               = "Windows"
  os_sku                = "Windows2019"
  windows_profile {
    outbound_nat_enabled = true
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesClusterNodePoolResource) nodeIPTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctestasg-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                   = "internal"
  kubernetes_cluster_id  = azurerm_kubernetes_cluster.test.id
  vm_size                = "Standard_D2s_v3"
  node_public_ip_enabled = true
  node_network_profile {
    node_public_ip_tags = {
      RoutingPreference = "Internet"
    }
  }
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesClusterNodePoolResource) networkProfileComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctestasg-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                   = "internal"
  kubernetes_cluster_id  = azurerm_kubernetes_cluster.test.id
  vm_size                = "Standard_D2s_v3"
  node_public_ip_enabled = true
  node_network_profile {
    allowed_host_ports {
      port_start = 8001
      port_end   = 8002
      protocol   = "UDP"
    }
    application_security_group_ids = [azurerm_application_security_group.test.id]
    node_public_ip_tags = {
      RoutingPreference = "Internet"
    }
  }
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesClusterNodePoolResource) snapshotSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "source" {
  name                  = "source"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesClusterNodePoolResource) snapshotRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "source" {
  name                  = "source"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_D2s_v3"
}

data "azurerm_kubernetes_node_pool_snapshot" "test" {
  name                = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "new"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  snapshot_id           = data.azurerm_kubernetes_node_pool_snapshot.test.id
  depends_on = [
    azurerm_kubernetes_cluster_node_pool.source
  ]
}
 `, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (KubernetesClusterNodePoolResource) gpuInstance(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2s_v3"
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_NC24ads_A100_v4"
  gpu_instance          = "MIG1g"
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesClusterNodePoolResource) virtualNetworkOwnershipRaceCondition(data acceptance.TestData) string {
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
  count                = 8
  name                 = "acctestsubnet%d${count.index}"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.${count.index}.0/24"]
  service_endpoints = [
    "Microsoft.Storage.Global",
    "Microsoft.AzureActiveDirectory",
    "Microsoft.KeyVault"
  ]

  lifecycle {
    # AKS automatically configures subnet delegations when the subnets are assigned
    # to node pools. We ignore changes so the terraform refresh run by Terraform's plugin-sdk,
    # at the end of the test, returns empty and leaves the test succeed.
    ignore_changes = [delegation]
  }
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
    vnet_subnet_id = azurerm_subnet.test["6"].id
    pod_subnet_id  = azurerm_subnet.test["7"].id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "standard"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test1" {
  name                  = "internal1"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_L8s_v3"
  node_count            = 1
  vnet_subnet_id        = azurerm_subnet.test[0].id
  pod_subnet_id         = azurerm_subnet.test[1].id
  zones                 = ["1"]
}

resource "azurerm_kubernetes_cluster_node_pool" "test2" {
  name                  = "internal2"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_L8s_v3"
  node_count            = 1
  vnet_subnet_id        = azurerm_subnet.test[2].id
  pod_subnet_id         = azurerm_subnet.test[3].id
  zones                 = ["1"]
}

resource "azurerm_kubernetes_cluster_node_pool" "test3" {
  name                  = "internal3"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_L8s_v3"
  node_count            = 1
  vnet_subnet_id        = azurerm_subnet.test[4].id
  pod_subnet_id         = azurerm_subnet.test[5].id
  zones                 = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
