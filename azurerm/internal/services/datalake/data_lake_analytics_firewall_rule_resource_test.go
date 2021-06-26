package datalake_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataLakeAnalyticsFirewallRuleResource struct {
}

func TestAccDataLakeAnalyticsFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_firewall_rule", "test")
	r := DataLakeAnalyticsFirewallRuleResource{}
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, startIP, endIP),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue(startIP),
				check.That(data.ResourceName).Key("end_ip_address").HasValue(endIP),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataLakeAnalyticsFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_firewall_rule", "test")
	r := DataLakeAnalyticsFirewallRuleResource{}
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, startIP, endIP),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue(startIP),
				check.That(data.ResourceName).Key("end_ip_address").HasValue(endIP),
			),
		},
		{
			Config:      r.requiresImport(data, startIP, endIP),
			ExpectError: acceptance.RequiresImportError("azurerm_data_lake_analytics_firewall_rule"),
		},
	})
}

func TestAccDataLakeAnalyticsFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_firewall_rule", "test")
	r := DataLakeAnalyticsFirewallRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "1.1.1.1", "2.2.2.2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue("1.1.1.1"),
				check.That(data.ResourceName).Key("end_ip_address").HasValue("2.2.2.2"),
			),
		},
		{
			Config: r.basic(data, "2.2.2.2", "3.3.3.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue("2.2.2.2"),
				check.That(data.ResourceName).Key("end_ip_address").HasValue("3.3.3.3"),
			),
		},
	})
}

func TestAccDataLakeAnalyticsFirewallRule_azureServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_analytics_firewall_rule", "test")
	r := DataLakeAnalyticsFirewallRuleResource{}
	azureServicesIP := "0.0.0.0"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, azureServicesIP, azureServicesIP),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue(azureServicesIP),
				check.That(data.ResourceName).Key("end_ip_address").HasValue(azureServicesIP),
			),
		},
		data.ImportStep(),
	})
}

func (t DataLakeAnalyticsFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	accountName := id.Path["accounts"]
	name := id.Path["firewallRules"]

	resp, err := clients.Datalake.AnalyticsFirewallRulesClient.Get(ctx, id.ResourceGroup, accountName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Date Lake Analytics Firewall Rule %q (Account %q / Resource Group: %q): %v", name, accountName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.FirewallRuleProperties != nil), nil
}

func (DataLakeAnalyticsFirewallRuleResource) basic(data acceptance.TestData, startIP, endIP string) string {
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

func (DataLakeAnalyticsFirewallRuleResource) requiresImport(data acceptance.TestData, startIP, endIP string) string {
	template := DataLakeAnalyticsFirewallRuleResource{}.basic(data, startIP, endIP)
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
