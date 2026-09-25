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

type DataProtectionBackupPolicyCosmosdbAccountResource struct{}

func TestAccDataProtectionBackupPolicyCosmosdbAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_cosmosdb_account", "test")
	r := DataProtectionBackupPolicyCosmosdbAccountResource{}
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

func TestAccDataProtectionBackupPolicyCosmosdbAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_cosmosdb_account", "test")
	r := DataProtectionBackupPolicyCosmosdbAccountResource{}
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

func TestAccDataProtectionBackupPolicyCosmosdbAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_cosmosdb_account", "test")
	r := DataProtectionBackupPolicyCosmosdbAccountResource{}
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

func (r DataProtectionBackupPolicyCosmosdbAccountResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r DataProtectionBackupPolicyCosmosdbAccountResource) template(data acceptance.TestData) string {
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

func (r DataProtectionBackupPolicyCosmosdbAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_cosmosdb_account" "test" {
  name                            = "acctest-dbp-cosmos-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  backup_schedule                 = ["R/2026-02-08T10:00:00+00:00/P1W"]
  default_retention_duration      = "P10Y"
  time_zone                       = "UTC"
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_cosmosdb_account" "import" {
  name                            = azurerm_data_protection_backup_policy_cosmosdb_account.test.name
  data_protection_backup_vault_id = azurerm_data_protection_backup_policy_cosmosdb_account.test.data_protection_backup_vault_id
  backup_schedule                 = azurerm_data_protection_backup_policy_cosmosdb_account.test.backup_schedule
  default_retention_duration      = azurerm_data_protection_backup_policy_cosmosdb_account.test.default_retention_duration
  time_zone                       = azurerm_data_protection_backup_policy_cosmosdb_account.test.time_zone
}
`, r.basic(data))
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_cosmosdb_account" "test" {
  name                            = "acctest-dbp-cosmos-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  backup_schedule                 = ["R/2026-02-08T10:00:00+00:00/P1W"]
  default_retention_duration      = "P10Y"
  time_zone                       = "UTC"

  retention_rule {
    name              = "Monthly"
    duration          = "P10Y"
    absolute_criteria = "FirstOfMonth"
  }

  retention_rule {
    name                   = "Weekly"
    duration               = "P4M"
    days_of_week           = ["Sunday"]
    scheduled_backup_times = ["2026-02-08T10:00:00Z"]
  }

  retention_rule {
    name                   = "Yearly"
    duration               = "P1Y"
    days_of_week           = ["Monday"]
    months_of_year         = ["January"]
    weeks_of_month         = ["First"]
    scheduled_backup_times = ["2026-02-08T10:00:00Z"]
  }
}
`, r.template(data), data.RandomInteger)
}
