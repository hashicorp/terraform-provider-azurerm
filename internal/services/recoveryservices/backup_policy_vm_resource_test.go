// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BackupProtectionPolicyVMResource struct{}

func TestAccBackupProtectionPolicyVM_policyTypeDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.policyTypeDefault(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("policy_type").HasValue("V1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_basicDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDaily(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_withInstantRestoreRetentionRangeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDailyWithInstantRestoreRetentionRange(data, 1, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicDailyWithInstantRestoreRetentionRange(data, 5, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDaily(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBackupProtectionPolicyVM_basicWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWeekly(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_completeDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeDaily(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_completeWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWeekly(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Weekly"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("backup.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("42"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("77"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDaily(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeDaily(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("42"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("77"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWeekly(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeWeekly(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Weekly"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("backup.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("42"),
				check.That(data.ResourceName).Key("retention_weekly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("77"),
				check.That(data.ResourceName).Key("retention_yearly.0.weekdays.#").HasValue("4"),
				check.That(data.ResourceName).Key("retention_yearly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateDailyToWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDaily(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWeekly(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateWeeklyToDaily(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWeekly(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicDaily(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateWeeklyToPartial(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWeekly(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeWeeklyPartial(data, "V1"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_basicHourlyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicHourly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_basicDailyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDaily(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_withInstantRestoreRetentionRangeUpdateV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDailyWithInstantRestoreRetentionRange(data, 1, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicDailyWithInstantRestoreRetentionRange(data, 30, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_basicWeeklyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWeekly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_completeHourlyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeHourly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_completeDailyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeDaily(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_completeWeeklyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWeekly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateHourlyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicHourly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeHourly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateDailyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDaily(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeDaily(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateWeeklyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWeekly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.completeWeekly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateHourlyAndWeeklyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicHourly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWeekly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicHourly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateDailyAndWeeklyV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDaily(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWeekly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicDaily(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateWeeklyToPartialV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWeekly(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeWeeklyPartial(data, "V2"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_withCustomRGName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomResourceGroup(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_basicDays(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.daysBasic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_completeDays(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.daysComplete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_tieringPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tieringPolicy(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackupProtectionPolicyVM_updateBackupTime(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.policyTypeDefault(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateBackupTime(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t BackupProtectionPolicyVMResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := protectionpolicies.ParseBackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.ProtectionPoliciesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service Protection Policy (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (BackupProtectionPolicyVMResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-backup-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) policyTypeDefault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, r.template(data), data.RandomInteger)
}

// nolint: unparam
func (r BackupProtectionPolicyVMResource) basicHourly(data acceptance.TestData, policyType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency     = "Hourly"
    time          = "23:00"
    hour_interval = 4
    hour_duration = 4
  }

  retention_daily {
    count = 10
  }

  policy_type = "%s"
}
`, r.template(data), data.RandomInteger, policyType)
}

func (r BackupProtectionPolicyVMResource) basicDaily(data acceptance.TestData, policyType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }

  policy_type = "%s"
}
`, r.template(data), data.RandomInteger, policyType)
}

func (r BackupProtectionPolicyVMResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "import" {
  name                = azurerm_backup_policy_vm.test.name
  resource_group_name = azurerm_backup_policy_vm.test.resource_group_name
  recovery_vault_name = azurerm_backup_policy_vm.test.recovery_vault_name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, r.basicDaily(data, "V1"))
}

func (r BackupProtectionPolicyVMResource) basicWeekly(data acceptance.TestData, policyType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday"]
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday"]
  }

  policy_type = "%s"
}
`, r.template(data), data.RandomInteger, policyType)
}

func (r BackupProtectionPolicyVMResource) completeHourly(data acceptance.TestData, policyType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  timezone            = "UTC"

  backup {
    frequency     = "Hourly"
    time          = "23:00"
    hour_interval = 12
    hour_duration = 24
  }

  retention_daily {
    count = 10
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }

  policy_type = "%s"
}
`, r.template(data), data.RandomInteger, policyType)
}

func (r BackupProtectionPolicyVMResource) completeDaily(data acceptance.TestData, policyType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  timezone            = "UTC"
  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }

  policy_type = "%s"
}
`, r.template(data), data.RandomInteger, policyType)
}

func (r BackupProtectionPolicyVMResource) completeWeekly(data acceptance.TestData, policyType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }

  policy_type = "%s"
}
`, r.template(data), data.RandomInteger, policyType)
}

func (r BackupProtectionPolicyVMResource) completeWeeklyPartial(data acceptance.TestData, policyType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday", "Friday"]
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday"]
    weeks    = ["Last"]
    months   = ["January"]
  }

  policy_type = "%s"
}
`, r.template(data), data.RandomInteger, policyType)
}

func (r BackupProtectionPolicyVMResource) basicDailyWithInstantRestoreRetentionRange(data acceptance.TestData, retentionRange int, policyType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                           = "acctest-BPVM-%d"
  resource_group_name            = azurerm_resource_group.test.name
  recovery_vault_name            = azurerm_recovery_services_vault.test.name
  instant_restore_retention_days = %d
  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 31
  }

  policy_type = "%s"
}
`, r.template(data), data.RandomInteger, retentionRange, policyType)
}

func (r BackupProtectionPolicyVMResource) withCustomResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  instant_restore_resource_group {
    prefix = "acctest"
    suffix = "suffix"
  }

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }
}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) daysBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }

  retention_monthly {
    count = 10
    days  = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  }

  retention_yearly {
    count  = 10
    months = ["January", "July"]
    days   = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  }

}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) daysComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "23:00"
  }

  retention_daily {
    count = 10
  }

  retention_monthly {
    count             = 10
    days              = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    include_last_days = true
  }

  retention_yearly {
    count             = 10
    months            = ["January", "July"]
    days              = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    include_last_days = true
  }

}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) tieringPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-bpvm-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Weekly"
    time      = "23:00"
    weekdays  = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_weekly {
    count    = 42
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
  }

  retention_monthly {
    count    = 7
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
    weeks    = ["First", "Last"]
  }

  retention_yearly {
    count    = 77
    weekdays = ["Sunday", "Wednesday", "Friday", "Saturday"]
    weeks    = ["First", "Last"]
    months   = ["January", "July"]
  }

  tiering_policy {
    archived_restore_point {
      duration      = 5
      duration_type = "Months"
      mode          = "TierAfter"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r BackupProtectionPolicyVMResource) updateBackupTime(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_backup_policy_vm" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name

  backup {
    frequency = "Daily"
    time      = "20:00"
  }

  retention_daily {
    count = 10
  }
}
`, r.template(data), data.RandomInteger)
}
