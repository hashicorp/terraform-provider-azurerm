// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MySQLServerResource struct{}

func TestAccMySQLServer_basicFiveSeven(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "5.7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
	})
}

func TestAccMySQLServer_basicEightZero(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "8.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
	})
}

func TestAccMySQLServer_autogrowOnly(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}
	mysqlVersion := "5.7"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autogrow(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
		{
			Config: r.basic(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
	})
}

func TestAccMySQLServer_requiresImport(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "5.7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMySQLServer_complete(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "8.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
	})
}

func TestAccMySQLServer_update(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}
	mysqlVersion := "8.0"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
		{
			Config: r.complete(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
		{
			Config: r.complete2(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "threat_detection_policy.0.storage_account_access_key"),
		{
			Config: r.complete3(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password", "threat_detection_policy.0.storage_account_access_key"),
		{
			Config: r.basic(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
	})
}

func TestAccMySQLServer_updateSKU(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
		{
			Config: r.sku(data, "MO_Gen5_16"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
	})
}

func TestAccMySQLServer_createReplica(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}
	mysqlVersion := "8.0"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			Config: r.createReplica(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMySQLServer_createPointInTimeRestore(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}
	restoreTime := time.Now().Add(11 * time.Minute)
	mysqlVersion := "8.0"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, mysqlVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
		{
			PreConfig: func() { time.Sleep(restoreTime.Sub(time.Now().Add(-7 * time.Minute))) },
			Config:    r.createPointInTimeRestore(data, mysqlVersion, restoreTime.Format(time.RFC3339)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"),
	})
}

func TestAccMySQLServer_infrastructureEncryption(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipping since `azurerm_mysql_server` is deprecated and will be removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_mysql_server", "test")
	r := MySQLServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.infrastructureEncryption(data, "5.7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_login_password"), // not returned as sensitive
	})
}

func (t MySQLServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servers.ParseServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.MySQL.MySqlClient.Servers.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading MySQL Server (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (MySQLServerResource) basic(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
  version                          = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func (MySQLServerResource) infrastructureEncryption(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                              = "acctestmysqlsvr-%d"
  location                          = azurerm_resource_group.test.location
  resource_group_name               = azurerm_resource_group.test.name
  sku_name                          = "GP_Gen5_2"
  administrator_login               = "acctestun"
  administrator_login_password      = "H@Sh1CoR3!"
  ssl_enforcement_enabled           = true
  ssl_minimal_tls_version_enforced  = "TLS1_1"
  storage_mb                        = 51200
  version                           = "%s"
  infrastructure_encryption_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func (MySQLServerResource) complete(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%[1]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  auto_grow_enabled                = true
  backup_retention_days            = 7
  create_mode                      = "Default"
  geo_redundant_backup_enabled     = false
  public_network_access_enabled    = true
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_2"
  storage_mb                       = 51200
  version                          = "%[3]s"
  threat_detection_policy {
    enabled              = true
    disabled_alerts      = ["Sql_Injection", "Data_Exfiltration"]
    email_account_admins = true
    email_addresses      = ["pearcec@example.com", "admin@example.com"]
    retention_days       = 7
  }
  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func (MySQLServerResource) complete2(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mysql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%[1]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!updated"
  auto_grow_enabled                = true
  backup_retention_days            = 7
  create_mode                      = "Default"
  geo_redundant_backup_enabled     = false
  ssl_enforcement_enabled          = false
  ssl_minimal_tls_version_enforced = "TLSEnforcementDisabled"
  storage_mb                       = 51200
  version                          = "%[3]s"
  threat_detection_policy {
    enabled                    = true
    disabled_alerts            = ["Sql_Injection"]
    email_account_admins       = true
    email_addresses            = ["pearcec@example.com"]
    retention_days             = 7
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func (MySQLServerResource) complete3(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mysql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_mysql_server" "test" {
  name                             = "acctestmysqlsvr-%[1]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!updated"
  auto_grow_enabled                = true
  backup_retention_days            = 7
  create_mode                      = "Default"
  geo_redundant_backup_enabled     = false
  ssl_enforcement_enabled          = false
  ssl_minimal_tls_version_enforced = "TLSEnforcementDisabled"
  storage_mb                       = 51200
  version                          = "%[3]s"
  threat_detection_policy {
    enabled                    = true
    email_account_admins       = true
    retention_days             = 7
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
  }
}
`, data.RandomInteger, data.Locations.Primary, version)
}

func (r MySQLServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "import" {
  name                = azurerm_mysql_server.test.name
  location            = azurerm_mysql_server.test.location
  resource_group_name = azurerm_mysql_server.test.resource_group_name
  sku_name            = "GP_Gen5_2"
  version             = "5.7"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}
`, r.basic(data, "5.7"))
}

func (MySQLServerResource) sku(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "%s"
  version             = "5.7"

  storage_mb                   = 16777216
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  ssl_enforcement_enabled      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sku)
}

func (MySQLServerResource) autogrow(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mysql_server" "test" {
  name                = "acctestmysqlsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "GP_Gen5_2"
  version             = "%s"

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  auto_grow_enabled            = true
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  ssl_enforcement_enabled      = true
  storage_mb                   = 51200
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, version)
}

func (r MySQLServerResource) createReplica(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "replica" {
  name                = "acctestmysqlsvr-%d-replica"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "GP_Gen5_2"
  version             = "%s"
  storage_mb          = 51200

  create_mode                      = "Replica"
  creation_source_server_id        = azurerm_mysql_server.test.id
  public_network_access_enabled    = false
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
}
`, r.basic(data, version), data.RandomInteger, version)
}

func (r MySQLServerResource) createPointInTimeRestore(data acceptance.TestData, version, restoreTime string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mysql_server" "restore" {
  name                = "acctestmysqlsvr-%d-restore"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "GP_Gen5_2"
  version             = "%s"

  create_mode                      = "PointInTimeRestore"
  creation_source_server_id        = azurerm_mysql_server.test.id
  restore_point_in_time            = "%s"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
}
`, r.basic(data, version), data.RandomInteger, version, restoreTime)
}
