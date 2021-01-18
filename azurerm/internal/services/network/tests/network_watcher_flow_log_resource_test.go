package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	azureNetwork "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkWatcherFlowLogResource struct {
}

func testAccNetworkWatcherFlowLog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("0"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_disabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.disabledConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").Exists(),
				check.That(data.ResourceName).Key("retention_policy.0.days").Exists(),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		// disabled flow logs don't import all values
	})
}

func testAccNetworkWatcherFlowLog_reenabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.disabledConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").Exists(),
				check.That(data.ResourceName).Key("retention_policy.0.days").Exists(),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("0"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_retentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("0"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.retentionPolicyConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("7"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_updateStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.retentionPolicyConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("7"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.retentionPolicyConfigUpdateStorageAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("7"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_trafficAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("false"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("0"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.TrafficAnalyticsDisabledConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("7"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.TrafficAnalyticsEnabledConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("7"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("traffic_analytics.#").HasValue("1"),
				check.That(data.ResourceName).Key("traffic_analytics.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("traffic_analytics.0.interval_in_minutes").HasValue("60"),
				check.That(data.ResourceName).Key("traffic_analytics.0.workspace_id").Exists(),
				check.That(data.ResourceName).Key("traffic_analytics.0.workspace_region").Exists(),
				check.That(data.ResourceName).Key("traffic_analytics.0.workspace_resource_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.TrafficAnalyticsUpdateInterval(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("7"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("traffic_analytics.#").HasValue("1"),
				check.That(data.ResourceName).Key("traffic_analytics.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("traffic_analytics.0.interval_in_minutes").HasValue("10"),
				check.That(data.ResourceName).Key("traffic_analytics.0.workspace_id").Exists(),
				check.That(data.ResourceName).Key("traffic_analytics.0.workspace_region").Exists(),
				check.That(data.ResourceName).Key("traffic_analytics.0.workspace_resource_id").Exists(),
			),
		},
		data.ImportStep(),
		// flow log must be disabled before destroy
		{
			Config: r.TrafficAnalyticsDisabledConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_watcher_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("retention_policy.0.days").HasValue("7"),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
	})
}

func testAccNetworkWatcherFlowLog_version(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.versionConfig(data, 1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.versionConfig(data, 2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (t NetworkWatcherFlowLogResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azureNetwork.ParseNetworkWatcherFlowLogID(state.ID)
	if err != nil {
		return nil, err
	}

	// Get current flow log status
	statusParameters := network.FlowLogStatusParameters{
		TargetResourceID: &id.NetworkSecurityGroupID,
	}

	future, err := clients.Network.WatcherClient.GetFlowLogStatus(ctx, id.ResourceGroup, id.NetworkWatcherName, statusParameters)
	if err != nil {
		return nil, fmt.Errorf("reading Network Watcher Flow Log (%s): %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, clients.Network.WatcherClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for retrieval of Flow Log Configuration for target %q: %+v", id, err)
	}

	fli, err := future.Result(*clients.Network.WatcherClient)
	if err != nil {
		return nil, fmt.Errorf("retrieving Flow Log Configuration for target %q: %+v", id, err)
	}

	return utils.Bool(fli.TargetResourceID != nil), nil
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger%1000000)
}

func (r NetworkWatcherFlowLogResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.prerequisites(data))
}

func (r NetworkWatcherFlowLogResource) retentionPolicyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.prerequisites(data))
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

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.testb.id
  enabled                   = true

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.prerequisites(data), data.RandomInteger%1000000+1)
}

func (r NetworkWatcherFlowLogResource) disabledConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = false

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.prerequisites(data))
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
`, r.prerequisites(data), data.RandomInteger)
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
`, r.prerequisites(data), data.RandomInteger)
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
`, r.prerequisites(data), data.RandomInteger)
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
`, r.prerequisites(data), data.RandomInteger, version)
}
