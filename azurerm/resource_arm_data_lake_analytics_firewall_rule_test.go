package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDataLakeAnalyticsFirewallRule_basic(t *testing.T) {
	resourceName := "azurerm_data_lake_analytics_firewall_rule.test"
	ri := tf.AccRandTimeInt()
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsFirewallRule_basic(ri, acceptance.Location(), startIP, endIP),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsFirewallRuleExists(resourceName),
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

func TestAccAzureRMDataLakeAnalyticsFirewallRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_data_lake_analytics_firewall_rule.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsFirewallRule_basic(ri, location, startIP, endIP),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", startIP),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", endIP),
				),
			},
			{
				Config:      testAccAzureRMDataLakeAnalyticsFirewallRule_requiresImport(ri, location, startIP, endIP),
				ExpectError: acceptance.RequiresImportError("azurerm_data_lake_analytics_firewall_rule"),
			},
		},
	})
}

func TestAccAzureRMDataLakeAnalyticsFirewallRule_update(t *testing.T) {
	resourceName := "azurerm_data_lake_analytics_firewall_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsFirewallRule_basic(ri, acceptance.Location(), "1.1.1.1", "2.2.2.2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", "2.2.2.2"),
				),
			},
			{
				Config: testAccAzureRMDataLakeAnalyticsFirewallRule_basic(ri, acceptance.Location(), "2.2.2.2", "3.3.3.3"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsFirewallRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "start_ip_address", "2.2.2.2"),
					resource.TestCheckResourceAttr(resourceName, "end_ip_address", "3.3.3.3"),
				),
			},
		},
	})
}

func TestAccAzureRMDataLakeAnalyticsFirewallRule_azureServices(t *testing.T) {
	resourceName := "azurerm_data_lake_analytics_firewall_rule.test"
	ri := tf.AccRandTimeInt()
	azureServicesIP := "0.0.0.0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataLakeAnalyticsFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataLakeAnalyticsFirewallRule_basic(ri, acceptance.Location(), azureServicesIP, azureServicesIP),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataLakeAnalyticsFirewallRuleExists(resourceName),
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

func testCheckAzureRMDataLakeAnalyticsFirewallRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Datalake.AnalyticsFirewallRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMDataLakeAnalyticsFirewallRuleDestroy(s *terraform.State) error {
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

func testAccAzureRMDataLakeAnalyticsFirewallRule_basic(rInt int, location, startIP, endIP string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_data_lake_analytics_account" "test" {
  name                = "acctest%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  default_store_account_name = "${azurerm_data_lake_store.test.name}"
}

resource "azurerm_data_lake_analytics_firewall_rule" "test" {
  name                = "acctest%[3]s"
  account_name        = "${azurerm_data_lake_analytics_account.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  start_ip_address    = "%[4]s"
  end_ip_address      = "%[5]s"
}
`, rInt, location, strconv.Itoa(rInt)[10:17], startIP, endIP)
}

func testAccAzureRMDataLakeAnalyticsFirewallRule_requiresImport(rInt int, location, startIP, endIP string) string {
	template := testAccAzureRMDataLakeAnalyticsFirewallRule_basic(rInt, location, startIP, endIP)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_analytics_firewall_rule" "import" {
  name                = "${azurerm_data_lake_analytics_firewall_rule.test.name}"
  account_name        = "${azurerm_data_lake_analytics_firewall_rule.test.account_name}"
  resource_group_name = "${azurerm_data_lake_analytics_firewall_rule.test.resource_group_name}"
  start_ip_address    = "${azurerm_data_lake_analytics_firewall_rule.test.start_ip_address}"
  end_ip_address      = "${azurerm_data_lake_analytics_firewall_rule.test.end_ip_address}"
}
`, template)
}
