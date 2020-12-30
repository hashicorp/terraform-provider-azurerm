package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MySqlAdministratorResource struct {
}

func TestAccMySqlAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_active_directory_administrator", "test")
	r := MySqlAdministratorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("login").HasValue("sqladmin"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUpdates(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("login").HasValue("sqladmin2"),
			),
		},
	})
}

func TestAccMySqlAdministrator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_active_directory_administrator", "test")
	r := MySqlAdministratorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("login").HasValue("sqladmin"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_mysql_active_directory_administrator"),
		},
	})
}

func TestAccMySqlAdministrator_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_active_directory_administrator", "test")
	r := MySqlAdministratorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckAzureMySqlAdministratorDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func (t MySqlAdministratorResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]

	resp, err := clients.MySQL.ServerAdministratorsClient.Get(ctx, resourceGroup, serverName)
	if err != nil {
		return nil, fmt.Errorf("reading MySQL Administrator (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckAzureMySqlAdministratorDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ServerAdministratorsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		if _, err := client.Delete(ctx, resourceGroup, serverName); err != nil {
			return fmt.Errorf("Bad: Delete on mysqlAdministratorClient: %+v", err)
		}

		return nil
	}
}

func (MySqlAdministratorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  version                      = "5.7"
  ssl_enforcement_enabled      = true
}

resource "azurerm_mysql_active_directory_administrator" "test" {
  server_name         = azurerm_mysql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = "sqladmin"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r MySqlAdministratorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_active_directory_administrator" "import" {
  server_name         = azurerm_mysql_active_directory_administrator.test.server_name
  resource_group_name = azurerm_mysql_active_directory_administrator.test.resource_group_name
  login               = azurerm_mysql_active_directory_administrator.test.login
  tenant_id           = azurerm_mysql_active_directory_administrator.test.tenant_id
  object_id           = azurerm_mysql_active_directory_administrator.test.object_id
}
`, r.basic(data))
}

func (MySqlAdministratorResource) withUpdates(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
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
  version                      = "5.7"
  ssl_enforcement_enabled      = true
}

resource "azurerm_mysql_active_directory_administrator" "test" {
  server_name         = azurerm_mysql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = "sqladmin2"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
