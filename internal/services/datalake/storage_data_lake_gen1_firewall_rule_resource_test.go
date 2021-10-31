package datalake_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageDataLakeGen1FirewallRuleResource struct {
}

func TestAccStorageDataLakeGen1FirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_firewall_rule", "test")
	r := StorageDataLakeGen1FirewallRuleResource{}
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

//

func TestAccStorageDataLakeGen1FirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_firewall_rule", "test")
	r := StorageDataLakeGen1FirewallRuleResource{}
	startIP := "1.1.1.1"
	endIP := "2.2.2.2"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, startIP, endIP),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, startIP, endIP),
			ExpectError: acceptance.RequiresImportError("azurerm_storage_data_lake_gen1_firewall_rule"),
		},
	})
}

func TestAccStorageDataLakeGen1FirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_firewall_rule", "test")
	r := StorageDataLakeGen1FirewallRuleResource{}

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

func TestAccStorageDataLakeGen1FirewallRule_azureServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_firewall_rule", "test")
	r := StorageDataLakeGen1FirewallRuleResource{}
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

func (t StorageDataLakeGen1FirewallRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Datalake.StoreFirewallRulesClient.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Storage Date Lake Gen1 Firewall Rule %s: %v", id, err)
	}

	return utils.Bool(resp.FirewallRuleProperties != nil), nil
}

func (StorageDataLakeGen1FirewallRuleResource) basic(data acceptance.TestData, startIP, endIP string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_storage_data_lake_gen1_firewall_rule" "test" {
  name                = "acctest"
  account_name        = azurerm_storage_data_lake_gen1_filesystem.test.name
  resource_group_name = azurerm_resource_group.test.name
  start_ip_address    = "%s"
  end_ip_address      = "%s"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17], startIP, endIP)
}

func (StorageDataLakeGen1FirewallRuleResource) requiresImport(data acceptance.TestData, startIP, endIP string) string {
	template := StorageDataLakeGen1FirewallRuleResource{}.basic(data, startIP, endIP)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen1_firewall_rule" "import" {
  name                = azurerm_storage_data_lake_gen1_firewall_rule.test.name
  account_name        = azurerm_storage_data_lake_gen1_firewall_rule.test.account_name
  resource_group_name = azurerm_storage_data_lake_gen1_firewall_rule.test.resource_group_name
  start_ip_address    = azurerm_storage_data_lake_gen1_firewall_rule.test.start_ip_address
  end_ip_address      = azurerm_storage_data_lake_gen1_firewall_rule.test.end_ip_address
}
`, template)
}
