// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataProtectionBackupPolicyBlobStorageResource struct{}

func TestAccDataProtectionBackupPolicyBlobStorage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_blob_storage", "test")
	r := DataProtectionBackupPolicyBlobStorageResource{}
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

func TestAccDataProtectionBackupPolicyBlobStorage_vaultbackup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_blob_storage", "test")
	r := DataProtectionBackupPolicyBlobStorageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vaultBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupPolicyBlobStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_blob_storage", "test")
	r := DataProtectionBackupPolicyBlobStorageResource{}
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

func (r DataProtectionBackupPolicyBlobStorageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backuppolicies.ParseBackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupPolicyClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving DataProtection BackupPolicy (%q): %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r DataProtectionBackupPolicyBlobStorageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

func (r DataProtectionBackupPolicyBlobStorageResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_blob_storage" "test" {
  name                                   = "acctest-dbp-%d"
  vault_id                               = azurerm_data_protection_backup_vault.test.id
  operational_default_retention_duration = "P30D"
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupPolicyBlobStorageResource) vaultBackup(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_blob_storage" "test" {
  name                                   = "acctest-dbp-%d"
  vault_id                               = azurerm_data_protection_backup_vault.test.id
  operational_default_retention_duration = "P30D"
  vault_default_retention_duration       = "P7D"
  backup_repeating_time_intervals        = ["R/2024-05-08T11:30:00+00:00/P1W"]

  retention_rule {
    name     = "Monthly"
    priority = 15
    life_cycle {
      duration        = "P6M"
      data_store_type = "VaultStore"
    }
    criteria {
      days_of_month = [1, 2, 0]
    }
  }

  retention_rule {
    name     = "Daily"
    priority = 25
    life_cycle {
      duration        = "P7D"
      data_store_type = "VaultStore"
    }
    criteria {
      days_of_week           = ["Thursday"]
      months_of_year         = ["November"]
      weeks_of_month         = ["First"]
      scheduled_backup_times = ["2024-05-08T02:30:00Z"]
    }
  }
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupPolicyBlobStorageResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_blob_storage" "import" {
  name                                   = azurerm_data_protection_backup_policy_blob_storage.test.name
  vault_id                               = azurerm_data_protection_backup_policy_blob_storage.test.vault_id
  operational_default_retention_duration = "P30D"
}
`, config)
}
