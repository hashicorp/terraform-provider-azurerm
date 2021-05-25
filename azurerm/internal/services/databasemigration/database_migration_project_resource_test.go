package databasemigration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databasemigration/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DatabaseMigrationProjectResource struct {
}

func TestAccDatabaseMigrationProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_platform").HasValue("SQL"),
				check.That(data.ResourceName).Key("target_platform").HasValue("SQLDB"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabaseMigrationProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_platform").HasValue("SQL"),
				check.That(data.ResourceName).Key("target_platform").HasValue("SQLDB"),
				check.That(data.ResourceName).Key("tags.name").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabaseMigrationProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDatabaseMigrationProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.name").HasValue("Test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t DatabaseMigrationProjectResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ProjectID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DatabaseMigration.ProjectsClient.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Database Migration Project (Project Name %q / Service Name %q / Group Name %q) does not exist", id.Name, id.ServiceName, id.ResourceGroup)
	}

	return utils.Bool(resp.ProjectProperties != nil), nil
}

func (DatabaseMigrationProjectResource) basic(data acceptance.TestData) string {
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
`, DatabaseMigrationServiceResource{}.basic(data), data.RandomInteger)
}

func (DatabaseMigrationProjectResource) complete(data acceptance.TestData) string {
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
`, DatabaseMigrationServiceResource{}.basic(data), data.RandomInteger)
}

func (DatabaseMigrationProjectResource) requiresImport(data acceptance.TestData) string {
	template := DatabaseMigrationProjectResource{}.basic(data)
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
