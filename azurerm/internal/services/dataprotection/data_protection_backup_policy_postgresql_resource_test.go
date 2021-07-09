package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataProtectionBackupPolicyPostgreSQLResource struct{}

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

func TestAccDataProtectionBackupPolicyPostgreSQL_update(t *testing.T) {
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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DataProtectionBackupPolicyPostgreSQLResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r DataProtectionBackupPolicyPostgreSQLResource) basic(data acceptance.TestData) string {
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

func (r DataProtectionBackupPolicyPostgreSQLResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
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

func (r DataProtectionBackupPolicyPostgreSQLResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_postgresql" "test" {
  name                = "acctest-dbp-%d"
  resource_group_name = azurerm_resource_group.test.name
  vault_name          = azurerm_data_protection_backup_vault.test.name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  default_retention_duration      = "P4M"
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
