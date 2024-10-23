// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/agentpools"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func TestAccKubernetesCluster_updateVmSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withHostEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateVmSize(data, "Standard_DS3_v2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_node_pool.0.temporary_name_for_rotation"),
	})
}

func TestAccKubernetesCluster_updateVmSizeAfterFailureWithTempAndDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// create the temporary node pool to simulate the case where both old default node pool and temp node pool exist
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					if _, ok := ctx.Deadline(); !ok {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 1*time.Hour)
						defer cancel()
					}

					client := clients.Containers.AgentPoolsClient

					id, err := commonids.ParseKubernetesClusterID(state.Attributes["id"])
					if err != nil {
						return err
					}

					defaultNodePoolId := agentpools.NewAgentPoolID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, state.Attributes["default_node_pool.0.name"])

					resp, err := client.Get(ctx, defaultNodePoolId)
					if err != nil {
						return fmt.Errorf("retrieving %s: %+v", defaultNodePoolId, err)
					}
					if resp.Model == nil {
						return fmt.Errorf("retrieving %s: model was nil", defaultNodePoolId)
					}

					tempNodePoolName := "temp"
					profile := resp.Model
					profile.Name = &tempNodePoolName
					profile.Properties.VMSize = pointer.To("Standard_DS3_v2")

					tempNodePoolId := agentpools.NewAgentPoolID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, tempNodePoolName)
					if err := client.CreateOrUpdateThenPoll(ctx, tempNodePoolId, *profile); err != nil {
						return fmt.Errorf("creating %s: %+v", tempNodePoolId, err)
					}

					return nil
				}, data.ResourceName),
			),
		},
		{
			Config: r.updateVmSize(data, "Standard_DS3_v2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_node_pool.0.temporary_name_for_rotation"),
	})

}

func TestAccKubernetesCluster_updateVmSizeAfterFailureWithTempWithoutDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithTempName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// create the temporary node pool and delete the old default node pool to simulate the case where resizing fails when trying to bring up the new node pool
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					if _, ok := ctx.Deadline(); !ok {
						var cancel context.CancelFunc
						ctx, cancel = context.WithTimeout(ctx, 1*time.Hour)
						defer cancel()
					}

					client := clients.Containers.AgentPoolsClient

					id, err := commonids.ParseKubernetesClusterID(state.Attributes["id"])
					if err != nil {
						return err
					}

					defaultNodePoolId := agentpools.NewAgentPoolID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, state.Attributes["default_node_pool.0.name"])

					resp, err := client.Get(ctx, defaultNodePoolId)
					if err != nil {
						return fmt.Errorf("retrieving %s: %+v", defaultNodePoolId, err)
					}
					if resp.Model == nil {
						return fmt.Errorf("retrieving %s: model was nil", defaultNodePoolId)
					}

					tempNodePoolName := "temp"
					profile := resp.Model
					profile.Name = &tempNodePoolName
					profile.Properties.VMSize = pointer.To("Standard_DS3_v2")

					tempNodePoolId := agentpools.NewAgentPoolID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, tempNodePoolName)
					if err := client.CreateOrUpdateThenPoll(ctx, tempNodePoolId, *profile); err != nil {
						return fmt.Errorf("creating %s: %+v", tempNodePoolId, err)
					}

					if err := client.DeleteThenPoll(ctx, defaultNodePoolId); err != nil {
						return fmt.Errorf("deleting default %s: %+v", defaultNodePoolId, err)
					}

					return nil
				}, data.ResourceName),
			),
			// the plan will show that the default node pool name has been set to "temp" and we're trying to set it back to "default"
			ExpectNonEmptyPlan: true,
		},
		{
			Config: r.updateVmSize(data, "Standard_DS3_v2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_node_pool.0.temporary_name_for_rotation"),
	})
}

func TestAccKubernetesCluster_cycleSystemNodePool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withHostTempDiskVmSize(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateOsDisk(data, "Ephemeral", 75),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_node_pool.0.temporary_name_for_rotation"),
		{
			Config: r.updateZones(data, "Standard_D2ads_v5", "[1,2,3]"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_node_pool.0.temporary_name_for_rotation"),
		{
			Config: r.updateLinuxKernelSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_node_pool.0.temporary_name_for_rotation"),
	})
}

func TestAccKubernetesCluster_cycleSystemNodePoolFipsEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableFips(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("default_node_pool.0.temporary_name_for_rotation"),
	})
}

func TestAccKubernetesCluster_addAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.addAgentConfig(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("1"),
			),
		},
		{
			Config: r.addAgentConfig(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
	})
}

func TestAccKubernetesCluster_manualScaleIgnoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualScaleIgnoreChangesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("1"),
				data.CheckWithClient(r.updateDefaultNodePoolAgentCount(2)),
			),
		},
		{
			Config: r.manualScaleIgnoreChangesConfigUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
	})
}

func TestAccKubernetesCluster_removeAgent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.addAgentConfig(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
		{
			Config: r.addAgentConfig(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("1"),
			),
		},
	})
}

func TestAccKubernetesCluster_autoScalingError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoScalingEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_count").HasValue("2"),
			),
		},
		{
			Config: r.autoScalingEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("cannot change `node_count` when `auto_scaling_enabled` is set to `true`"),
		},
	})
}

func TestAccKubernetesCluster_autoScalingErrorMax(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoScalingEnabledUpdateMax(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("`node_count`\\(11\\) must be equal to or less than `max_count`\\(10\\) when `auto_scaling_enabled` is set to `true`"),
		},
	})
}

func TestAccKubernetesCluster_autoScalingWithMaxCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoScalingWithMaxCountConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_autoScalingErrorMin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoScalingEnabledUpdateMin(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectError: regexp.MustCompile("`node_count`\\(1\\) must be equal to or greater than `min_count`\\(2\\) when `auto_scaling_enabled` is set to `true`"),
		},
	})
}

func TestAccKubernetesCluster_autoScalingNodeCountUnset(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoscaleNodeCountUnsetConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.min_count").HasValue("2"),
				check.That(data.ResourceName).Key("default_node_pool.0.max_count").HasValue("4"),
				check.That(data.ResourceName).Key("default_node_pool.0.auto_scaling_enabled").HasValue("true"),
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
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoscaleNoAvailabilityZonesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.type").HasValue("VirtualMachineScaleSets"),
				check.That(data.ResourceName).Key("default_node_pool.0.min_count").HasValue("1"),
				check.That(data.ResourceName).Key("default_node_pool.0.max_count").HasValue("2"),
				check.That(data.ResourceName).Key("default_node_pool.0.auto_scaling_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_autoScalingWithAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoscaleWithAvailabilityZonesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesCluster_autoScalingProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoScalingProfileConfigMinimal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.auto_scaling_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.expander").HasValue("random"),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoScalingProfileConfigComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.auto_scaling_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.expander").HasValue("least-waste"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.max_graceful_termination_sec").HasValue("15"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.max_node_provisioning_time").HasValue("10m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.max_unready_percentage").HasValue("50"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.new_pod_scale_up_delay").HasValue("10s"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.max_unready_nodes").HasValue("5"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_add").HasValue("10m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_delete").HasValue("10s"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_delay_after_failure").HasValue("15m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_unneeded").HasValue("15m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_unready").HasValue("15m"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scale_down_utilization_threshold").HasValue("0.5"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.empty_bulk_delete_max").HasValue("50"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.scan_interval").HasValue("10s"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.skip_nodes_with_local_storage").HasValue("false"),
				check.That(data.ResourceName).Key("auto_scaler_profile.0.skip_nodes_with_system_pods").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (KubernetesClusterResource) basic(data acceptance.TestData) string {
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

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) withHostEncryption(data acceptance.TestData) string {
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
    name                    = "default"
    node_count              = 1
    vm_size                 = "Standard_DS2_v2"
    host_encryption_enabled = true
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) withHostTempDiskVmSize(data acceptance.TestData) string {
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
    vm_size    = "Standard_D2ads_v5"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) basicWithTempName(data acceptance.TestData) string {
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
    name                        = "default"
    temporary_name_for_rotation = "temp"
    node_count                  = 1
    vm_size                     = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) updateVmSize(data acceptance.TestData, vmSize string) string {
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
    name                        = "default"
    temporary_name_for_rotation = "temp"
    node_count                  = 1
    vm_size                     = "%s"
    host_encryption_enabled     = false
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, vmSize)
}

func (KubernetesClusterResource) updateZones(data acceptance.TestData, vmSize, zones string) string {
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
    name                         = "default"
    temporary_name_for_rotation  = "temp"
    node_count                   = 1
    vm_size                      = "%s"
    zones                        = %s
    node_public_ip_enabled       = true
    max_pods                     = 60
    only_critical_addons_enabled = true

    kubelet_config {
      pod_max_pid = 12346
    }

    linux_os_config {
      sysctl_config {
        vm_swappiness = 40
      }
    }
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, vmSize, zones)
}

func (KubernetesClusterResource) updateLinuxKernelSettings(data acceptance.TestData) string {
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
    name                         = "default"
    temporary_name_for_rotation  = "temp"
    node_count                   = 1
    vm_size                      = "Standard_D2ads_v5"
    node_public_ip_enabled       = true
    max_pods                     = 60
    only_critical_addons_enabled = true

    kubelet_config {
      pod_max_pid = 12347
    }

    linux_os_config {
      sysctl_config {
        vm_swappiness = 45
      }
    }

    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) enableFips(data acceptance.TestData) string {
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
    fips_enabled = true
    name         = "default"
    node_count   = 1
    vm_size      = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesClusterResource) updateOsDisk(data acceptance.TestData, osDiskType string, osDiskSize int) string {
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
    name                        = "default"
    temporary_name_for_rotation = "temp"
    node_count                  = 1
    os_disk_type                = "%s"
    os_disk_size_gb             = %d
    vm_size                     = "Standard_D2ads_v5"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, osDiskType, osDiskSize)
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
    upgrade_settings {
      max_surge = "10%%"
    }
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
    upgrade_settings {
      max_surge = "10%%"
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
    upgrade_settings {
      max_surge = "10%%"
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
    name                 = "default"
    auto_scaling_enabled = true
    min_count            = 2
    max_count            = 4
    vm_size              = "Standard_DS2_v2"
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
    name                 = "pool1"
    min_count            = 1
    max_count            = 2
    auto_scaling_enabled = true
    vm_size              = "Standard_DS2_v2"
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
    name                 = "pool1"
    min_count            = 1
    max_count            = 2
    auto_scaling_enabled = true
    vm_size              = "Standard_DS2_v2"
    zones                = ["1", "2"]
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, olderKubernetesVersion)
}
func (KubernetesClusterResource) autoScalingProfileConfigMinimal(data acceptance.TestData) string {
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
    name                 = "default"
    auto_scaling_enabled = true
    min_count            = 2
    max_count            = 4
    vm_size              = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  auto_scaler_profile {
    skip_nodes_with_local_storage = false
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, currentKubernetesVersion)
}

func (KubernetesClusterResource) autoScalingProfileConfigComplete(data acceptance.TestData) string {
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
    name                 = "default"
    auto_scaling_enabled = true
    min_count            = 2
    max_count            = 4
    vm_size              = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  auto_scaler_profile {
    balance_similar_node_groups      = true
    expander                         = "least-waste"
    max_graceful_termination_sec     = 15
    max_node_provisioning_time       = "10m"
    max_unready_nodes                = 5
    max_unready_percentage           = 50
    new_pod_scale_up_delay           = "10s"
    scan_interval                    = "10s"
    scale_down_delay_after_add       = "10m"
    scale_down_delay_after_delete    = "10s"
    scale_down_delay_after_failure   = "15m"
    scale_down_unneeded              = "15m"
    scale_down_unready               = "15m"
    scale_down_utilization_threshold = "0.5"
    empty_bulk_delete_max            = "50"
    skip_nodes_with_local_storage    = false
    skip_nodes_with_system_pods      = false
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, currentKubernetesVersion)
}
