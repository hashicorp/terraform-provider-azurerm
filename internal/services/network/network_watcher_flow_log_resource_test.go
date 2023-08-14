// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/flowlogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkWatcherFlowLogResource struct{}

func testAccNetworkWatcherFlowLog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_network_watcher_flow_log"),
		},
	})
}

func testAccNetworkWatcherFlowLog_disabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.disabledConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_reenabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.disabledConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_retentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.retentionPolicyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_updateStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.retentionPolicyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.retentionPolicyConfigUpdateStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_trafficAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.TrafficAnalyticsDisabledConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.TrafficAnalyticsEnabledConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.TrafficAnalyticsUpdateInterval(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// flow log must be disabled before destroy
		{
			Config: r.TrafficAnalyticsDisabledConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func testAccNetworkWatcherFlowLog_version(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.versionConfig(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.versionConfig(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_location(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.location(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data, "Test"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags(data, "Prod"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkWatcherFlowLogResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := flowlogs.ParseFlowLogID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.FlowLogs.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Network Watcher Flow Log (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (NetworkWatcherFlowLogResource) prerequisites(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-watcher-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestNSG%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_watcher" "test" {
  name                = "acctest-NW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsa%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  account_tier              = "Standard"
  account_kind              = "StorageV2"
  account_replication_type  = "LRS"
  enable_https_traffic_only = true
}
`, data.RandomIntOfLength(10), data.Locations.Primary, data.RandomIntOfLength(10), data.RandomInteger, data.RandomInteger%1000000)
}

func (r NetworkWatcherFlowLogResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.prerequisites(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "import" {
  network_watcher_name = azurerm_network_watcher_flow_log.test.network_watcher_name
  resource_group_name  = azurerm_network_watcher_flow_log.test.resource_group_name
  name                 = azurerm_network_watcher_flow_log.test.name

  network_security_group_id = azurerm_network_watcher_flow_log.test.network_security_group_id
  storage_account_id        = azurerm_network_watcher_flow_log.test.storage_account_id
  enabled                   = azurerm_network_watcher_flow_log.test.enabled

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.basicConfig(data))
}

func (r NetworkWatcherFlowLogResource) retentionPolicyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.prerequisites(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) retentionPolicyConfigUpdateStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "testb" {
  name                = "acctestsab%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  account_tier              = "Standard"
  account_kind              = "StorageV2"
  account_replication_type  = "LRS"
  enable_https_traffic_only = true
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.testb.id
  enabled                   = true

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.prerequisites(data), data.RandomInteger%1000000+1, data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) disabledConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"
  location             = azurerm_network_watcher.test.location

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = false

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.prerequisites(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) TrafficAnalyticsEnabledConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"
  location             = azurerm_network_watcher.test.location

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = true
    days    = 7
  }

  traffic_analytics {
    enabled               = true
    workspace_id          = azurerm_log_analytics_workspace.test.workspace_id
    workspace_region      = azurerm_log_analytics_workspace.test.location
    workspace_resource_id = azurerm_log_analytics_workspace.test.id
  }
}
`, r.prerequisites(data), data.RandomInteger, data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) TrafficAnalyticsUpdateInterval(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"
  location             = azurerm_network_watcher.test.location

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = true
    days    = 7
  }

  traffic_analytics {
    enabled               = true
    workspace_id          = azurerm_log_analytics_workspace.test.workspace_id
    workspace_region      = azurerm_log_analytics_workspace.test.location
    workspace_resource_id = azurerm_log_analytics_workspace.test.id
    interval_in_minutes   = 10
  }
}
`, r.prerequisites(data), data.RandomInteger, data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) TrafficAnalyticsDisabledConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"
  location             = azurerm_network_watcher.test.location

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = true
    days    = 7
  }

  traffic_analytics {
    enabled               = false
    workspace_id          = azurerm_log_analytics_workspace.test.workspace_id
    workspace_region      = azurerm_log_analytics_workspace.test.location
    workspace_resource_id = azurerm_log_analytics_workspace.test.id
  }
}
`, r.prerequisites(data), data.RandomInteger, data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) versionConfig(data acceptance.TestData, version int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"
  location             = azurerm_network_watcher.test.location

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true
  version                   = %d

  retention_policy {
    enabled = true
    days    = 7
  }

  traffic_analytics {
    enabled               = true
    workspace_id          = azurerm_log_analytics_workspace.test.workspace_id
    workspace_region      = azurerm_log_analytics_workspace.test.location
    workspace_resource_id = azurerm_log_analytics_workspace.test.id
  }
}
`, r.prerequisites(data), data.RandomInteger, data.RandomInteger, version)
}

func (r NetworkWatcherFlowLogResource) location(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"
  location             = azurerm_resource_group.test.location

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.prerequisites(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) tags(data acceptance.TestData, v string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = false
    days    = 0
  }

  tags = {
    env = "%s"
  }
}
`, r.prerequisites(data), data.RandomInteger, v)
}
