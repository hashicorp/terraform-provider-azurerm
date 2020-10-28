package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databasemigration/parse"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDatabaseMigrationService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabaseMigrationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabaseMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationServiceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard_1vCores"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDatabaseMigrationService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabaseMigrationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabaseMigrationService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationServiceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard_1vCores"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.name", "test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDatabaseMigrationService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabaseMigrationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabaseMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationServiceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDatabaseMigrationService_requiresImport),
		},
	})
}

func TestAccAzureRMDatabaseMigrationService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabaseMigrationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabaseMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDatabaseMigrationService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.name", "test"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDatabaseMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabaseMigrationServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDatabaseMigrationServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Database Migration Service not found: %s", resourceName)
		}

		id, err := parse.DatabaseMigrationServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).DatabaseMigration.ServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Database Migration Service (Service Name %q / Group Name %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on ServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDatabaseMigrationServiceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DatabaseMigration.ServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_database_migration_service" {
			continue
		}

		id, err := parse.DatabaseMigrationServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on ServicesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDatabaseMigrationService_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dbms-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVnet-dbms-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name               = "acctestSubnet-dbms-%d"
  virtual_network_id = azurerm_virtual_network.test.id
  address_prefix     = "10.0.1.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDatabaseMigrationService_basic(data acceptance.TestData) string {
	template := testAccAzureRMDatabaseMigrationService_base(data)

	return fmt.Sprintf(`
%s

resource "azurerm_database_migration_service" "test" {
  name                = "acctestDbms-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_1vCores"
}
`, template, data.RandomInteger)
}

func testAccAzureRMDatabaseMigrationService_complete(data acceptance.TestData) string {
	template := testAccAzureRMDatabaseMigrationService_base(data)

	return fmt.Sprintf(`
%s

resource "azurerm_database_migration_service" "test" {
  name                = "acctestDbms-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_1vCores"
  tags = {
    name = "test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMDatabaseMigrationService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDatabaseMigrationService_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_database_migration_service" "import" {
  name                = azurerm_database_migration_service.test.name
  location            = azurerm_database_migration_service.test.location
  resource_group_name = azurerm_database_migration_service.test.resource_group_name
  subnet_id           = azurerm_database_migration_service.test.subnet_id
  sku_name            = azurerm_database_migration_service.test.sku_name
}
`, template)
}
