package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databasemigration/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDatabaseMigrationProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabaseMigrationProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabaseMigrationProject_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationProjectExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "source_platform", "SQL"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_platform", "SQLDB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDatabaseMigrationProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabaseMigrationProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabaseMigrationProject_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationProjectExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "source_platform", "SQL"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_platform", "SQLDB"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.name", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDatabaseMigrationProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabaseMigrationProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabaseMigrationProject_basic(data),
			},
			data.RequiresImportErrorStep(testAccAzureRMDatabaseMigrationProject_requiresImport),
		},
	})
}

func TestAccAzureRMDatabaseMigrationProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabaseMigrationProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabaseMigrationProject_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationProjectExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDatabaseMigrationProject_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationProjectExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.name", "Test"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDatabaseMigrationProject_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationProjectExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDatabaseMigrationProjectExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Database Migration Project not found: %s", resourceName)
		}

		id, err := parse.ProjectID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).DatabaseMigration.ProjectsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Database Migration Project (Project Name %q / Service Name %q / Group Name %q) does not exist", id.Name, id.ServiceName, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on ProjectsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDatabaseMigrationProjectDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DatabaseMigration.ProjectsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_database_migration_project" {
			continue
		}

		id, err := parse.ProjectID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on ProjectsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDatabaseMigrationProject_basic(data acceptance.TestData) string {
	template := testAccAzureRMDatabaseMigrationService_basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_database_migration_project" "test" {
  name                = "acctestDbmsProject-%d"
  service_name        = azurerm_database_migration_service.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  source_platform     = "SQL"
  target_platform     = "SQLDB"
}
`, template, data.RandomInteger)
}

func testAccAzureRMDatabaseMigrationProject_complete(data acceptance.TestData) string {
	template := testAccAzureRMDatabaseMigrationService_basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_database_migration_project" "test" {
  name                = "acctestDbmsProject-%d"
  service_name        = azurerm_database_migration_service.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  source_platform     = "SQL"
  target_platform     = "SQLDB"
  tags = {
    name = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDatabaseMigrationProject_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDatabaseMigrationProject_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_database_migration_project" "import" {
  name                = azurerm_database_migration_project.test.name
  service_name        = azurerm_database_migration_project.test.service_name
  resource_group_name = azurerm_database_migration_project.test.resource_group_name
  location            = azurerm_database_migration_project.test.location
  source_platform     = azurerm_database_migration_project.test.source_platform
  target_platform     = azurerm_database_migration_project.test.target_platform
}
`, template)
}
