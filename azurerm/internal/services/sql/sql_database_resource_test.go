package sql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SqlDatabaseResource struct{}

func TestAccSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("create_mode"),
	})
}

func TestAccSqlDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSqlDatabase_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccSqlDatabase_elasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.elasticPool(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("elastic_pool_name").HasValue(fmt.Sprintf("acctestep%d", data.RandomInteger)),
			),
		},
		data.ImportStep("create_mode"),
	})
}

func TestAccSqlDatabase_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccSqlDatabase_dataWarehouse(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataWarehouse(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_sql_database.test").ExistsInAzure(r),
			),
		},
		data.ImportStep("create_mode"),
	})
}

func TestAccSqlDatabase_restorePointInTime(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	restorePointInTime := time.Now().Add(15 * time.Minute)
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:                    r.basic(data),
			PreventPostDestroyRefresh: true,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			// we need to make sure we wait long enough to prevent this:
			// The source database, 'acctestdb201215000744760737', has not existed long enough to support restores.
			PreConfig: func() { time.Sleep(restorePointInTime.Sub(time.Now().Add(45 * time.Minute))) },
			Config:    r.restorePointInTime(data, restorePointInTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_sql_database.test_restore").ExistsInAzure(r),
			),
		},
	})
}

func TestAccSqlDatabase_collation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_Latin1_General_CP1_CI_AS"),
			),
		},
		{
			Config: r.collationUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("Japanese_Bushu_Kakusu_100_CS_AS_KS_WS"),
			),
		},
	})
}

func TestAccSqlDatabase_requestedServiceObjectiveName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.requestedServiceObjectiveName(data, "S0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("requested_service_objective_name").HasValue("S0"),
			),
		},
		{
			Config: r.requestedServiceObjectiveName(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("requested_service_objective_name").HasValue("S1"),
			),
		},
	})
}

func TestAccSqlDatabase_threatDetectionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.threatDetectionPolicy(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.state").HasValue("Enabled"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.retention_days").HasValue("15"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.disabled_alerts.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.email_account_admins").HasValue("Enabled"),
			),
		},
		data.ImportStep("create_mode", "threat_detection_policy.0.storage_account_access_key"),
		{
			Config: r.threatDetectionPolicy(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.state").HasValue("Disabled"),
			),
		},
	})
}

func TestAccSqlDatabase_readScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.readScale(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_scale").HasValue("true"),
			),
		},
		{
			Config: r.readScale(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_scale").HasValue("false"),
			),
		},
	})
}

func TestAccSqlDatabase_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundant(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("true"),
			),
		},
		{
			Config: r.zoneRedundant(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("false"),
			),
		},
	})
}

func TestAccSqlDatabase_bacpac(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bacpac(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_sql_database.test").ExistsInAzure(r),
			),
		},
	})
}

func TestAccSqlDatabase_withBlobAuditingPolices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "create_mode"),
		{
			Config: r.withBlobAuditingPolices(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key", "create_mode"),
		{
			Config: r.withBlobAuditingPolicesDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "create_mode"),
	})
}

func TestAccSqlDatabase_withBlobAuditingPolicesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBlobAuditingPolices(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extended_auditing_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("extended_auditing_policy.0.storage_account_access_key_is_secondary").HasValue("true"),
				check.That(data.ResourceName).Key("extended_auditing_policy.0.retention_in_days").HasValue("6"),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key", "create_mode"),
		{
			Config: r.withBlobAuditingPolicesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("extended_auditing_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("extended_auditing_policy.0.storage_account_access_key_is_secondary").HasValue("false"),
				check.That(data.ResourceName).Key("extended_auditing_policy.0.retention_in_days").HasValue("11"),
			),
		},
		data.ImportStep("administrator_login_password", "extended_auditing_policy.0.storage_account_access_key", "create_mode"),
	})
}

func TestAccSqlDatabase_onlineSecondaryMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sql_database", "test")
	r := SqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.onlineSecondaryMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "create_mode"),
	})
}

func (r SqlDatabaseResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DatabaseID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Sql.DatabasesClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving SQL Database %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlDatabaseResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DatabaseID(state.ID)
	if err != nil {
		return nil, err
	}
	if _, err := client.Sql.DatabasesClient.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name); err != nil {
		return nil, fmt.Errorf("deleting SQL Database %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r SqlDatabaseResource) basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r SqlDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sql_database" "import" {
  name                             = azurerm_sql_database.test.name
  resource_group_name              = azurerm_sql_database.test.resource_group_name
  server_name                      = azurerm_sql_database.test.server_name
  location                         = azurerm_sql_database.test.location
  edition                          = azurerm_sql_database.test.edition
  collation                        = azurerm_sql_database.test.collation
  max_size_bytes                   = azurerm_sql_database.test.max_size_bytes
  requested_service_objective_name = azurerm_sql_database.test.requested_service_objective_name
}
`, r.basic(data))
}

func (r SqlDatabaseResource) withTags(data acceptance.TestData) string {
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

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r SqlDatabaseResource) withTagsUpdate(data acceptance.TestData) string {
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

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r SqlDatabaseResource) dataWarehouse(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest_rg_%d"
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
  edition                          = "DataWarehouse"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  requested_service_objective_name = "DW400c"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r SqlDatabaseResource) restorePointInTime(data acceptance.TestData, restorePointInTime time.Time) string {
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

resource "azurerm_sql_database" "test_restore" {
  name                  = "acctestdb_restore%d"
  resource_group_name   = azurerm_resource_group.test.name
  server_name           = azurerm_sql_server.test.name
  location              = azurerm_resource_group.test.location
  create_mode           = "PointInTimeRestore"
  source_database_id    = azurerm_sql_database.test.id
  restore_point_in_time = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, restorePointInTime.UTC().Format(time.RFC3339))
}

func (r SqlDatabaseResource) elasticPool(data acceptance.TestData) string {
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

resource "azurerm_sql_elasticpool" "test" {
  name                = "acctestep%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  edition             = "Basic"
  dtu                 = 50
  pool_size           = 5000
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test.name
  location                         = azurerm_resource_group.test.location
  edition                          = azurerm_sql_elasticpool.test.edition
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  elastic_pool_name                = azurerm_sql_elasticpool.test.name
  requested_service_objective_name = "ElasticPool"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SqlDatabaseResource) collationUpdate(data acceptance.TestData) string {
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
  collation                        = "Japanese_Bushu_Kakusu_100_CS_AS_KS_WS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r SqlDatabaseResource) bacpac(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG_%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "bacpac"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "test.bacpac"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source                 = "testdata/sql_import.bacpac"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_firewall_rule" "test" {
  name                = "allowazure"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
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

  import {
    storage_uri                  = azurerm_storage_blob.test.url
    storage_key                  = azurerm_storage_account.test.primary_access_key
    storage_key_type             = "StorageAccessKey"
    administrator_login          = azurerm_sql_server.test.administrator_login
    administrator_login_password = azurerm_sql_server.test.administrator_login_password
    authentication_type          = "SQL"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SqlDatabaseResource) requestedServiceObjectiveName(data acceptance.TestData, requestedServiceObjectiveName string) string {
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
  requested_service_objective_name = %q
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, requestedServiceObjectiveName)
}

func (r SqlDatabaseResource) threatDetectionPolicy(data acceptance.TestData, state string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "test%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
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
  name                = "acctestdb%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name
  location            = azurerm_resource_group.test.location
  edition             = "Standard"
  collation           = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes      = "1073741824"

  threat_detection_policy {
    retention_days             = 15
    state                      = "%s"
    disabled_alerts            = ["Sql_Injection"]
    email_account_admins       = "Enabled"
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, state)
}

func (r SqlDatabaseResource) readScale(data acceptance.TestData, readScale bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "readscaletestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "readscaletestsqlserver%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                = "readscaletestdb%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name
  location            = azurerm_resource_group.test.location
  edition             = "Premium"
  collation           = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes      = "1073741824"
  read_scale          = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, readScale)
}

func (r SqlDatabaseResource) zoneRedundant(data acceptance.TestData, zoneRedundant bool) string {
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
  name                = "acctestdb%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name
  location            = azurerm_resource_group.test.location
  edition             = "Premium"
  collation           = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes      = "1073741824"
  zone_redundant      = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, zoneRedundant)
}

func (r SqlDatabaseResource) withBlobAuditingPolices(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test.name
  location                         = azurerm_resource_group.test.location
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"

  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.test.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(15))
}

func (r SqlDatabaseResource) withBlobAuditingPolicesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test.name
  location                         = azurerm_resource_group.test.location
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"

  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.test2.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.test2.primary_access_key
    storage_account_access_key_is_secondary = false
    retention_in_days                       = 11
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(15))
}

func (r SqlDatabaseResource) withBlobAuditingPolicesDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test.name
  location                         = azurerm_resource_group.test.location
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"

  extended_auditing_policy = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(15))
}

func (r SqlDatabaseResource) onlineSecondaryMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sql-%[1]d-p"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  server_name                      = azurerm_sql_server.test.name
  location                         = azurerm_resource_group.test.location
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-sql-%[1]d-s"
  location = "%[3]s"
}

resource "azurerm_sql_server" "test2" {
  name                         = "acctestsqlserver2%[1]d"
  resource_group_name          = azurerm_resource_group.test2.name
  location                     = azurerm_resource_group.test2.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_database" "test2" {
  name                = "acctestdb2%[1]d"
  resource_group_name = azurerm_resource_group.test2.name
  server_name         = azurerm_sql_server.test2.name
  location            = azurerm_resource_group.test2.location
  create_mode         = "OnlineSecondary"
  source_database_id  = azurerm_sql_database.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}
