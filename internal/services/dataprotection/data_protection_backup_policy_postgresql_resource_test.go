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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataProtectionBackupPolicyPostgreSQLResource struct{}

func TestAccDataProtectionBackupPolicyPostgreSQL_basicDeprecatedInV4(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("this test requires 3.0 mode")
	}
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSQLResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecatedInV4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"default_retention_duration",
			"default_retention_rule.#",
			"default_retention_rule.0.%",
			"default_retention_rule.0.life_cycle.#",
			"default_retention_rule.0.life_cycle.0.%",
			"default_retention_rule.0.life_cycle.0.data_store_type",
			"default_retention_rule.0.life_cycle.0.duration",
		),
	})
}

func TestAccDataProtectionBackupPolicyPostgreSQL_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSQLResource{}
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

func TestAccDataProtectionBackupPolicyPostgreSQL_requiresImportDeprecatedInV4(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("this test requires 3.0 mode")
	}
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSQLResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecatedInV4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImportDeprecatedInV4),
	})
}

func TestAccDataProtectionBackupPolicyPostgreSQL_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSQLResource{}
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

func TestAccDataProtectionBackupPolicyPostgreSQL_completeDeprecatedInV4(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skip("this test requires 3.0 mode")
	}
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSQLResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeDeprecatedInV4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"default_retention_duration",
			"default_retention_rule.#",
			"default_retention_rule.0.%",
			"default_retention_rule.0.life_cycle.#",
			"default_retention_rule.0.life_cycle.0.%",
			"default_retention_rule.0.life_cycle.0.data_store_type",
			"default_retention_rule.0.life_cycle.0.duration",
			"retention_rule.0.duration",
			"retention_rule.0.life_cycle.#",
			"retention_rule.0.life_cycle.0.%",
			"retention_rule.0.life_cycle.0.data_store_type",
			"retention_rule.0.life_cycle.0.duration",
			"retention_rule.1.duration",
			"retention_rule.1.life_cycle.#",
			"retention_rule.1.life_cycle.0.%",
			"retention_rule.1.life_cycle.0.data_store_type",
			"retention_rule.1.life_cycle.0.duration",
			"retention_rule.2.duration",
			"retention_rule.2.life_cycle.#",
			"retention_rule.2.life_cycle.0.%",
			"retention_rule.2.life_cycle.0.data_store_type",
			"retention_rule.2.life_cycle.0.duration",
		),
	})
}

func TestAccDataProtectionBackupPolicyPostgreSQL_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSQLResource{}
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

func (r DataProtectionBackupPolicyPostgreSQLResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r DataProtectionBackupPolicyPostgreSQLResource) template(data acceptance.TestData) string {
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

func (r DataProtectionBackupPolicyPostgreSQLResource) basicDeprecatedInV4(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "test" {
  name                = "acctest-dbp-%d"
  resource_group_name = azurerm_resource_group.test.name
  vault_name          = azurerm_data_protection_backup_vault.test.name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P4M"
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupPolicyPostgreSQLResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "test" {
  name                = "acctest-dbp-%d"
  resource_group_name = azurerm_resource_group.test.name
  vault_name          = azurerm_data_protection_backup_vault.test.name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupPolicyPostgreSQLResource) requiresImportDeprecatedInV4(data acceptance.TestData) string {
	config := r.basicDeprecatedInV4(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "import" {
  name                = azurerm_data_protection_backup_policy_postgresql.test.name
  resource_group_name = azurerm_data_protection_backup_policy_postgresql.test.resource_group_name
  vault_name          = azurerm_data_protection_backup_policy_postgresql.test.vault_name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P4M"
}
`, config)
}

func (r DataProtectionBackupPolicyPostgreSQLResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "import" {
  name                = azurerm_data_protection_backup_policy_postgresql.test.name
  resource_group_name = azurerm_data_protection_backup_policy_postgresql.test.resource_group_name
  vault_name          = azurerm_data_protection_backup_policy_postgresql.test.vault_name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }
}
`, config)
}

func (r DataProtectionBackupPolicyPostgreSQLResource) completeDeprecatedInV4(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "test" {
  name                = "acctest-dbp-%d"
  resource_group_name = azurerm_resource_group.test.name
  vault_name          = azurerm_data_protection_backup_vault.test.name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P4M"
  time_zone                       = "India Standard Time"

  retention_rule {
    name     = "weekly"
    duration = "P6M"
    priority = 20
    criteria {
      absolute_criteria = "FirstOfWeek"
    }
  }

  retention_rule {
    name     = "thursday"
    duration = "P1W"
    priority = 25
    criteria {
      days_of_week           = ["Thursday"]
      months_of_year         = ["November"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }

  retention_rule {
    name     = "monthly"
    duration = "P1D"
    priority = 30
    criteria {
      weeks_of_month         = ["First", "Last"]
      days_of_week           = ["Tuesday"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupPolicyPostgreSQLResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "test" {
  name                            = "acctest-dbp-%d"
  resource_group_name             = azurerm_resource_group.test.name
  vault_name                      = azurerm_data_protection_backup_vault.test.name
  backup_repeating_time_intervals = ["R/2023-12-31T10:00:00+05:30/P1W"]

  retention_rule {
    name     = "Weekly"
    priority = 30
    life_cycle {
      duration        = "P12W"
      data_store_type = "VaultStore"
      target_copy {
        option_json = jsonencode({
          objectType = "CopyOnExpiryOption"
        })
        data_store_type = "ArchiveStore"
      }
    }
    life_cycle {
      duration        = "P27W"
      data_store_type = "ArchiveStore"
    }
    criteria {
      weeks_of_month         = ["First", "Last"]
      days_of_week           = ["Tuesday"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }
  default_retention_rule {
    life_cycle {
      duration        = "P12M"
      data_store_type = "VaultStore"
      target_copy {
        option_json = jsonencode({
          objectType = "CopyOnExpiryOption"
        })
        data_store_type = "ArchiveStore"
      }
    }
    life_cycle {
      duration        = "P27M"
      data_store_type = "ArchiveStore"
    }
  }
}
`, template, data.RandomInteger)
}
