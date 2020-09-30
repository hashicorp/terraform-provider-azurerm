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

func TestAccAzureRMMsSqlSyncGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_sync_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlSyncGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:    testAccAzureRMMsSqlSyncGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.sync"),
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.hub"),
					testCheckAzureRMMsSqlSyncGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlSyncGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_sync_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlSyncGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:    testAccAzureRMMsSqlSyncGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.sync"),
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.hub"),
					testCheckAzureRMMsSqlSyncGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMsSqlSyncGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_sync_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlSyncGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlSyncGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlSyncGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMsSqlSyncGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError(data.ResourceType),
			},
		},
	})
}

func TestAccAzureRMMsSqlSyncGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlSyncGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlSyncGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlSyncGroupExists(data.ResourceName),
					testCheckAzureRMMsSqlSyncGroupDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMMsSqlSyncGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.SyncGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["database_name"]
		syncGroupName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, databaseName, syncGroupName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("SQL Sync Group %q (database %q / server %q / resource group %q) was not found", syncGroupName, databaseName, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMMsSqlSyncGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.SyncGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_sync_group" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["database_name"]
		syncGroupName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, resourceGroup, serverName, databaseName, syncGroupName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("SQL Sync Group %q (database %q / server %q / resource group %q) still exists: %+v", syncGroupName, databaseName, serverName, resourceGroup, resp)
	}

	return nil
}

func testCheckAzureRMMsSqlSyncGroupDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MSSQL.SyncGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["database_name"]
		syncGroupName := rs.Primary.Attributes["name"]

		future, err := client.Delete(ctx, resourceGroup, serverName, databaseName, syncGroupName)

		if err != nil {
			return err
		}

		return future.WaitForCompletionRef(ctx, client.Client)
	}
}

func testAccAzureRMMsSqlSyncGroup_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "sync" {
  name      = "syncStore%[2]d"
  server_id = azurerm_mssql_server.test.id
  sku_name  = "S1"
}

resource "azurerm_mssql_database" "hub" {
  name        = "syncHub%[2]d"
  server_id   = azurerm_mssql_server.test.id
  sku_name    = "S2"
  sample_name = "AdventureWorksLT"
}
`, testAccAzureRMMsSqlServer_basic(data), data.RandomInteger)
}

func testAccAzureRMMsSqlSyncGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_sync_group" "test" {
  name        = "acctest-syncgroup-%[2]d"
  database_id = azurerm_mssql_database.hub.id

  conflict_resolution_policy = "HubWin"
  sync_database_id           = azurerm_mssql_database.sync.id

  hub_database_username = azurerm_mssql_server.test.administrator_login
  hub_database_password = "thisIsKat11"
}

`, testAccAzureRMMsSqlSyncGroup_template(data), data.RandomInteger)
}

func testAccAzureRMMsSqlSyncGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_sync_group" "test" {
  name        = "acctest-syncgroup-%[2]d"
  database_id = azurerm_mssql_database.hub.id

  conflict_resolution_policy = "HubWin"
  sync_database_id           = azurerm_mssql_database.sync.id

  hub_database_username = azurerm_mssql_server.test.administrator_login
  hub_database_password = "thisIsKat11"

  //primary_sync_member_name = azurerm_mssql_database.hub.name

  table {
    name = "[SalesLT].[Product]"

    column {
      name      = "[ProductID]"
      data_size = "4"
      data_type = "int"
    }

    column {
      name      = "[ProductNumber]"
      data_size = "25"
      data_type = "nvarchar"
    }

    column {
      name      = "[Color]"
      data_size = "15"
      data_type = "nvarchar"
    }
  }
}

`, testAccAzureRMMsSqlSyncGroup_template(data), data.RandomInteger)
}

func testAccAzureRMMsSqlSyncGroup_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_sync_group" "import" {
  name                       = azurerm_mssql_sync_group.test.name
  resource_group_name        = azurerm_mssql_sync_group.test.resource_group_name
  server_name                = azurerm_mssql_sync_group.test.server_name
  database_name              = azurerm_mssql_sync_group.test.database_name
  conflict_resolution_policy = azurerm_mssql_sync_group.test.conflict_resolution_policy
  interval                   = azurerm_mssql_sync_group.test.interval
  sync_database_id           = azurerm_mssql_sync_group.test.sync_database_id
  hub_database_username      = azurerm_mssql_sync_group.test.hub_database_username
  hub_database_password      = azurerm_mssql_sync_group.test.hub_database_password
  primary_sync_member_name   = azurerm_mssql_sync_group.test.primary_sync_member_name
}
`, testAccAzureRMMsSqlSyncGroup_basic(data))
}
