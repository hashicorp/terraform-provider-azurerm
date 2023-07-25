// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataProtectionBackupInstancePostgreSQLResource struct{}

func TestAccDataProtectionBackupInstancePostgreSQL_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql", "test")
	r := DataProtectionBackupInstancePostgreSQLResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstancePostgreSQL_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql", "test")
	r := DataProtectionBackupInstancePostgreSQLResource{}
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

func TestAccDataProtectionBackupInstancePostgreSQL_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql", "test")
	r := DataProtectionBackupInstancePostgreSQLResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstancePostgreSQL_keyVaultAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql", "test")
	r := DataProtectionBackupInstancePostgreSQLResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyVaultAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstancePostgreSQL_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql", "test")
	r := DataProtectionBackupInstancePostgreSQLResource{}
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

func (r DataProtectionBackupInstancePostgreSQLResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstances.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r DataProtectionBackupInstancePostgreSQLResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_server" "test" {
  name                = "acctest-postgresql-server-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "B_Gen5_2"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_firewall_rule" "test" {
  name                = "AllowAllWindowsAzureIps"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctest-postgresql-database-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  charset             = "UTF8"
  collation           = "English_United States.1252"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dataprotection-vault-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "test" {
  name                       = "acctest%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = ["Create", "Get"]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }

  access_policy {
    tenant_id = azurerm_data_protection_backup_vault.test.identity.0.tenant_id
    object_id = azurerm_data_protection_backup_vault.test.identity.0.principal_id

    key_permissions = ["Create", "Get"]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "acctestsecret%[1]d"
  value        = "Server=${azurerm_postgresql_server.test.name}.postgres.database.azure.com;Database=${azurerm_postgresql_database.test.name};Port=5432;User Id=psqladmin@${azurerm_postgresql_server.test.name};Password=H@Sh1CoR3!;Ssl Mode=Require;"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_postgresql_server.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity.0.principal_id
}

resource "azurerm_data_protection_backup_policy_postgresql" "test" {
  name                            = "acctest-dp-%[1]d"
  resource_group_name             = azurerm_resource_group.test.name
  vault_name                      = azurerm_data_protection_backup_vault.test.name
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P4M"
}

resource "azurerm_data_protection_backup_policy_postgresql" "another" {
  name                            = "acctest-dp-second-%[1]d"
  resource_group_name             = azurerm_resource_group.test.name
  vault_name                      = azurerm_data_protection_backup_vault.test.name
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P3M"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(16))
}

func (r DataProtectionBackupInstancePostgreSQLResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql" "test" {
  name             = "acctest-dbi-%d"
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_vault.test.id
  database_id      = azurerm_postgresql_database.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql.test.id
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstancePostgreSQLResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql" "import" {
  name             = azurerm_data_protection_backup_instance_postgresql.test.name
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_instance_postgresql.test.vault_id
  database_id      = azurerm_postgresql_database.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql.test.id
}
`, config)
}

func (r DataProtectionBackupInstancePostgreSQLResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql" "test" {
  name             = "acctest-dbi-%d"
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_vault.test.id
  database_id      = azurerm_postgresql_database.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql.another.id
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstancePostgreSQLResource) keyVaultAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql" "test" {
  name                                    = "acctest-dbi-%d"
  location                                = azurerm_resource_group.test.location
  vault_id                                = azurerm_data_protection_backup_vault.test.id
  database_id                             = azurerm_postgresql_database.test.id
  backup_policy_id                        = azurerm_data_protection_backup_policy_postgresql.another.id
  database_credential_key_vault_secret_id = azurerm_key_vault_secret.test.versionless_id
}
`, template, data.RandomInteger)
}
