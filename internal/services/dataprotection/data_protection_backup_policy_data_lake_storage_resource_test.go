// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupPolicyDataLakeStorageResource struct{}

func TestAccDataProtectionBackupPolicyDataLakeStorage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_data_lake_storage", "test")
	r := DataProtectionBackupPolicyDataLakeStorageResource{}

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

func TestAccDataProtectionBackupPolicyDataLakeStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_data_lake_storage", "test")
	r := DataProtectionBackupPolicyDataLakeStorageResource{}

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

func TestAccDataProtectionBackupPolicyDataLakeStorage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_data_lake_storage", "test")
	r := DataProtectionBackupPolicyDataLakeStorageResource{}

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

func (r DataProtectionBackupPolicyDataLakeStorageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := basebackuppolicyresources.ParseBackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupPolicyClient.BackupPoliciesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupPolicyDataLakeStorageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dbv-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DataProtectionBackupPolicyDataLakeStorageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_data_protection_backup_policy_data_lake_storage" "test" {
  name                            = "acctest-dbp-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  backup_schedule                 = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_duration = "P4M"
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupPolicyDataLakeStorageResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_data_lake_storage" "import" {
  name                            = azurerm_data_protection_backup_policy_data_lake_storage.test.name
  data_protection_backup_vault_id = azurerm_data_protection_backup_policy_data_lake_storage.test.data_protection_backup_vault_id
  backup_schedule                 = azurerm_data_protection_backup_policy_data_lake_storage.test.backup_schedule
  default_retention_duration      = azurerm_data_protection_backup_policy_data_lake_storage.test.default_retention_duration
}
`, r.basic(data))
}

func (r DataProtectionBackupPolicyDataLakeStorageResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_data_protection_backup_policy_data_lake_storage" "test" {
  name                            = "acctest-dbp-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  backup_schedule                 = ["R/2021-05-23T02:30:00+00:00/P1W", "R/2021-05-24T03:40:00+00:00/P1W"]
  time_zone                       = "Coordinated Universal Time"

  default_retention_duration = "P4M"

  retention_rule {
    name              = "weekly"
    duration          = "P6M"
    absolute_criteria = "FirstOfWeek"
  }

  retention_rule {
    name                   = "thursday"
    duration               = "P1W"
    days_of_week           = ["Thursday", "Friday"]
    months_of_year         = ["November", "December"]
    scheduled_backup_times = ["2021-05-23T02:30:00Z"]
  }

  retention_rule {
    name                   = "monthly"
    duration               = "P1D"
    weeks_of_month         = ["First", "Last"]
    days_of_week           = ["Tuesday"]
    scheduled_backup_times = ["2021-05-23T02:30:00Z", "2021-05-24T03:40:00Z"]
  }
}
`, r.template(data), data.RandomInteger)
}
