package postgres_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PostgreSQLVirtualNetworkRuleResource struct {
}

func TestAccPostgreSQLVirtualNetworkRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	r := PostgreSQLVirtualNetworkRuleResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPostgreSQLVirtualNetworkRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	r := PostgreSQLVirtualNetworkRuleResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPostgreSQLVirtualNetworkRule_switchSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	r := PostgreSQLVirtualNetworkRuleResource{}

	// Create regex strings that will ensure that one subnet name exists, but not the other
	preConfigRegex := regexp.MustCompile(fmt.Sprintf("(acctest-SN1-%d)$|(acctest-SN[^2]-%d)$", data.RandomInteger, data.RandomInteger))  // subnet 1 but not 2
	postConfigRegex := regexp.MustCompile(fmt.Sprintf("(acctest-SN2-%d)$|(acctest-SN-[^1]%d)$", data.RandomInteger, data.RandomInteger)) // subnet 2 but not 1

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subnetSwitchPre(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestMatchResourceAttr(data.ResourceName, "subnet_id", preConfigRegex),
			),
		},
		{
			Config: r.subnetSwitchPost(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestMatchResourceAttr(data.ResourceName, "subnet_id", postConfigRegex),
			),
		},
	})
}

func TestAccPostgreSQLVirtualNetworkRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	r := PostgreSQLVirtualNetworkRuleResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckPostgreSQLVirtualNetworkRuleDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccPostgreSQLVirtualNetworkRule_multipleSubnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "rule1")
	r := PostgreSQLVirtualNetworkRuleResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleSubnets(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccPostgreSQLVirtualNetworkRule_IgnoreEndpointValid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_virtual_network_rule", "test")
	r := PostgreSQLVirtualNetworkRuleResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreEndpointValid(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ignore_missing_vnet_service_endpoint").HasValue("true"),
			),
		},
	})
}

func (t PostgreSQLVirtualNetworkRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.VirtualNetworkRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.VirtualNetworkRulesClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Postgresql Virtual Network Rule (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckPostgreSQLVirtualNetworkRuleDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Postgres.VirtualNetworkRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		ruleName := rs.Primary.Attributes["name"]

		future, err := client.Delete(ctx, resourceGroup, serverName, ruleName)
		if err != nil {
			// If the error is that the resource we want to delete does not exist in the first
			// place (404), then just return with no error.
			if response.WasNotFound(future.Response()) {
				return nil
			}

			return fmt.Errorf("Error deleting PostgreSQL Virtual Network Rule: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			// Same deal as before. Just in case.
			if response.WasNotFound(future.Response()) {
				return nil
			}

			return fmt.Errorf("Error deleting PostgreSQL Virtual Network Rule: %+v", err)
		}

		return nil
	}
}

func (PostgreSQLVirtualNetworkRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-SN-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctest-PSQL-VNET-rule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  subnet_id           = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PostgreSQLVirtualNetworkRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_virtual_network_rule" "import" {
  name                = azurerm_postgresql_virtual_network_rule.test.name
  resource_group_name = azurerm_postgresql_virtual_network_rule.test.resource_group_name
  server_name         = azurerm_postgresql_virtual_network_rule.test.server_name
  subnet_id           = azurerm_postgresql_virtual_network_rule.test.subnet_id
}
`, r.basic(data))
}

func (PostgreSQLVirtualNetworkRuleResource) subnetSwitchPre(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "acctest-SN1-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-SN2-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctest-PSQL-VNET-rule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  subnet_id           = azurerm_subnet.test1.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PostgreSQLVirtualNetworkRuleResource) subnetSwitchPost(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.7.29.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "acctest-SN1-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctest-SN2-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.128/25"
  service_endpoints    = ["Microsoft.Sql"]
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                = "acctest-PSQL-VNET-rule-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  subnet_id           = azurerm_subnet.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PostgreSQLVirtualNetworkRuleResource) multipleSubnets(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
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

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }
}

resource "azurerm_postgresql_virtual_network_rule" "rule1" {
  name                                 = "acctestpostgresqlvnetrule1%d"
  resource_group_name                  = azurerm_resource_group.test.name
  server_name                          = azurerm_postgresql_server.test.name
  subnet_id                            = azurerm_subnet.vnet1_subnet1.id
  ignore_missing_vnet_service_endpoint = false
}

resource "azurerm_postgresql_virtual_network_rule" "rule2" {
  name                                 = "acctestpostgresqlvnetrule2%d"
  resource_group_name                  = azurerm_resource_group.test.name
  server_name                          = azurerm_postgresql_server.test.name
  subnet_id                            = azurerm_subnet.vnet1_subnet2.id
  ignore_missing_vnet_service_endpoint = false
}

resource "azurerm_postgresql_virtual_network_rule" "rule3" {
  name                                 = "acctestpostgresqlvnetrule3%d"
  resource_group_name                  = azurerm_resource_group.test.name
  server_name                          = azurerm_postgresql_server.test.name
  subnet_id                            = azurerm_subnet.vnet2_subnet1.id
  ignore_missing_vnet_service_endpoint = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (PostgreSQLVirtualNetworkRuleResource) ignoreEndpointValid(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-psql-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.7.29.0/29"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-SN-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.7.29.0/29"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctestpostgresqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "GP_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_virtual_network_rule" "test" {
  name                                 = "acctestpostgresqlvnetrule%d"
  resource_group_name                  = azurerm_resource_group.test.name
  server_name                          = azurerm_postgresql_server.test.name
  subnet_id                            = azurerm_subnet.test.id
  ignore_missing_vnet_service_endpoint = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
