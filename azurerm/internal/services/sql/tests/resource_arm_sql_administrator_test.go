package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlAdministrator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_active_directory_administrator", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlAdministrator_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "login", "sqladmin"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSqlAdministrator_withUpdates(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "login", "sqladmin2"),
				),
			},
		},
	})
}

func TestAccAzureRMSqlAdministrator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_active_directory_administrator", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlAdministrator_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "login", "sqladmin"),
				),
			},
			{
				Config:      testAccAzureRMSqlAdministrator_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_sql_active_directory_administrator"),
			},
		},
	})
}

func TestAccAzureRMSqlAdministrator_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_active_directory_administrator", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlAdministratorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlAdministrator_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlAdministratorExists(data.ResourceName),
					testCheckAzureRMSqlAdministratorDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSqlAdministratorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Sql.ServerAzureADAdministratorsClient
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

func testCheckAzureRMSqlAdministratorDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Sql.ServerAzureADAdministratorsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		if _, err := client.Delete(ctx, resourceGroup, serverName); err != nil {
			return fmt.Errorf("Bad: Delete on sqlAdministratorClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSqlAdministratorDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Sql.ServerAzureADAdministratorsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_sql_active_directory_administrator" {
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

		return fmt.Errorf("SQL Administrator (server %q / resource group %q) still exists: %+v", serverName, resourceGroup, resp)
	}

	return nil
}

func testAccAzureRMSqlAdministrator_basic(data acceptance.TestData) string {
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

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name         = azurerm_sql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = "sqladmin"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMSqlAdministrator_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_active_directory_administrator" "import" {
  server_name         = azurerm_sql_active_directory_administrator.test.server_name
  resource_group_name = azurerm_sql_active_directory_administrator.test.resource_group_name
  login               = azurerm_sql_active_directory_administrator.test.login
  tenant_id           = azurerm_sql_active_directory_administrator.test.tenant_id
  object_id           = azurerm_sql_active_directory_administrator.test.object_id
}
`, testAccAzureRMSqlAdministrator_basic(data))
}

func testAccAzureRMSqlAdministrator_withUpdates(data acceptance.TestData) string {
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

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_active_directory_administrator" "test" {
  server_name         = azurerm_sql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
  login               = "sqladmin2"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  object_id           = data.azurerm_client_config.current.client_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
