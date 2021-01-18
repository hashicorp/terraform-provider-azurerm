package mssql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MsSqlDatabaseResource struct{}

func TestAccMsSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

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

func TestAccMsSqlDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("max_size_gb").HasValue("1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
		data.ImportStep("sample_name"),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("license_type").HasValue("LicenseIncluded"),
				check.That(data.ResourceName).Key("max_size_gb").HasValue("2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Staging"),
			),
		},
		data.ImportStep("sample_name"),
	})
}

func TestAccMsSqlDatabase_elasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.elasticPool(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("elastic_pool_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("ElasticPool"),
			),
		},
		data.ImportStep(),
		{
			Config: r.elasticPoolDisassociation(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_GP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.gp(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_GP_Serverless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.gpServerless(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_pause_delay_in_minutes").HasValue("70"),
				check.That(data.ResourceName).Key("min_capacity").HasValue("0.75"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_S_Gen5_2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.gpServerlessUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_pause_delay_in_minutes").HasValue("90"),
				check.That(data.ResourceName).Key("min_capacity").HasValue("1.25"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_S_Gen5_2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_BC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.bc(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_scale").HasValue("true"),
				check.That(data.ResourceName).Key("sku_name").HasValue("BC_Gen5_2"),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.bcUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_scale").HasValue("false"),
				check.That(data.ResourceName).Key("sku_name").HasValue("BC_Gen5_2"),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_HS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.hs(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_replica_count").HasValue("2"),
				check.That(data.ResourceName).Key("sku_name").HasValue("HS_Gen5_2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.hsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_replica_count").HasValue("4"),
				check.That(data.ResourceName).Key("sku_name").HasValue("HS_Gen5_2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_createCopyMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "copy")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.createCopyMode(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep("create_mode", "creation_source_database_id"),
	})
}

func TestAccMsSqlDatabase_createPITRMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),

		{
			PreConfig: func() { time.Sleep(7 * time.Minute) },
			Config:    r.createPITRMode(data),
			Check: resource.ComposeTestCheckFunc(
				check.That("azurerm_mssql_database.pitr").ExistsInAzure(r),
			),
		},

		data.ImportStep("create_mode", "creation_source_database_id", "restore_point_in_time"),
	})
}

func TestAccMsSqlDatabase_createSecondaryMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "secondary")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.createSecondaryMode(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep("create_mode", "creation_source_database_id", "sample_name"),
	})
}

func TestAccMsSqlDatabase_createRestoreMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.createRestoreMode(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("create_mode", "creation_source_database_id"),

		{
			PreConfig: func() { time.Sleep(8 * time.Minute) },
			Config:    r.createRestoreModeDBDeleted(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		data.ImportStep(),

		{
			PreConfig: func() { time.Sleep(8 * time.Minute) },
			Config:    r.createRestoreModeDBRestored(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_mssql_database.restore").ExistsInAzure(r),
			),
		},

		data.ImportStep("create_mode", "restore_dropped_database_id"),
	})
}

func TestAccMsSqlDatabase_threatDetectionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.threatDetectionPolicy(data, "Enabled"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.state").HasValue("Enabled"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.retention_days").HasValue("15"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.disabled_alerts.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.email_account_admins").HasValue("Enabled"),
			),
		},
		data.ImportStep("sample_name", "threat_detection_policy.0.storage_account_access_key"),
		{
			Config: r.threatDetectionPolicy(data, "Disabled"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.state").HasValue("Disabled"),
			),
		},
		data.ImportStep("sample_name", "threat_detection_policy.0.storage_account_access_key"),
	})
}

func TestAccMsSqlDatabase_withBlobAuditingPolices(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBlobAuditingPolices(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("extended_auditing_policy.0.storage_account_access_key"),
		{
			Config: r.withBlobAuditingPolicesUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("extended_auditing_policy.0.storage_account_access_key"),
		{
			Config: r.withBlobAuditingPolicesDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateSku(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateSku2(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_minCapacity0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.minCapacity0(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_withLongTermRetentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withLongTermRetentionPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withLongTermRetentionPolicyUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_withShortTermRetentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withShortTermRetentionPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withShortTermRetentionPolicyUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (MsSqlDatabaseResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.DatabasesClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("SQL Database %q (Server %q, Resource Group %q) does not exist", id.Name, id.ServerName, id.ResourceGroup)
		}

		return nil, fmt.Errorf("reading SQL Database %q (Server %q, Resource Group %q): %v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (MsSqlDatabaseResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_sql_server.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "import" {
  name      = azurerm_mssql_database.test.name
  server_id = azurerm_sql_server.test.id
}
`, r.basic(data))
}

func (r MsSqlDatabaseResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_sql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  max_size_gb  = 1
  sample_name  = "AdventureWorksLT"
  sku_name     = "GP_Gen5_2"

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_sql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "LicenseIncluded"
  max_size_gb  = 2
  sku_name     = "GP_Gen5_2"

  tags = {
    ENV = "Staging"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) elasticPool(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name            = "acctest-db-%[2]d"
  server_id       = azurerm_sql_server.test.id
  elastic_pool_id = azurerm_mssql_elasticpool.test.id
  sku_name        = "ElasticPool"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) elasticPoolDisassociation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_sql_server.test.id
  sku_name  = "GP_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) gp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_sql_server.test.id
  sku_name  = "GP_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) gpServerless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name                        = "acctest-db-%[2]d"
  server_id                   = azurerm_sql_server.test.id
  auto_pause_delay_in_minutes = 70
  min_capacity                = 0.75
  sku_name                    = "GP_S_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) gpServerlessUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name                        = "acctest-db-%[2]d"
  server_id                   = azurerm_sql_server.test.id
  auto_pause_delay_in_minutes = 90
  min_capacity                = 1.25
  sku_name                    = "GP_S_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) hs(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[2]d"
  server_id          = azurerm_sql_server.test.id
  read_replica_count = 2
  sku_name           = "HS_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) hsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[2]d"
  server_id          = azurerm_sql_server.test.id
  read_replica_count = 4
  sku_name           = "HS_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) bc(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name           = "acctest-db-%[2]d"
  server_id      = azurerm_sql_server.test.id
  read_scale     = true
  sku_name       = "BC_Gen5_2"
  zone_redundant = true
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) bcUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name           = "acctest-db-%[2]d"
  server_id      = azurerm_sql_server.test.id
  read_scale     = false
  sku_name       = "BC_Gen5_2"
  zone_redundant = false
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) createCopyMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "copy" {
  name                        = "acctest-dbc-%[2]d"
  server_id                   = azurerm_sql_server.test.id
  create_mode                 = "Copy"
  creation_source_database_id = azurerm_mssql_database.test.id
}
`, r.complete(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) createPITRMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "pitr" {
  name                        = "acctest-dbp-%[2]d"
  server_id                   = azurerm_sql_server.test.id
  create_mode                 = "PointInTimeRestore"
  restore_point_in_time       = "%[3]s"
  creation_source_database_id = azurerm_mssql_database.test.id

}
`, r.basic(data), data.RandomInteger, time.Now().Add(time.Duration(7)*time.Minute).UTC().Format(time.RFC3339))
}

func (r MsSqlDatabaseResource) createSecondaryMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "second" {
  name     = "acctestRG-mssql2-%[2]d"
  location = "%[3]s"
}

resource "azurerm_sql_server" "second" {
  name                         = "acctest-sqlserver2-%[2]d"
  resource_group_name          = azurerm_resource_group.second.name
  location                     = azurerm_resource_group.second.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-dbs-%[2]d"
  server_id                   = azurerm_sql_server.second.id
  create_mode                 = "Secondary"
  creation_source_database_id = azurerm_mssql_database.test.id

}
`, r.complete(data), data.RandomInteger, data.Locations.Secondary)
}

func (MsSqlDatabaseResource) createRestoreMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_mssql_database" "copy" {
  name                        = "acctest-dbc-%[1]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Copy"
  creation_source_database_id = azurerm_mssql_database.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlDatabaseResource) createRestoreModeDBDeleted(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlDatabaseResource) createRestoreModeDBRestored(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_mssql_database" "restore" {
  name                        = "acctest-dbr-%[1]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Restore"
  restore_dropped_database_id = azurerm_mssql_server.test.restorable_dropped_database_ids[0]
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlDatabaseResource) threatDetectionPolicy(data acceptance.TestData, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "test%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_sql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  max_size_gb  = 1
  sample_name  = "AdventureWorksLT"
  sku_name     = "GP_Gen5_2"

  threat_detection_policy {
    retention_days             = 15
    state                      = "%[3]s"
    disabled_alerts            = ["Sql_Injection"]
    email_account_admins       = "Enabled"
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
    use_server_default         = "Disabled"
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger, state)
}

func (r MsSqlDatabaseResource) withBlobAuditingPolices(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.test.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withBlobAuditingPolicesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.test2.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.test2.primary_access_key
    storage_account_access_key_is_secondary = false
    retention_in_days                       = 3
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withBlobAuditingPolicesDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name                     = "acctest-db-%[3]d"
  server_id                = azurerm_sql_server.test.id
  extended_auditing_policy = []
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) updateSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_sql_server.test.id
  sku_name  = "HS_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) updateSku2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_sql_server.test.id
  sku_name  = "HS_Gen5_4"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) minCapacity0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_sql_server.test.id

  min_capacity = 0
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withLongTermRetentionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 1
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withLongTermRetentionPolicyUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  long_term_retention_policy {
    weekly_retention = "P1W"
    yearly_retention = "P1Y"
    week_of_year     = 2
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withShortTermRetentionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  short_term_retention_policy {
    retention_days = 8
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withShortTermRetentionPolicyUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_sql_server.test.id
  short_term_retention_policy {
    retention_days = 10
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}
