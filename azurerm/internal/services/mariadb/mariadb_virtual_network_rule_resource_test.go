package mariadb_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MariaDbVirtualNetworkRuleResource struct {
}

func TestAccMariaDbVirtualNetworkRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_virtual_network_rule", "test")
	r := MariaDbVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccMariaDbVirtualNetworkRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_virtual_network_rule", "test")
	r := MariaDbVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_mariadb_virtual_network_rule"),
		},
	})
}

func TestAccMariaDbVirtualNetworkRule_switchSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_virtual_network_rule", "test")
	r := MariaDbVirtualNetworkRuleResource{}

	// Create regex strings that will ensure that one subnet name exists, but not the other
	preConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet1%d)$|(subnet[^2]%d)$", data.RandomInteger, data.RandomInteger))  // subnet 1 but not 2
	postConfigRegex := regexp.MustCompile(fmt.Sprintf("(subnet2%d)$|(subnet[^1]%d)$", data.RandomInteger, data.RandomInteger)) // subnet 2 but not 1

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subnetSwitchPre(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestMatchResourceAttr(data.ResourceName, "subnet_id", preConfigRegex),
			),
		},
		{
			Config: r.subnetSwitchPost(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestMatchResourceAttr(data.ResourceName, "subnet_id", postConfigRegex),
			),
		},
	})
}

func TestAccMariaDbVirtualNetworkRule_multipleSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mariadb_virtual_network_rule", "rule1")
	r := MariaDbVirtualNetworkRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleSubnets(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_mariadb_virtual_network_rule.rule2").ExistsInAzure(r),
				check.That("azurerm_mariadb_virtual_network_rule.rule3").ExistsInAzure(r),
			),
		},
	})
}

func (MariaDbVirtualNetworkRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	serverName := id.Path["servers"]
	name := id.Path["virtualNetworkRules"]

	resp, err := clients.MariaDB.VirtualNetworkRulesClient.Get(ctx, id.ResourceGroup, serverName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving MariaDB Virtual Network Rule %q (Server %q / Resource Group %q): %v", name, serverName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.VirtualNetworkRuleProperties != nil), nil
}

func (MariaDbVirtualNetworkRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mariadb_server" "test" {
  name                         = "acctestmariadbsvr-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement_enabled      = true
  sku_name                     = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_mariadb_virtual_network_rule" "test" {
  name                = "acctestmariadbvnetrule%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mariadb_server.test.name
  subnet_id           = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r MariaDbVirtualNetworkRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mariadb_virtual_network_rule" "import" {
  name                = azurerm_mariadb_virtual_network_rule.test.name
  resource_group_name = azurerm_mariadb_virtual_network_rule.test.resource_group_name
  server_name         = azurerm_mariadb_virtual_network_rule.test.server_name
  subnet_id           = azurerm_mariadb_virtual_network_rule.test.subnet_id
}
`, r.basic(data))
}

func (MariaDbVirtualNetworkRuleResource) subnetSwitchPre(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mariadb_server" "test" {
  name                         = "acctestmariadbsvr-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement_enabled      = true
  sku_name                     = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_mariadb_virtual_network_rule" "test" {
  name                = "acctestmariadbvnetrule%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mariadb_server.test.name
  subnet_id           = azurerm_subnet.test1.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (MariaDbVirtualNetworkRuleResource) subnetSwitchPost(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "subnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "subnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mariadb_server" "test" {
  name                         = "acctestmariadbsvr-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement_enabled      = true
  sku_name                     = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_mariadb_virtual_network_rule" "test" {
  name                = "acctestmariadbvnetrule%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mariadb_server.test.name
  subnet_id           = azurerm_subnet.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (MariaDbVirtualNetworkRuleResource) multipleSubnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "vnet1" {
  name                = "acctestvnet1%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "vnet2" {
  name                = "acctestvnet2%d"
  address_space       = ["10.1.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "vnet1_subnet1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet1.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "vnet1_subnet2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet1.name
  address_prefix       = "10.7.29.128/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "vnet2_subnet1" {
  name                 = "acctestsubnet3%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.vnet2.name
  address_prefix       = "10.1.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_mariadb_server" "test" {
  name                         = "acctestmariadbsvr-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement_enabled      = true
  sku_name                     = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_mariadb_virtual_network_rule" "rule1" {
  name                = "acctestmariadbvnetrule1%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mariadb_server.test.name
  subnet_id           = azurerm_subnet.vnet1_subnet1.id
}

resource "azurerm_mariadb_virtual_network_rule" "rule2" {
  name                = "acctestmariadbvnetrule2%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mariadb_server.test.name
  subnet_id           = azurerm_subnet.vnet1_subnet2.id
}

resource "azurerm_mariadb_virtual_network_rule" "rule3" {
  name                = "acctestmariadbvnetrule3%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_mariadb_server.test.name
  subnet_id           = azurerm_subnet.vnet2_subnet1.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
