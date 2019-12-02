package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-07-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func testAccAzureRMNetworkWatcherFlowLog_basic(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_reenabled(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

func testAccAzureRMNetworkWatcherFlowLog_trafficAnalytics(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
		},
	})
}

func testCheckAzureRMNetworkWatcherFlowLogExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).Network.WatcherClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		parsedID, err := parseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		resourceGroupName := parsedID.ResourceGroup
		networkWatcherName := parsedID.Path["networkWatchers"]
		networkSecurityGroupID := parsedID.Path["networkSecurityGroupId"]

		statusParameters := network.FlowLogStatusParameters{
			TargetResourceID: &networkSecurityGroupID,
		}
		future, err := client.GetFlowLogStatus(ctx, resourceGroupName, networkWatcherName, statusParameters)
		if err != nil {
			return fmt.Errorf("Error retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for retrieval of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
		}

		if _, err := future.Result(*client); err != nil {
			return fmt.Errorf("Error retrieving of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
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
    name                = "acctestnsg%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_watcher" "test" {
    name                = "acctestnw-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_storage_account" "test" {
    name                = "acctestsa%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location            = "${azurerm_resource_group.test.location}"

    account_tier              = "Standard"
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
    name                = "acctestlaw-%d"
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
    name                = "acctestlaw-%d"
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
