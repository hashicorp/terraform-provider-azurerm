// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataProtectionBackupPolicyMySQLFlexibleServerResource struct{}

func TestAccDataProtectionBackupPolicyMySQLFlexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_mysql_flexible_server", "test")
	r := DataProtectionBackupPolicyMySQLFlexibleServerResource{}

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

func TestAccDataProtectionBackupPolicyMySQLFlexibleServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_mysql_flexible_server", "test")
	r := DataProtectionBackupPolicyMySQLFlexibleServerResource{}

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

func TestAccDataProtectionBackupPolicyMySQLFlexibleServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_mysql_flexible_server", "test")
	r := DataProtectionBackupPolicyMySQLFlexibleServerResource{}

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

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backuppolicies.ParseBackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupPolicyClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) template(data acceptance.TestData) string {
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

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_mysql_flexible_server" "test" {
  name                            = "acctest-dbp-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_mysql_flexible_server" "import" {
  name                            = azurerm_data_protection_backup_policy_mysql_flexible_server.test.name
  vault_id                        = azurerm_data_protection_backup_policy_mysql_flexible_server.test.vault_id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }
}
`, r.basic(data))
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_mysql_flexible_server" "test" {
  name                            = "acctest-dbp-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W", "R/2021-05-24T03:40:00+00:00/P1W"]
  time_zone                       = "India Standard Time"

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }

  retention_rule {
    name     = "weekly"
    priority = 20

    life_cycle {
      duration        = "P6M"
      data_store_type = "VaultStore"
    }

    criteria {
      absolute_criteria = "FirstOfWeek"
    }
  }

  retention_rule {
    name     = "thursday"
    priority = 25

    life_cycle {
      duration        = "P1W"
      data_store_type = "VaultStore"
    }

    criteria {
      days_of_week           = ["Thursday", "Friday"]
      months_of_year         = ["November", "December"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }

  retention_rule {
    name     = "monthly"
    priority = 30

    life_cycle {
      duration        = "P1D"
      data_store_type = "VaultStore"
    }

    criteria {
      weeks_of_month         = ["First", "Last"]
      days_of_week           = ["Tuesday"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z", "2021-05-24T03:40:00Z"]
    }
  }
}
`, r.template(data), data.RandomInteger)
}
