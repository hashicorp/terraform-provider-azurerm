package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataProtectionBackupPolicyPostgreSqlResource struct{}

func TestAccDataProtectionBackupPolicyPostgreSql_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSqlResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupPolicyPostgreSql_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSqlResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataProtectionBackupPolicyPostgreSql_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSqlResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupPolicyPostgreSql_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_postgresql", "test")
	r := DataProtectionBackupPolicyPostgreSqlResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DataProtectionBackupPolicyPostgreSqlResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.BackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupPolicyClient.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving DataProtection BackupPolicy (%q): %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r DataProtectionBackupPolicyPostgreSqlResource) template(data acceptance.TestData) string {
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

func (r DataProtectionBackupPolicyPostgreSqlResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "test" {
  name                = "acctest-dbp-%d"
  resource_group_name = azurerm_resource_group.test.name
  vault_name          = azurerm_data_protection_backup_vault.test.name

  backup_rules {
    name                     = "backup"
    repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  }
  default_retention_duration = "P4M"
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupPolicyPostgreSqlResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "import" {
  name                = azurerm_data_protection_backup_policy_postgresql.test.name
  resource_group_name = azurerm_data_protection_backup_policy_postgresql.test.resource_group_name
  vault_name          = azurerm_data_protection_backup_policy_postgresql.test.vault_name
  backup_rules {
    name                     = "backup"
    repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  }
  default_retention_duration = "P4M"
}
`, config)
}

func (r DataProtectionBackupPolicyPostgreSqlResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "test" {
  name                = "acctest-dbp-%d"
  resource_group_name = azurerm_resource_group.test.name
  vault_name          = azurerm_data_protection_backup_vault.test.name

  backup_rules {
    name                     = "backup"
    repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  }
  default_retention_duration = "P4M"
  retention_rules {
    name             = "weekly"
    duration         = "P6M"
    tagging_priority = 20
    tagging_criteria {
      absolute_criteria = "FirstOfWeek"
    }
  }

  retention_rules {
    name             = "thursday"
    duration         = "P1W"
    tagging_priority = 25
    tagging_criteria {
      days_of_the_week = ["Thursday"]
      schedule_times   = ["2021-05-23T02:30:00Z"]
    }
  }

  retention_rules {
    name             = "monthly"
    duration         = "P1D"
    tagging_priority = 30
    tagging_criteria {
      weeks_of_the_month = ["First", "Last"]
      days_of_the_week   = ["Tuesday"]
      schedule_times     = ["2021-05-23T02:30:00Z"]
    }
  }
}
`, template, data.RandomInteger)
}
