package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	nw "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
)

func testAccAzureRMNetworkWatcherFlowLog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_disabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_disabledConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			// disabled flow logs don't import all values
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_reenabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_disabledConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_retentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_updateStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfigUpdateStorageAccount(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_trafficAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsDisabledConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsEnabledConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_analytics.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_analytics.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_analytics.0.interval_in_minutes", "60"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "traffic_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "traffic_analytics.0.workspace_region"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "traffic_analytics.0.workspace_resource_id"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsUpdateInterval(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_analytics.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_analytics.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_analytics.0.interval_in_minutes", "10"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "traffic_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "traffic_analytics.0.workspace_region"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "traffic_analytics.0.workspace_resource_id"),
				),
			},
			data.ImportStep(),
			// flow log must be disabled before destroy
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsDisabledConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_version(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_versionConfig(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_versionConfig(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetworkWatcherFlowLogExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.WatcherClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := nw.ParseNetworkWatcherFlowLogID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		statusParameters := network.FlowLogStatusParameters{
			TargetResourceID: &id.NetworkSecurityGroupID,
		}
		future, err := client.GetFlowLogStatus(ctx, id.ResourceGroup, id.NetworkWatcherName, statusParameters)
		if err != nil {
			return fmt.Errorf("Error retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for retrieval of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
		}

		if _, err := future.Result(*client); err != nil {
			return fmt.Errorf("Error retrieving of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
		}

		return nil
	}
}

func testAccAzureRMNetworkWatcherFlowLog_prerequisites(data acceptance.TestData) string {
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

func testAccAzureRMNetworkWatcherFlowLog_basicConfig(data acceptance.TestData) string {
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
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(data))
}

func testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(data acceptance.TestData) string {
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
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(data))
}

func testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfigUpdateStorageAccount(data acceptance.TestData) string {
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
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(data), data.RandomInteger%1000000+1)
}

func testAccAzureRMNetworkWatcherFlowLog_disabledConfig(data acceptance.TestData) string {
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
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(data))
}

func testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsEnabledConfig(data acceptance.TestData) string {
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
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(data), data.RandomInteger)
}

func testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsUpdateInterval(data acceptance.TestData) string {
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
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(data), data.RandomInteger)
}

func testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsDisabledConfig(data acceptance.TestData) string {
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
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(data), data.RandomInteger)
}

func testAccAzureRMNetworkWatcherFlowLog_versionConfig(data acceptance.TestData, version int) string {
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
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(data), data.RandomInteger, version)
}
