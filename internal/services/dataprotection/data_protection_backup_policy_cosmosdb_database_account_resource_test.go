// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupPolicyCosmosdbDatabaseAccountResource struct{}

func TestAccDataProtectionBackupPolicyCosmosDBDatabaseAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_cosmosdb_database_account", "test")
	r := DataProtectionBackupPolicyCosmosdbDatabaseAccountResource{}
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

func TestAccDataProtectionBackupPolicyCosmosDBDatabaseAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_cosmosdb_database_account", "test")
	r := DataProtectionBackupPolicyCosmosdbDatabaseAccountResource{}
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

func TestAccDataProtectionBackupPolicyCosmosDBDatabaseAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_cosmosdb_database_account", "test")
	r := DataProtectionBackupPolicyCosmosdbDatabaseAccountResource{}
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

func (r DataProtectionBackupPolicyCosmosdbDatabaseAccountResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := basebackuppolicyresources.ParseBackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupPolicyClient20260601.BackupPoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupPolicyCosmosdbDatabaseAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dp-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dbv-cosmos-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DataProtectionBackupPolicyCosmosdbDatabaseAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_cosmosdb_database_account" "test" {
  name                            = "acctest-dbp-cosmos-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2026-02-08T10:00:00+00:00/P1W"]
  time_zone                       = "UTC"

  default_retention_rule {
    life_cycle {
      duration        = "P10Y"
      data_store_type = "VaultStore"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupPolicyCosmosdbDatabaseAccountResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_cosmosdb_database_account" "import" {
  name                            = azurerm_data_protection_backup_policy_cosmosdb_database_account.test.name
  vault_id                        = azurerm_data_protection_backup_policy_cosmosdb_database_account.test.vault_id
  backup_repeating_time_intervals = azurerm_data_protection_backup_policy_cosmosdb_database_account.test.backup_repeating_time_intervals
  time_zone                       = azurerm_data_protection_backup_policy_cosmosdb_database_account.test.time_zone

  default_retention_rule {
    life_cycle {
      duration        = "P10Y"
      data_store_type = "VaultStore"
    }
  }
}
`, r.basic(data))
}

func (r DataProtectionBackupPolicyCosmosdbDatabaseAccountResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_cosmosdb_database_account" "test" {
  name                            = "acctest-dbp-cosmos-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2026-02-08T10:00:00+00:00/P1W"]
  time_zone                       = "UTC"

  default_retention_rule {
    life_cycle {
      duration        = "P10Y"
      data_store_type = "VaultStore"
    }
  }

  retention_rule {
    name     = "Monthly"
    priority = 10

    criteria {
      absolute_criteria = "FirstOfMonth"
    }

    life_cycle {
      duration        = "P10Y"
      data_store_type = "VaultStore"
    }
  }
}
`, r.template(data), data.RandomInteger)
}
