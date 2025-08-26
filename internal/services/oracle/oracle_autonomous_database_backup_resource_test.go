// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutonomousDatabaseBackupResource struct{}

func TestAutonomousDatabaseBackupResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseBackupResource{}.ResourceType(), "test")
	r := AutonomousDatabaseBackupResource{}
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

func TestAutonomousDatabaseBackupResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseBackupResource{}.ResourceType(), "test")
	r := AutonomousDatabaseBackupResource{}
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

func TestAutonomousDatabaseBackupResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseBackupResource{}.ResourceType(), "test")
	r := AutonomousDatabaseBackupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAutonomousDatabaseBackupResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseBackupResource{}.ResourceType(), "test")
	r := AutonomousDatabaseBackupResource{}
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

func (a AutonomousDatabaseBackupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabasebackups.ParseAutonomousDatabaseBackupID(state.ID)
	if err != nil {
		return nil, err
	}

	adbId := autonomousdatabasebackups.NewAutonomousDatabaseID(
		id.SubscriptionId,
		id.ResourceGroupName,
		id.AutonomousDatabaseName,
	)
	backupId := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(id.SubscriptionId, id.ResourceGroupName, id.AutonomousDatabaseName, id.AutonomousDatabaseBackupName)

	backup, err := getBackupFromOCI(ctx, client.Oracle.OracleClient.AutonomousDatabaseBackups, autonomousdatabases.AutonomousDatabaseId(adbId), backupId)
	if err != nil {
		return nil, fmt.Errorf("checking backup existence: %+v", err)
	}

	return pointer.To(backup != nil), nil
}

func getBackupFromOCI(ctx context.Context, client *autonomousdatabasebackups.AutonomousDatabaseBackupsClient, adbId autonomousdatabases.AutonomousDatabaseId, backupId autonomousdatabasebackups.AutonomousDatabaseBackupId) (*autonomousdatabasebackups.AutonomousDatabaseBackup, error) {
	resp, err := client.ListByParent(ctx, autonomousdatabasebackups.AutonomousDatabaseId(adbId))
	if err != nil {
		return nil, fmt.Errorf("listing backups for %s: %+v", adbId.ID(), err)
	}

	id := backupId.ID()

	if model := resp.Model; model != nil {
		for _, backup := range *model {
			if backup.Id != nil && strings.EqualFold(*backup.Id, id) {
				return &backup, nil
			}
		}
	}

	return nil, nil
}

func (a AutonomousDatabaseBackupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_backup" "test" {
  name                     = "backup%[2]d"
  autonomous_database_id   = azurerm_oracle_autonomous_database.test.id
  retention_period_in_days = 120
  type                     = "LongTerm"
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a AutonomousDatabaseBackupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s


resource "azurerm_oracle_autonomous_database_backup" "test" {
  name                     = "backup%[2]d"
  autonomous_database_id   = azurerm_oracle_autonomous_database.test.id
  retention_period_in_days = 120
  type                     = "LongTerm"
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a AutonomousDatabaseBackupResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_backup" "test" {
  name                     = "backup%[2]d"
  autonomous_database_id   = azurerm_oracle_autonomous_database.test.id
  retention_period_in_days = 160
  type                     = "LongTerm"
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a AutonomousDatabaseBackupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database_backup" "import" {
  name                     = azurerm_oracle_autonomous_database_backup.test.name
  autonomous_database_id   = azurerm_oracle_autonomous_database_backup.test.autonomous_database_id
  retention_period_in_days = azurerm_oracle_autonomous_database_backup.test.retention_period_in_days
}
`, a.basic(data))
}

func (a AutonomousDatabaseBackupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-oadbb-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "eacctest%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Oracle.Database/networkAttachments"
    }
  }
}

resource "azurerm_oracle_autonomous_database" "test" {
  name                             = "OFake%[1]d"
  display_name                     = "OFake%[1]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = "%[2]s"
  compute_model                    = "ECPU"
  compute_count                    = 2
  license_model                    = "BringYourOwnLicense"
  backup_retention_period_in_days  = 12
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  data_storage_size_in_tbs         = 1
  db_workload                      = "OLTP"
  admin_password                   = "TestPass#2024#"
  db_version                       = "19c"
  character_set                    = "AL32UTF8"
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.test.id
  virtual_network_id               = azurerm_virtual_network.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
