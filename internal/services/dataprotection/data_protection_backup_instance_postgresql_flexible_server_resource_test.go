// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataProtectionBackupInstancePostgreSQLFlexibleServerResource struct{}

func TestAccDataProtectionBackupInstancePostgreSQLFlexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql_flexible_server", "test")
	r := DataProtectionBackupInstancePostgreSQLFlexibleServerResource{}

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

func TestAccDataProtectionBackupInstancePostgreSQLFlexibleServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql_flexible_server", "test")
	r := DataProtectionBackupInstancePostgreSQLFlexibleServerResource{}

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

func TestAccDataProtectionBackupInstancePostgreSQLFlexibleServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql_flexible_server", "test")
	r := DataProtectionBackupInstancePostgreSQLFlexibleServerResource{}

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

func TestAccDataProtectionBackupInstancePostgreSQLFlexibleServer_keyVaultAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql_flexible_server", "test")
	r := DataProtectionBackupInstancePostgreSQLFlexibleServerResource{}

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

func TestAccDataProtectionBackupInstancePostgreSQLFlexibleServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_postgresql_flexible_server", "test")
	r := DataProtectionBackupInstancePostgreSQLFlexibleServerResource{}

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

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%[1]d"
  location = "%[2]s"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-postgresqlfs-%[1]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}

resource "azurerm_postgresql_flexible_server_firewall_rule" "test" {
  name             = "AllowAllWindowsAzureIps"
  server_id        = azurerm_postgresql_flexible_server.test.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "0.0.0.0"
}

resource "azurerm_postgresql_flexible_server_database" "test" {
  name      = "acctest-postgresql-database-%[1]d"
  server_id = azurerm_postgresql_flexible_server.test.id
  charset   = "UTF8"
  collation = "English_United States.1252"
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
  value        = "Server=${azurerm_postgresql_flexible_server.test.name}.postgres.database.azure.com;Database=${azurerm_postgresql_flexible_server_database.test.name};Port=5432;User Id=psqladmin@${azurerm_postgresql_flexible_server.test.name};Password=H@Sh1CoR3!;Ssl Mode=Require;"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_postgresql_flexible_server.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity.0.principal_id
}

resource "azurerm_data_protection_backup_policy_postgresql_flexible_server" "test" {
  name                            = "acctest-dp-%[1]d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P4M"
}

resource "azurerm_data_protection_backup_policy_postgresql_flexible_server" "another" {
  name                            = "acctest-dp-second-%[1]d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P3M"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(16))
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "test" {
  name             = "acctest-dbi-%d"
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_vault.test.id
  database_id      = azurerm_postgresql_flexible_server_database.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql_flexible_server.test.id
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "import" {
  name             = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.name
  location         = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.location
  vault_id         = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.vault_id
  database_id      = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.database_id
  backup_policy_id = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.backup_policy_id
}
`, config)
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "test" {
  name             = "acctest-dbi-%d"
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_vault.test.id
  database_id      = azurerm_postgresql_flexible_server_database.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql_flexible_server.another.id
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) keyVaultAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "test" {
  name                                    = "acctest-dbi-%d"
  location                                = azurerm_resource_group.test.location
  vault_id                                = azurerm_data_protection_backup_vault.test.id
  database_id                             = azurerm_postgresql_flexible_server_database.test.id
  backup_policy_id                        = azurerm_data_protection_backup_policy_postgresql_flexible_server.another.id
  database_credential_key_vault_secret_id = azurerm_key_vault_secret.test.versionless_id
}
`, template, data.RandomInteger)
}
