package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureMySqlAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_active_directory_administrator", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureMySqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureMySqlAdministrator_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureMySqlAdministratorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "login", "sqladmin"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureMySqlAdministrator_withUpdates(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureMySqlAdministratorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "login", "sqladmin2"),
				),
			},
		},
	})
}

func TestAccAzureMySqlAdministrator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_active_directory_administrator", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureMySqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureMySqlAdministrator_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureMySqlAdministratorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "login", "sqladmin"),
				),
			},
			{
				Config:      testAccAzureMySqlAdministrator_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_mysql_active_directory_administrator"),
			},
		},
	})
}

func TestAccAzureMySqlAdministrator_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mysql_active_directory_administrator", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureMySqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureMySqlAdministrator_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureMySqlAdministratorExists(data.ResourceName),
					testCheckAzureMySqlAdministratorDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureMySqlAdministratorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ServerAdministratorsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		_, err := client.Get(ctx, resourceGroup, serverName)
		return err
	}
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

func testCheckAzureMySqlAdministratorDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MySQL.ServerAdministratorsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mysql_active_directory_administrator" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("MySQL AD Administrator (server %q / resource group %q) still exists: %+v", serverName, resourceGroup, resp)
	}

	return nil
}

func testAccAzureMySqlAdministrator_basic(data acceptance.TestData) string {
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

func testAccAzureMySqlAdministrator_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_active_directory_administrator" "import" {
  server_name         = azurerm_mysql_active_directory_administrator.test.server_name
  resource_group_name = azurerm_mysql_active_directory_administrator.test.resource_group_name
  login               = azurerm_mysql_active_directory_administrator.test.login
  tenant_id           = azurerm_mysql_active_directory_administrator.test.tenant_id
  object_id           = azurerm_mysql_active_directory_administrator.test.object_id
}
`, testAccAzureMySqlAdministrator_basic(data))
}

func testAccAzureMySqlAdministrator_withUpdates(data acceptance.TestData) string {
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
