package databasemigration_test

import (
	`context`
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databasemigration/parse"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DatabaseMigrationServiceResource struct {
}

func TestAccDatabaseMigrationService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_service", "test")
	r := DatabaseMigrationServiceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard_1vCores"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabaseMigrationService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_service", "test")
	r := DatabaseMigrationServiceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subnet_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard_1vCores"),
				check.That(data.ResourceName).Key("tags.name").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabaseMigrationService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_service", "test")
	r := DatabaseMigrationServiceResource{}

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

func TestAccDatabaseMigrationService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_service", "test")
	r := DatabaseMigrationServiceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.name").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t DatabaseMigrationServiceResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DatabaseMigration.ServicesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Database Migration Service (Service Name %q / Group Name %q) does not exist", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.ServiceProperties != nil), nil
}

func (DatabaseMigrationServiceResource) base(data acceptance.TestData) string {
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
  name                 = "acctestSubnet-dbms-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DatabaseMigrationServiceResource) basic(data acceptance.TestData) string {
	template := DatabaseMigrationServiceResource{}.base(data)

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

func (DatabaseMigrationServiceResource) complete(data acceptance.TestData) string {
	template := DatabaseMigrationServiceResource{}.base(data)

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

func (DatabaseMigrationServiceResource) requiresImport(data acceptance.TestData) string {
	template := DatabaseMigrationServiceResource{}.basic(data)
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
