package datalake_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccDataLakeAnalyticsFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_firewall_rule", "test")
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataLakeAnalyticsFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLakeAnalyticsFirewallRule_basic(data, startIP, endIP),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataLakeAnalyticsFirewallRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "start_ip_address", startIP),
					resource.TestCheckResourceAttr(data.ResourceName, "end_ip_address", endIP),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataLakeAnalyticsFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_firewall_rule", "test")
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataLakeAnalyticsFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLakeAnalyticsFirewallRule_basic(data, startIP, endIP),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataLakeAnalyticsFirewallRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "start_ip_address", startIP),
					resource.TestCheckResourceAttr(data.ResourceName, "end_ip_address", endIP),
				),
			},
			{
				Config:      testAccDataLakeAnalyticsFirewallRule_requiresImport(data, startIP, endIP),
				ExpectError: acceptance.RequiresImportError("azurerm_data_lake_analytics_firewall_rule"),
			},
		},
	})
}

func TestAccDataLakeAnalyticsFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_firewall_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataLakeAnalyticsFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLakeAnalyticsFirewallRule_basic(data, "1.1.1.1", "2.2.2.2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataLakeAnalyticsFirewallRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "start_ip_address", "1.1.1.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "end_ip_address", "2.2.2.2"),
				),
			},
			{
				Config: testAccDataLakeAnalyticsFirewallRule_basic(data, "2.2.2.2", "3.3.3.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataLakeAnalyticsFirewallRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "start_ip_address", "2.2.2.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "end_ip_address", "3.3.3.3"),
				),
			},
		},
	})
}

func TestAccDataLakeAnalyticsFirewallRule_azureServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_firewall_rule", "test")
	azureServicesIP := "0.0.0.0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataLakeAnalyticsFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLakeAnalyticsFirewallRule_basic(data, azureServicesIP, azureServicesIP),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataLakeAnalyticsFirewallRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "start_ip_address", azureServicesIP),
					resource.TestCheckResourceAttr(data.ResourceName, "end_ip_address", azureServicesIP),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckDataLakeAnalyticsFirewallRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.AnalyticsFirewallRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		firewallRuleName := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for data lake store firewall rule: %s", firewallRuleName)
		}

		resp, err := conn.Get(ctx, resourceGroup, accountName, firewallRuleName)
		if err != nil {
			return fmt.Errorf("Bad: Get on dataLakeAnalyticsFirewallRulesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Date Lake Analytics Firewall Rule %q (Account %q / Resource Group: %q) does not exist", firewallRuleName, accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckDataLakeAnalyticsFirewallRuleDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.AnalyticsFirewallRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_lake_analytics_firewall_rule" {
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

		return fmt.Errorf("Data Lake Analytics Firewall Rule still exists:\n%#v", resp)
	}

	return nil
}

func testAccDataLakeAnalyticsFirewallRule_basic(data acceptance.TestData, startIP, endIP string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  default_store_account_name = azurerm_data_lake_store.test.name
}

resource "azurerm_data_lake_analytics_firewall_rule" "test" {
  name                = "acctest%[3]s"
  account_name        = azurerm_data_lake_analytics_account.test.name
  resource_group_name = azurerm_resource_group.test.name
  start_ip_address    = "%[4]s"
  end_ip_address      = "%[5]s"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[10:17], startIP, endIP)
}

func testAccDataLakeAnalyticsFirewallRule_requiresImport(data acceptance.TestData, startIP, endIP string) string {
	template := testAccDataLakeAnalyticsFirewallRule_basic(data, startIP, endIP)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_firewall_rule" "import" {
  name                = azurerm_data_lake_analytics_firewall_rule.test.name
  account_name        = azurerm_data_lake_analytics_firewall_rule.test.account_name
  resource_group_name = azurerm_data_lake_analytics_firewall_rule.test.resource_group_name
  start_ip_address    = azurerm_data_lake_analytics_firewall_rule.test.start_ip_address
  end_ip_address      = azurerm_data_lake_analytics_firewall_rule.test.end_ip_address
}
`, template)
}
