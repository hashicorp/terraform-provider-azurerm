// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type DataProtectionBackupInstanceProtectAction struct{}

func TestAccDataProtectionBackupInstanceProtectAction_stopProtection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_protect", "test")
	a := DataProtectionBackupInstanceProtectAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.protectActionStopProtection(data),
			},

			{
				RefreshState: true,
				Check:        check.That("azurerm_data_protection_backup_instance_postgresql_flexible_server.test").Key("protection_state").HasValue("ProtectionStopped"),
			},
		},
	})
}

func TestAccDataProtectionBackupInstanceProtectAction_resumeProtection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_protect", "test")
	a := DataProtectionBackupInstanceProtectAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.protectActionResumeProtection(data),
			},

			{
				RefreshState: true,
				Check:        check.That("azurerm_data_protection_backup_instance_postgresql_flexible_server.test").Key("protection_state").HasValue("ProtectionConfigured"),
			},
		},
	})
}

func TestAccDataProtectionBackupInstanceProtectAction_suspendBackups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_protect", "test")
	a := DataProtectionBackupInstanceProtectAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.protectActionSuspendBackups(data),
			},

			{
				RefreshState: true,
				Check:        check.That("azurerm_data_protection_backup_instance_postgresql_flexible_server.test").Key("protection_state").HasValue("BackupsSuspended"),
			},
		},
	})
}

func TestAccDataProtectionBackupInstanceProtectAction_resumeBackups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_protect", "test")
	a := DataProtectionBackupInstanceProtectAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.protectActionResumeBackups(data),
			},

			{
				RefreshState: true,
				Check:        check.That("azurerm_data_protection_backup_instance_postgresql_flexible_server.test").Key("protection_state").HasValue("BackupsResumed"),
			},
		},
	})
}

func (a *DataProtectionBackupInstanceProtectAction) templatePostgres(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-postgresqlfs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  version                = "12"
  sku_name               = "GP_Standard_D4s_v3"
  zone                   = "2"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dataprotection-vault-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "test2" {
  scope                = azurerm_postgresql_flexible_server.test.id
  role_definition_name = "PostgreSQL Flexible Server Long Term Retention Backup Role"
  principal_id         = azurerm_data_protection_backup_vault.test.identity.0.principal_id
}

resource "azurerm_data_protection_backup_policy_postgresql_flexible_server" "test" {
  name                            = "acctest-dp-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }

  depends_on = [azurerm_role_assignment.test, azurerm_role_assignment.test2]
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (a *DataProtectionBackupInstanceProtectAction) protectActionStopProtection(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "test" {
  name             = "acctest-dbi-%d"
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_vault.test.id
  server_id        = azurerm_postgresql_flexible_server.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql_flexible_server.test.id

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_data_protection_backup_instance_protect.stop_protection]
    }
  }
}

action "azurerm_data_protection_backup_instance_protect" "stop_protection" {
  config {
    backup_instance_id = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.id
    protect_action     = "stop_protection"
  }
}
`, a.templatePostgres(data), data.RandomInteger)
}

func (a *DataProtectionBackupInstanceProtectAction) protectActionResumeProtection(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "test" {
  name             = "acctest-dbi-%d"
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_vault.test.id
  server_id        = azurerm_postgresql_flexible_server.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql_flexible_server.test.id

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_data_protection_backup_instance_protect.resume_protection]
    }
  }
}

action "azurerm_data_protection_backup_instance_protect" "resume_protection" {
  config {
    backup_instance_id = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.id
    protect_action     = "resume_protection"
  }
}
`, a.templatePostgres(data), data.RandomInteger)
}

func (a *DataProtectionBackupInstanceProtectAction) protectActionSuspendBackups(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "test" {
  name             = "acctest-dbi-%d"
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_vault.test.id
  server_id        = azurerm_postgresql_flexible_server.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql_flexible_server.test.id

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_data_protection_backup_instance_protect.suspend_backups]
    }
  }
}

action "azurerm_data_protection_backup_instance_protect" "suspend_backups" {
  config {
    backup_instance_id = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.id
    protect_action     = "suspend_backups"
  }
}
`, a.templatePostgres(data), data.RandomInteger)
}

func (a *DataProtectionBackupInstanceProtectAction) protectActionResumeBackups(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "test" {
  name             = "acctest-dbi-%d"
  location         = azurerm_resource_group.test.location
  vault_id         = azurerm_data_protection_backup_vault.test.id
  server_id        = azurerm_postgresql_flexible_server.test.id
  backup_policy_id = azurerm_data_protection_backup_policy_postgresql_flexible_server.test.id

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_data_protection_backup_instance_protect.resume_backups]
    }
  }
}

action "azurerm_data_protection_backup_instance_protect" "resume_backups" {
  config {
    backup_instance_id = azurerm_data_protection_backup_instance_postgresql_flexible_server.test.id
    protect_action     = "resume_backups"
  }
}
`, a.templatePostgres(data), data.RandomInteger)
}
