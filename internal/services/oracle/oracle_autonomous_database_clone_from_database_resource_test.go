// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseCloneFromDatabaseResource struct{}

func TestAccAutonomousDatabaseCloneFromDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone_from_database", "test")
	r := AutonomousDatabaseCloneFromDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccAutonomousDatabaseCloneFromDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone_from_database", "test")
	r := AutonomousDatabaseCloneFromDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutonomousDatabaseCloneFromDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone_from_database", "test")
	r := AutonomousDatabaseCloneFromDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password", "refreshable_model"),
	})
}

func TestAccAutonomousDatabaseCloneFromDatabase_metadataClone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_oracle_autonomous_database_clone_from_database", "test")
	r := AutonomousDatabaseCloneFromDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metadataClone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("clone_type").HasValue("Metadata"),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func (r AutonomousDatabaseCloneFromDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Oracle.OracleClient.AutonomousDatabases.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r AutonomousDatabaseCloneFromDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone_from_database" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  clone_type                    = "Full"

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tb          = 1
  database_version                 = "19c"
  database_workload                = "OLTP"
  display_name                     = "ADB%[2]dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = true
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  allowed_ip_addresses             = ["140.204.126.129", "140.204.125.0/24"]

  depends_on = [azurerm_oracle_autonomous_database.test]
}
`, AdbsRegularResource{}.publicAccess(data), data.RandomInteger)
}

func (r AutonomousDatabaseCloneFromDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone_from_database" "import" {
  name                             = azurerm_oracle_autonomous_database_clone_from_database.test.name
  resource_group_name              = azurerm_oracle_autonomous_database_clone_from_database.test.resource_group_name
  location                         = azurerm_oracle_autonomous_database_clone_from_database.test.location
  source_autonomous_database_id    = azurerm_oracle_autonomous_database_clone_from_database.test.source_autonomous_database_id
  clone_type                       = azurerm_oracle_autonomous_database_clone_from_database.test.clone_type
  admin_password                   = azurerm_oracle_autonomous_database_clone_from_database.test.admin_password
  backup_retention_period_in_days  = azurerm_oracle_autonomous_database_clone_from_database.test.backup_retention_period_in_days
  character_set                    = azurerm_oracle_autonomous_database_clone_from_database.test.character_set
  compute_count                    = azurerm_oracle_autonomous_database_clone_from_database.test.compute_count
  compute_model                    = azurerm_oracle_autonomous_database_clone_from_database.test.compute_model
  data_storage_size_in_tb          = azurerm_oracle_autonomous_database_clone_from_database.test.data_storage_size_in_tb
  database_version                 = azurerm_oracle_autonomous_database_clone_from_database.test.database_version
  database_workload                = azurerm_oracle_autonomous_database_clone_from_database.test.database_workload
  display_name                     = azurerm_oracle_autonomous_database_clone_from_database.test.display_name
  license_model                    = azurerm_oracle_autonomous_database_clone_from_database.test.license_model
  auto_scaling_enabled             = azurerm_oracle_autonomous_database_clone_from_database.test.auto_scaling_enabled
  auto_scaling_for_storage_enabled = azurerm_oracle_autonomous_database_clone_from_database.test.auto_scaling_for_storage_enabled
  mtls_connection_required         = azurerm_oracle_autonomous_database_clone_from_database.test.mtls_connection_required
  national_character_set           = azurerm_oracle_autonomous_database_clone_from_database.test.national_character_set
  subnet_id                        = azurerm_oracle_autonomous_database_clone_from_database.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database_clone_from_database.test.virtual_network_id
  tags                             = azurerm_oracle_autonomous_database_clone_from_database.test.tags
}
`, r.complete(data))
}

func (r AutonomousDatabaseCloneFromDatabaseResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone_from_database" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  clone_type                    = "Full"

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 15
  character_set                    = "AL32UTF8"
  compute_count                    = 4.0
  compute_model                    = "ECPU"
  data_storage_size_in_tb          = 2
  database_version                 = "19c"
  database_workload                = "DW"
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

func (r AutonomousDatabaseCloneFromDatabaseResource) metadataClone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_clone_from_database" "test" {
  name                = "ADB%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  clone_type                    = "Metadata"

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tb          = 1
  database_version                 = "19c"
  database_workload                = "OLTP"
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
