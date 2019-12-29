package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataMigrationService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_migration_service", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataMigrationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataMigrationServiceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard_1vCores"),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "Cloud"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataMigrationService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_migration_service", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataMigrationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataMigrationService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataMigrationServiceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard_1vCores"),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "Cloud"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.name", "test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataMigrationService_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_data_migration_service", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataMigrationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataMigrationServiceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMDataMigrationService_requiresImport),
		},
	})
}

func TestAccAzureRMDataMigrationService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_migration_service", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataMigrationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataMigrationServiceExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMDataMigrationService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataMigrationServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.name", "test"),
				),
			},
		},
	})
}

func testCheckAzureRMDataMigrationServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Data Migration Service not found: %s", resourceName)
		}

		groupName := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).DataMigration.ServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, groupName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Data Migration Service (Service Name %q / Group Name %q) does not exist", name, groupName)
			}
			return fmt.Errorf("Bad: Get on ServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDataMigrationServiceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataMigration.ServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_migration_service" {
			continue
		}

		groupName := rs.Primary.Attributes["resource_group_name"]
		name := rs.Primary.Attributes["name"]

		if resp, err := client.Get(ctx, groupName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on ServicesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDataMigrationService_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dms-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVnet-dms-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name 				   = "acctestSubnet-dms-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataMigrationService_basic(data acceptance.TestData) string {
	template := testAccAzureRMDataMigrationService_base(data)

	return fmt.Sprintf(`
%s

resource "azurerm_data_migration_service" "test" {
	name                = "acctestDms-%d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name
	subnet_id   		= azurerm_subnet.test.id
	sku_name            = "Standard_1vCores"
}
`, template, data.RandomInteger)
}

func testAccAzureRMDataMigrationService_complete(data acceptance.TestData) string {
	template := testAccAzureRMDataMigrationService_base(data)

	return fmt.Sprintf(`
%s

resource "azurerm_data_migration_service" "test" {
	name                = "acctestDms-%d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name
	subnet_id   		= azurerm_subnet.test.id
	sku_name            = "Standard_1vCores"
	tags 				= {
		name			= "test"
	}
}
`, template, data.RandomInteger)
}

func testAccAzureRMDataMigrationService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDataMigrationService_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_migration_service" "import" {
  name                = azurerm_data_migration_service.test.name
  location            = azurerm_data_migration_service.test.location
  resource_group_name = azurerm_data_migration_service.test.resource_group_name
  subnet_id   		  = azurerm_data_migration_service.test.subnet_id
  sku_name            = azurerm_data_migration_service.test.sku_name
}
`, template)
}
