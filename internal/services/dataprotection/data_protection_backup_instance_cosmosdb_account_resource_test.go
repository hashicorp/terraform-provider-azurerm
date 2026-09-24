// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupInstanceCosmosdbAccountResource struct{}

func TestAccDataProtectionBackupInstanceCosmosDBAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_cosmosdb_account", "test")
	r := DataProtectionBackupInstanceCosmosdbAccountResource{}
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

func TestAccDataProtectionBackupInstanceCosmosDBAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_cosmosdb_account", "test")
	r := DataProtectionBackupInstanceCosmosdbAccountResource{}
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

func TestAccDataProtectionBackupInstanceCosmosDBAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_cosmosdb_account", "test")
	r := DataProtectionBackupInstanceCosmosdbAccountResource{}
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

func TestAccDataProtectionBackupInstanceCosmosDBAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_cosmosdb_account", "test")
	r := DataProtectionBackupInstanceCosmosdbAccountResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DataProtectionBackupInstanceCosmosdbAccountResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstanceresources.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient20260601.BackupInstancesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupInstanceCosmosdbAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctestcosmos%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Session"
  }

  backup {
    type = "Continuous"
    tier = "Continuous7Days"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_data_protection_backup_policy_cosmosdb_account" "another" {
  name                            = "acctest-dbp-cosmos-other-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  backup_schedule                 = ["R/2026-02-09T10:00:00+00:00/P1W"]
  default_retention_duration      = "P5Y"
  time_zone                       = "UTC"
}

resource "azurerm_role_assignment" "reader" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "cosmos_operator" {
  scope                = azurerm_cosmosdb_account.test.id
  role_definition_name = "Cosmos DB Operator"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}
`, DataProtectionBackupPolicyCosmosDBAccountResource{}.basic(data), data.RandomInteger, data.RandomInteger)
}

func (r DataProtectionBackupInstanceCosmosdbAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_cosmosdb_account" "test" {
  name                              = "acctest-dbi-cosmos-%d"
  location                          = azurerm_resource_group.test.location
  data_protection_backup_vault_id   = azurerm_data_protection_backup_vault.test.id
  backup_policy_cosmosdb_account_id = azurerm_data_protection_backup_policy_cosmosdb_account.test.id
  cosmosdb_account_id               = azurerm_cosmosdb_account.test.id

  depends_on = [
    azurerm_role_assignment.reader,
    azurerm_role_assignment.cosmos_operator,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupInstanceCosmosdbAccountResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_cosmosdb_account" "import" {
  name                              = azurerm_data_protection_backup_instance_cosmosdb_account.test.name
  location                          = azurerm_data_protection_backup_instance_cosmosdb_account.test.location
  data_protection_backup_vault_id   = azurerm_data_protection_backup_instance_cosmosdb_account.test.data_protection_backup_vault_id
  backup_policy_cosmosdb_account_id = azurerm_data_protection_backup_instance_cosmosdb_account.test.backup_policy_cosmosdb_account_id
  cosmosdb_account_id               = azurerm_data_protection_backup_instance_cosmosdb_account.test.cosmosdb_account_id
}
`, r.basic(data))
}

func (r DataProtectionBackupInstanceCosmosdbAccountResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_cosmosdb_account" "test" {
  name                              = "acctest-dbi-cosmos-%d"
  location                          = azurerm_resource_group.test.location
  data_protection_backup_vault_id   = azurerm_data_protection_backup_vault.test.id
  backup_policy_cosmosdb_account_id = azurerm_data_protection_backup_policy_cosmosdb_account.another.id
  cosmosdb_account_id               = azurerm_cosmosdb_account.test.id

  depends_on = [
    azurerm_role_assignment.reader,
    azurerm_role_assignment.cosmos_operator,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupInstanceCosmosdbAccountResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_cosmosdb_account" "test" {
  name                              = "acctest-dbi-cosmos-%d"
  location                          = azurerm_resource_group.test.location
  data_protection_backup_vault_id   = azurerm_data_protection_backup_vault.test.id
  backup_policy_cosmosdb_account_id = azurerm_data_protection_backup_policy_cosmosdb_account.test.id
  cosmosdb_account_id               = azurerm_cosmosdb_account.test.id

  depends_on = [
    azurerm_role_assignment.reader,
    azurerm_role_assignment.cosmos_operator,
  ]
}
`, r.template(data), data.RandomInteger)
}
