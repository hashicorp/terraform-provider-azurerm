// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databasemigration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/projectresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DatabaseMigrationProjectResource struct{}

func TestAccDatabaseMigrationProject_basicSQLToSQLDB(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "SQL", "SQLDB"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_platform").HasValue("SQL"),
				check.That(data.ResourceName).Key("target_platform").HasValue("SQLDB"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabaseMigrationProject_basicPostgreSqlToAzureDbForPostgreSql(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "PostgreSql", "AzureDbForPostgreSql"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_platform").HasValue("PostgreSql"),
				check.That(data.ResourceName).Key("target_platform").HasValue("AzureDbForPostgreSql"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabaseMigrationProject_basicMySQLToAzureDbForMySql(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "MySQL", "AzureDbForMySql"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_platform").HasValue("MySQL"),
				check.That(data.ResourceName).Key("target_platform").HasValue("AzureDbForMySql"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabaseMigrationProject_basicMongoDbToMongoDb(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "MongoDb", "MongoDb"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("source_platform").HasValue("MongoDb"),
				check.That(data.ResourceName).Key("target_platform").HasValue("MongoDb"),
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
			Config: r.basic(data, "SQL", "SQLDB"),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDatabaseMigrationProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "SQL", "SQLDB"),
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
			Config: r.basic(data, "SQL", "SQLDB"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t DatabaseMigrationProjectResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := projectresource.ParseProjectID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DatabaseMigration.ProjectsClient.ProjectsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s", *id)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DatabaseMigrationProjectResource) basic(data acceptance.TestData, sourcePlatform string, targetPlatform string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_database_migration_project" "test" {
  name                = "acctestDbmsProject-%d"
  service_name        = azurerm_database_migration_service.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  source_platform     = "%s"
  target_platform     = "%s"
}
`, DatabaseMigrationServiceResource{}.basic(data), data.RandomInteger, sourcePlatform, targetPlatform)
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
	template := DatabaseMigrationProjectResource{}.basic(data, "SQL", "SQLDB")
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
