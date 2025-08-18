// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseCloneFromBackupResource struct{}

func TestAccAutonomousDatabaseCloneFromBackup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone_from_backup", "test")
	r := AutonomousDatabaseCloneFromBackupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "source", "clone_type", "use_latest_available_backup_time_stamp"),
	})
}

func TestAccAutonomousDatabaseCloneFromBackup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone_from_backup", "test")
	r := AutonomousDatabaseCloneFromBackupResource{}

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

func TestAccAutonomousDatabaseCloneFromBackup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone_from_backup", "test")
	r := AutonomousDatabaseCloneFromBackupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "source", "clone_type", "use_latest_available_backup_time_stamp"),
	})
}

func (r AutonomousDatabaseCloneFromBackupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Oracle.OracleClient.AutonomousDatabases.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r AutonomousDatabaseCloneFromBackupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone_from_backup" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  clone_type                    = "Full"
  source                        = "BackupFromTimestamp"

  use_latest_available_backup_time_stamp = true

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tb          = 1
  database_version                 = "19c"
  db_workload                      = "OLTP"
  display_name                     = "ADB%[2]dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = true
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_oracle_autonomous_database.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database.test.virtual_network_id

  tags = {
    Environment = "Test"
    Purpose     = "BackupTimestampClone"
  }
}
`, AdbsRegularResource{}.basic(data), data.RandomInteger)
}

func (r AutonomousDatabaseCloneFromBackupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone_from_backup" "import" {
  name                             = azurerm_oracle_autonomous_database_clone_from_backup.test.name
  resource_group_name              = azurerm_oracle_autonomous_database_clone_from_backup.test.resource_group_name
  location                         = azurerm_oracle_autonomous_database_clone_from_backup.test.location
  source_autonomous_database_id    = azurerm_oracle_autonomous_database_clone_from_backup.test.source_autonomous_database_id
  clone_type                       = azurerm_oracle_autonomous_database_clone_from_backup.test.clone_type
  source                           = azurerm_oracle_autonomous_database_clone_from_backup.test.source
  admin_password                   = azurerm_oracle_autonomous_database_clone_from_backup.test.admin_password
  backup_retention_period_in_days  = azurerm_oracle_autonomous_database_clone_from_backup.test.backup_retention_period_in_days
  character_set                    = azurerm_oracle_autonomous_database_clone_from_backup.test.character_set
  compute_count                    = azurerm_oracle_autonomous_database_clone_from_backup.test.compute_count
  compute_model                    = azurerm_oracle_autonomous_database_clone_from_backup.test.compute_model
  data_storage_size_in_tb          = azurerm_oracle_autonomous_database_clone_from_backup.test.data_storage_size_in_tb
  database_version                 = azurerm_oracle_autonomous_database_clone_from_backup.test.database_version
  db_workload                      = azurerm_oracle_autonomous_database_clone_from_backup.test.db_workload
  display_name                     = azurerm_oracle_autonomous_database_clone_from_backup.test.display_name
  license_model                    = azurerm_oracle_autonomous_database_clone_from_backup.test.license_model
  auto_scaling_enabled             = azurerm_oracle_autonomous_database_clone_from_backup.test.auto_scaling_enabled
  auto_scaling_for_storage_enabled = azurerm_oracle_autonomous_database_clone_from_backup.test.auto_scaling_for_storage_enabled
  mtls_connection_required         = azurerm_oracle_autonomous_database_clone_from_backup.test.mtls_connection_required
  national_character_set           = azurerm_oracle_autonomous_database_clone_from_backup.test.national_character_set
  subnet_id                        = azurerm_oracle_autonomous_database_clone_from_backup.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database_clone_from_backup.test.virtual_network_id
  tags                             = azurerm_oracle_autonomous_database_clone_from_backup.test.tags
}
`, r.basic(data))
}

func (r AutonomousDatabaseCloneFromBackupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone_from_backup" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  clone_type                    = "Metadata"
  source                        = "BackupFromTimestamp"

  # Use latest backup
  use_latest_available_backup_time_stamp = true

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tb          = 1
  database_version                 = "19c"
  db_workload                      = "OLTP"
  display_name                     = "ADB%[2]dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = true
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_oracle_autonomous_database.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database.test.virtual_network_id

  tags = {
    Environment = "Test"
    Purpose     = "BackupTimestampClone"
  }
}
`, AdbsRegularResource{}.basic(data), data.RandomInteger)
}
