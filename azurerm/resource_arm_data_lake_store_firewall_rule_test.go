package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"strconv"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDataLakeStoreFirewallRule_basic(t *testing.T) {
	resourceName := "azurerm_data_lake_store_firewall_rule.test"
	ri := acctest.RandInt()
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStoreFirewallRule_basic(ri, testLocation(), startIP, endIP),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", startIP),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", endIP),
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

func TestAccAzureRMDataLakeStoreFirewallRule_update(t *testing.T) {
	resourceName := "azurerm_data_lake_store_firewall_rule.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStoreFirewallRule_basic(ri, testLocation(), "1.1.1.1", "2.2.2.2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", "2.2.2.2"),
				),
			},
			{
				Config: testAccAzureRMDataLakeStoreFirewallRule_basic(ri, testLocation(), "2.2.2.2", "3.3.3.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", "2.2.2.2"),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", "3.3.3.3"),
				),
			},
		},
	})
}

func TestAccAzureRMDataLakeStoreFirewallRule_azureServices(t *testing.T) {
	resourceName := "azurerm_data_lake_store_firewall_rule.test"
	ri := acctest.RandInt()
	azureServicesIP := "0.0.0.0"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataLakeStoreFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeStoreFirewallRule_basic(ri, testLocation(), azureServicesIP, azureServicesIP),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeStoreFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", azureServicesIP),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", azureServicesIP),
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

func testCheckAzureRMDataLakeStoreFirewallRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		firewallRuleName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for data lake store firewall rule: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).dataLakeStoreFirewallRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, accountName, firewallRuleName)
		if err != nil {
			return fmt.Errorf("Bad: Get on dataLakeStoreFirewallRulesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Date Lake Store Firewall Rule %q (Account %q / Resource Group: %q) does not exist", firewallRuleName, accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataLakeStoreFirewallRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dataLakeStoreFirewallRulesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_lake_store_firewall_rule" {
			continue
		}

		firewallRuleName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, accountName, firewallRuleName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Data Lake Store Firewall Rule still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDataLakeStoreFirewallRule_basic(rInt int, location, startIP, endIP string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_data_lake_store_firewall_rule" "test" {
  name                = "acctest"
  account_name        = "${azurerm_data_lake_store.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  start_ip_address    = "%s"
  end_ip_address      = "%s"
}
`, rInt, location, strconv.Itoa(rInt)[0:15], startIP, endIP)
}
