package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSqlDatabaseLongTermPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database_long_term_retention_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseLongTermPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlDatabaseLongTermPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseLongTermPolicyExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSqlDatabaseLongTermPolicy_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_mssql_database_long_term_retention_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseLongTermPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSqlDatabaseLongTermPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseLongTermPolicyExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMSqlDatabaseLongTermPolicy_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_mssql_database_long_term_retention_policy"),
			},
		},
	})
}

func testCheckAzureRMSqlDatabaseLongTermPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Sql.BackupLongTermRetentionPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		databaseName := rs.Primary.Attributes["database_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, databaseName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Long Term Policy for SQL Database %q (server %q / resource group %q) was not found", databaseName, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMSqlDatabaseLongTermPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Sql.BackupLongTermRetentionPoliciesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_database_long_term_retention_policy" {
			continue
		}

		databaseName := rs.Primary.Attributes["database_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, databaseName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Long Term Policy for SQL Database %q (server %q / resource group %q) still exists: %+v", databaseName, serverName, resourceGroup, resp)
	}

	return nil
}

func testAccAzureRMSqlDatabaseLongTermPolicy_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test.name
  location                         = azurerm_resource_group.test.location
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_mssql_database_long_term_retention_policy" "test" {
  database_name       = azurerm_sql_database.test.name
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name

  weekly_retention  = "P1W"
  monthly_retention = "P1M"
  yearly_retention  = "P1Y"
  week_of_year      = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSqlDatabaseLongTermPolicy_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database_long_term_retention_policy" "import" {
  database_name       = azurerm_sql_database.test.name
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name

  weekly_retention  = "P1W"
  monthly_retention = "P1M"
  yearly_retention  = "P1Y"
  week_of_year      = 1
}
`, testAccAzureRMSqlDatabaseLongTermPolicy_basic(data))
}
