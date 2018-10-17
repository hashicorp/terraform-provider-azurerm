package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMNetworkWatcherFlowLog(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning one per region at once
	// (which our test suite can't easily workaround)
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"basic":                testAccAzureRMNetworkWatcherFlowLog_basic,
			"disabled":             testAccAzureRMNetworkWatcherFlowLog_disabled,
			"reenabled":            testAccAzureRMNetworkWatcherFlowLog_reenabled,
			"retentionPolicy":      testAccAzureRMNetworkWatcherFlowLog_retentionPolicy,
			"updateStorageAccount": testAccAzureRMNetworkWatcherFlowLog_updateStorageAccount,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccAzureRMNetworkWatcherFlowLog_basic(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(8)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_disabled(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(8)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_disabledConfig(ri, rs, location),
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
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_reenabled(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(8)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_disabledConfig(ri, rs, location),
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
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_retentionPolicy(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(8)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_basicConfig(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccAzureRMNetworkWatcherFlowLog_updateStorageAccount(t *testing.T) {
	resourceName := "azurerm_network_watcher_flow_log.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(8)
	rsNew := acctest.RandString(8)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkWatcherDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(ri, rsNew, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkWatcherFlowLogExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "network_watcher_name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "network_security_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "retention_policy.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "retention_policy.0.days"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testCheckAzureRMNetworkWatcherFlowLogExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		networkSecurityGroupID := rs.Primary.Attributes["network_security_group_id"]
		parsedID, err := parseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		resourceGroupName := parsedID.ResourceGroup
		networkWatcherName := parsedID.Path["networkWatchers"]

		client := testAccProvider.Meta().(*ArmClient).watcherClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		statusParameters := network.FlowLogStatusParameters{
			TargetResourceID: &networkSecurityGroupID,
		}
		future, err := client.GetFlowLogStatus(ctx, resourceGroupName, networkWatcherName, statusParameters)
		if err != nil {
			return fmt.Errorf("Error retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
		}

		err = future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for retrieval of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
		}

		_, err = future.Result(client)
		if err != nil {
			return fmt.Errorf("Error retrieving of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
		}

		return nil
	}
}

func testAccAzureRMNetworkWatcherFlowLog_basicConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
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
    name                = "acctestsa%s"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location            = "${azurerm_resource_group.test.location}"

    account_tier              = "Standard"
    account_replication_type  = "LRS"
    enable_https_traffic_only = true
}

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.test.id}"
    enabled                   = true

    retention_policy {
        days    = 0
        enabled = false
    }
}
`, rInt, location, rInt, rInt, rString)
}

func testAccAzureRMNetworkWatcherFlowLog_retentionPolicyConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
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
    name                = "acctestsa%s"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location            = "${azurerm_resource_group.test.location}"

    account_tier              = "Standard"
    account_replication_type  = "LRS"
    enable_https_traffic_only = true
}

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.test.id}"
    enabled                   = true
    
    retention_policy {
        days    = 7
        enabled = true
    }
}
`, rInt, location, rInt, rInt, rString)
}

func testAccAzureRMNetworkWatcherFlowLog_disabledConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name     = "acctestRG-%d"
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
    name                = "acctestsa%s"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location            = "${azurerm_resource_group.test.location}"

    account_tier              = "Standard"
    account_replication_type  = "LRS"
    enable_https_traffic_only = true
}

resource "azurerm_network_watcher_flow_log" "test" {
    network_watcher_name = "${azurerm_network_watcher.test.name}"
    resource_group_name  = "${azurerm_resource_group.test.name}"

    network_security_group_id = "${azurerm_network_security_group.test.id}"
    storage_account_id        = "${azurerm_storage_account.test.id}"
    enabled                   = false
    
    retention_policy {
        days    = 7
        enabled = true
    }
}
`, rInt, location, rInt, rInt, rString)
}
