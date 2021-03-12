package datalake_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataLakeStoreFirewallRuleResource struct {
}

func TestAccDataLakeStoreFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_firewall_rule", "test")
	r := DataLakeStoreFirewallRuleResource{}
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, startIP, endIP),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue(startIP),
				check.That(data.ResourceName).Key("end_ip_address").HasValue(endIP),
			),
		},
		data.ImportStep(),
	})
}

//

func TestAccDataLakeStoreFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_firewall_rule", "test")
	r := DataLakeStoreFirewallRuleResource{}
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, startIP, endIP),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, startIP, endIP),
			ExpectError: acceptance.RequiresImportError("azurerm_data_lake_store_firewall_rule"),
		},
	})
}

func TestAccDataLakeStoreFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_firewall_rule", "test")
	r := DataLakeStoreFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "1.1.1.1", "2.2.2.2"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue("1.1.1.1"),
				check.That(data.ResourceName).Key("end_ip_address").HasValue("2.2.2.2"),
			),
		},
		{
			Config: r.basic(data, "2.2.2.2", "3.3.3.3"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue("2.2.2.2"),
				check.That(data.ResourceName).Key("end_ip_address").HasValue("3.3.3.3"),
			),
		},
	})
}

func TestAccDataLakeStoreFirewallRule_azureServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_lake_store_firewall_rule", "test")
	r := DataLakeStoreFirewallRuleResource{}
	azureServicesIP := "0.0.0.0"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, azureServicesIP, azureServicesIP),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("start_ip_address").HasValue(azureServicesIP),
				check.That(data.ResourceName).Key("end_ip_address").HasValue(azureServicesIP),
			),
		},
		data.ImportStep(),
	})
}

func (t DataLakeStoreFirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["accounts"]
	name := id.Path["firewallRules"]

	resp, err := clients.Datalake.StoreFirewallRulesClient.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Date Lake Store Firewall Rule %q (Account %q / Resource Group: %q): %v", name, accountName, resourceGroup, err)
	}

	return utils.Bool(resp.FirewallRuleProperties != nil), nil
}

func (DataLakeStoreFirewallRuleResource) basic(data acceptance.TestData, startIP, endIP string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_data_lake_store_firewall_rule" "test" {
  name                = "acctest"
  account_name        = azurerm_data_lake_store.test.name
  resource_group_name = azurerm_resource_group.test.name
  start_ip_address    = "%s"
  end_ip_address      = "%s"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17], startIP, endIP)
}

func (DataLakeStoreFirewallRuleResource) requiresImport(data acceptance.TestData, startIP, endIP string) string {
	template := DataLakeStoreFirewallRuleResource{}.basic(data, startIP, endIP)
	return fmt.Sprintf(`
%s

resource "azurerm_data_lake_store_firewall_rule" "import" {
  name                = azurerm_data_lake_store_firewall_rule.test.name
  account_name        = azurerm_data_lake_store_firewall_rule.test.account_name
  resource_group_name = azurerm_data_lake_store_firewall_rule.test.resource_group_name
  start_ip_address    = azurerm_data_lake_store_firewall_rule.test.start_ip_address
  end_ip_address      = azurerm_data_lake_store_firewall_rule.test.end_ip_address
}
`, template)
}
