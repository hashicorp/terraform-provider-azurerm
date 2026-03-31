// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/flowlogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkWatcherFlowLogResource struct{}

func testAccNetworkWatcherFlowLog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// turns out while you can't create with an NSG, you can create a new flow log
			// and update it to an NSG, so we'll include that in the test for now ¯\_(ツ)_/¯
			Config: r.basicWithNSG(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_cannotCreateNewWithNSG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basicWithNSG(data),
			ExpectError: regexp.MustCompile("creation of new NSG flow logs is no longer supported by Azure as of June 30, 2025."),
		},
	})
}

func testAccNetworkWatcherFlowLog_basicWithSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfigWithSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkWatcherFlowLog_basicWithNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_watcher_flow_log", "test")
	r := NetworkWatcherFlowLogResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfigWithNIC(data),
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
			Config: r.basic(data),
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
			Config: r.basic(data),
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
			Config: r.basic(data),
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
			Config: r.basic(data),
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
		data.ImportStep(),
		{
			Config: r.basicIgnoreUnmanaged(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
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

func (r NetworkWatcherFlowLogResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := flowlogs.ParseFlowLogID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.FlowLogs.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (NetworkWatcherFlowLogResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-watcher-%[1]d"
  location = "%[2]s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestNSG%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_watcher" "test" {
  name                = "acctest-NW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsa%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  account_tier               = "Standard"
  account_kind               = "StorageV2"
  account_replication_type   = "LRS"
  https_traffic_only_enabled = true
}
`, data.RandomIntOfLength(10), data.Locations.Primary, data.RandomInteger%1000000)
}

func (r NetworkWatcherFlowLogResource) basicWithNSG(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%[2]d"

  network_security_group_id = azurerm_network_security_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  enabled                   = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.template(data), data.RandomInteger)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%[2]d"

  target_resource_id = azurerm_network_security_group.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%[2]d"

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) basicIgnoreUnmanaged(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      # "azurerm_log_analytics_workspace" leaves behind unmanaged resources created by Azure
      prevent_deletion_if_contains_resources = false
    }
  }
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%[2]d"

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) basicConfigWithSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%[2]d"

  target_resource_id = azurerm_subnet.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) basicConfigWithNIC(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "T@rraf0rmTest"
  disable_password_authentication = false

  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%[2]d"

  target_resource_id = azurerm_network_interface.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "import" {
  network_watcher_name = azurerm_network_watcher_flow_log.test.network_watcher_name
  resource_group_name  = azurerm_network_watcher_flow_log.test.resource_group_name
  name                 = azurerm_network_watcher_flow_log.test.name

  target_resource_id = azurerm_network_watcher_flow_log.test.target_resource_id
  storage_account_id = azurerm_network_watcher_flow_log.test.storage_account_id
  enabled            = azurerm_network_watcher_flow_log.test.enabled

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.basic(data))
}

func (r NetworkWatcherFlowLogResource) retentionPolicyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) retentionPolicyConfigUpdateStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_account" "testb" {
  name                = "acctestsab%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  account_tier               = "Standard"
  account_kind               = "StorageV2"
  account_replication_type   = "LRS"
  https_traffic_only_enabled = true
}

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.testb.id
  enabled            = true

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.template(data), data.RandomInteger%1000000+1, data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) disabledConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"
  location             = azurerm_network_watcher.test.location

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = false

  retention_policy {
    enabled = true
    days    = 7
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) TrafficAnalyticsEnabledConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      # "azurerm_log_analytics_workspace" leaves behind unmanaged resources created by Azure
      prevent_deletion_if_contains_resources = false
    }
  }
}

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

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

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
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) TrafficAnalyticsUpdateInterval(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      # "azurerm_log_analytics_workspace" leaves behind unmanaged resources created by Azure
      prevent_deletion_if_contains_resources = false
    }
  }
}

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

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

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
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) TrafficAnalyticsDisabledConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      # "azurerm_log_analytics_workspace" leaves behind unmanaged resources created by Azure
      prevent_deletion_if_contains_resources = false
    }
  }
}

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

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

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
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) versionConfig(data acceptance.TestData, version int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      # "azurerm_log_analytics_workspace" leaves behind unmanaged resources created by Azure
      prevent_deletion_if_contains_resources = false
    }
  }
}

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

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true
  version            = %d

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
`, r.template(data), data.RandomInteger, data.RandomInteger, version)
}

func (r NetworkWatcherFlowLogResource) location(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"
  location             = azurerm_resource_group.test.location

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = false
    days    = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetworkWatcherFlowLogResource) tags(data acceptance.TestData, v string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_network_watcher_flow_log" "test" {
  network_watcher_name = azurerm_network_watcher.test.name
  resource_group_name  = azurerm_resource_group.test.name
  name                 = "flowlog-%d"

  target_resource_id = azurerm_virtual_network.test.id
  storage_account_id = azurerm_storage_account.test.id
  enabled            = true

  retention_policy {
    enabled = false
    days    = 0
  }

  tags = {
    env = "%s"
  }
}
`, r.template(data), data.RandomInteger, v)
}
