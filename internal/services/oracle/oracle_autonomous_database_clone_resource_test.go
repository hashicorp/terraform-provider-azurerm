// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutonomousDatabaseCloneResource struct{}

func TestAccAutonomousDatabaseClone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("ADB%dclone", data.RandomInteger)),
				check.That(data.ResourceName).Key("data_base_type").HasValue("Clone"),
			),
		},
		data.ImportStep("admin_password", "source", "clone_type"),
	})
}

func TestAccAutonomousDatabaseClone_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

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

func TestAccAutonomousDatabaseClone_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "source", "clone_type", "refreshable_model"),
	})
}

func TestAccAutonomousDatabaseClone_metadataClone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metadataClone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("clone_type").HasValue("Metadata"),
			),
		},
		data.ImportStep("admin_password", "source", "clone_type"),
	})
}

func TestAccAutonomousDatabaseClone_backupTimestampLatestBackup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupTimestampCloneLatest(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("data_base_type").HasValue("CloneFromBackupTimestamp"),
			),
		},
		data.ImportStep("admin_password", "source", "clone_type", "use_latest_available_backup_time_stamp"),
	})
}

func (r AutonomousDatabaseCloneResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Oracle.OracleClient.AutonomousDatabases
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r AutonomousDatabaseCloneResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_id      = azurerm_oracle_autonomous_database.test.id
  clone_type     = "Full"
  source         = "Database"
  data_base_type = "Clone"

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "ADB%[2]dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = true
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_oracle_autonomous_database.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database.test.virtual_network_id

  tags = {
    Environment = "Test"
    Purpose     = "BasicClone"
  }

  depends_on = [azurerm_oracle_autonomous_database.test]
}
`, AdbsRegularResource{}.basic(data), data.RandomInteger)
}

func (r AutonomousDatabaseCloneResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "import" {
  name                             = azurerm_oracle_autonomous_database_clone.test.name
  resource_group_name              = azurerm_oracle_autonomous_database_clone.test.resource_group_name
  location                         = azurerm_oracle_autonomous_database_clone.test.location
  source_id                        = azurerm_oracle_autonomous_database_clone.test.source_id
  clone_type                       = azurerm_oracle_autonomous_database_clone.test.clone_type
  source                           = azurerm_oracle_autonomous_database_clone.test.source
  data_base_type                   = azurerm_oracle_autonomous_database_clone.test.data_base_type
  admin_password                   = azurerm_oracle_autonomous_database_clone.test.admin_password
  backup_retention_period_in_days  = azurerm_oracle_autonomous_database_clone.test.backup_retention_period_in_days
  character_set                    = azurerm_oracle_autonomous_database_clone.test.character_set
  compute_count                    = azurerm_oracle_autonomous_database_clone.test.compute_count
  compute_model                    = azurerm_oracle_autonomous_database_clone.test.compute_model
  data_storage_size_in_tbs         = azurerm_oracle_autonomous_database_clone.test.data_storage_size_in_tbs
  db_version                       = azurerm_oracle_autonomous_database_clone.test.db_version
  db_workload                      = azurerm_oracle_autonomous_database_clone.test.db_workload
  display_name                     = azurerm_oracle_autonomous_database_clone.test.display_name
  license_model                    = azurerm_oracle_autonomous_database_clone.test.license_model
  auto_scaling_enabled             = azurerm_oracle_autonomous_database_clone.test.auto_scaling_enabled
  auto_scaling_for_storage_enabled = azurerm_oracle_autonomous_database_clone.test.auto_scaling_for_storage_enabled
  mtls_connection_required         = azurerm_oracle_autonomous_database_clone.test.mtls_connection_required
  national_character_set           = azurerm_oracle_autonomous_database_clone.test.national_character_set
  subnet_id                        = azurerm_oracle_autonomous_database_clone.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database_clone.test.virtual_network_id
  tags                             = azurerm_oracle_autonomous_database_clone.test.tags
}
`, r.basic(data))
}

func (r AutonomousDatabaseCloneResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_id      = azurerm_oracle_autonomous_database.test.id
  clone_type     = "Full"
  source         = "Database"
  data_base_type = "Clone"

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 15
  character_set                    = "AL32UTF8"
  compute_count                    = 4.0
  compute_model                    = "ECPU"
  data_storage_size_in_tbs         = 2
  db_version                       = "19c"
  db_workload                      = "DW"
  display_name                     = "ADB%[2]dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = true
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_oracle_autonomous_database.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database.test.virtual_network_id

  # Clone-specific optional fields
  refreshable_model = "Manual"

  customer_contacts = ["test@example.com"]

  tags = {
    Environment = "Test"
    Purpose     = "CompleteClone"
    Type        = "Full"
  }
  depends_on = [azurerm_oracle_autonomous_database.test]
}
`, AdbsRegularResource{}.basic(data), data.RandomInteger)
}

func (r AutonomousDatabaseCloneResource) metadataClone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_id      = azurerm_oracle_autonomous_database.test.id
  clone_type     = "Metadata"
  source         = "Database"
  data_base_type = "Clone"

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
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
    Purpose     = "MetadataClone"
  }
}
`, AdbsRegularResource{}.basic(data), data.RandomInteger)
}

func (r AutonomousDatabaseCloneResource) backupTimestampCloneLatest(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_id      = azurerm_oracle_autonomous_database.test.id
  clone_type     = "Full"
  source         = "BackupFromTimestamp"
  data_base_type = "CloneFromBackupTimestamp"

  # Use latest backup
  use_latest_available_backup_time_stamp = true

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
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
