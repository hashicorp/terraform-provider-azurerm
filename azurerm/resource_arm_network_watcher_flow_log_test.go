package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func testAccAzureRMNetworkWatcherFlowLog_basic(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "0"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_disabled(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_disabledConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			// disabled flow logs don't import all values
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_reenabled(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_disabledConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "0"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_retentionPolicy(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "0"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_updateStorageAccount(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfigUpdateStorageAccount(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_trafficAnalytics(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "0"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsDisabledConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsEnabledConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "traffic_analytics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_analytics.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "traffic_analytics.0.workspace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "traffic_analytics.0.workspace_region"),
					resource.TestCheckResourceAttrSet(resourceName, "traffic_analytics.0.workspace_resource_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// flow log must be disabled before destroy
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsDisabledConfig(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.days", "7"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
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

		id, err := ParseNetworkWatcherFlowLogID(rs.Primary.Attributes["id"])
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

func testAccAzureRMNetworkWatcherFlowLog_prerequisites(rInt int, location string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-watcher-%d"
    location = "%s"
}

resource "azurerm_network_security_group" "test" {
    name                = "acctestNSG%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_watcher" "test" {
    name                = "acctest-NW-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_storage_account" "test" {
    name                = "acctestsa%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location            = "${azurerm_resource_group.test.location}"

    account_tier              = "Standard"
    account_kind              = "StorageV2"
    account_replication_type  = "LRS"
    enable_https_traffic_only = true
}
`, rInt, location, rInt, rInt, rInt%1000000)
}

func testAccAzureRMNetworkWatcherFlowLog_basicConfig(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.test.id}"
    enabled                   = true

    retention_policy {
        enabled = false
        days    = 0
    }
}
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(rInt, location))
}

func testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.test.id}"
    enabled                   = true
    
    retention_policy {
        enabled = true
        days    = 7
    }
}
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(rInt, location))
}

func testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfigUpdateStorageAccount(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "testb" {
    name                = "acctestsab%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location            = "${azurerm_resource_group.test.location}"

    account_tier              = "Standard"
    account_kind              = "StorageV2"
    account_replication_type  = "LRS"
    enable_https_traffic_only = true
}

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.testb.id}"
    enabled                   = true
    
    retention_policy {
        enabled = true
        days    = 7
    }
}
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(rInt, location), rInt%1000000+1)
}

func testAccAzureRMNetworkWatcherFlowLog_disabledConfig(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.test.id}"
    enabled                   = false
    
    retention_policy {
        enabled = true
        days    = 7
    }
}
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(rInt, location))
}

func testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsEnabledConfig(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
    name                = "acctestLAW-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.test.id}"
    enabled                   = true
    
    retention_policy {
        enabled = true
        days    = 7
    }

    traffic_analytics {
        enabled               = true
        workspace_id          = "${azurerm_log_analytics_workspace.test.workspace_id}"
        workspace_region      = "${azurerm_log_analytics_workspace.test.location}"
        workspace_resource_id = "${azurerm_log_analytics_workspace.test.id}"
    }
}
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(rInt, location), rInt)
}

func testAccAzureRMNetworkWatcherFlowLog_TrafficAnalyticsDisabledConfig(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
    name                = "acctestLAW-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku                 = "PerGB2018"
}

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.test.id}"
    enabled                   = true
    
    retention_policy {
        enabled = true
        days    = 7
    }

    traffic_analytics {
        enabled               = false
        workspace_id          = "${azurerm_log_analytics_workspace.test.workspace_id}"
        workspace_region      = "${azurerm_log_analytics_workspace.test.location}"
        workspace_resource_id = "${azurerm_log_analytics_workspace.test.id}"
    }
}
`, testAccAzureRMNetworkWatcherFlowLog_prerequisites(rInt, location), rInt)
}
