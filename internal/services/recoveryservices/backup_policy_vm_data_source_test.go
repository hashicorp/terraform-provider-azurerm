// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BackupProtectionPolicyVMDataSource struct{}

func TestAccDataSourceBackupPolicyVm_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("policy_type").HasValue("V1"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("timezone").HasValue("UTC"),
			),
		},
	})
}

func TestAccDataSourceBackupPolicyVm_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_backup_policy_vm", "test")
	r := BackupProtectionPolicyVMDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Weekly"),
				check.That(data.ResourceName).Key("consistency_type").HasValue("OnlyCrashConsistent"),
				check.That(data.ResourceName).Key("instant_restore_resource_group.0.prefix").HasValue("acctest"),
				check.That(data.ResourceName).Key("instant_restore_resource_group.0.suffix").HasValue("suffix"),
				check.That(data.ResourceName).Key("instant_restore_retention_days").HasValue("30"),
				check.That(data.ResourceName).Key("policy_type").HasValue("V2"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("42"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("77"),
				check.That(data.ResourceName).Key("tiering_policy.0.archived_restore_point.0.duration").HasValue("5"),
				check.That(data.ResourceName).Key("tiering_policy.0.archived_restore_point.0.duration_type").HasValue("Months"),
				check.That(data.ResourceName).Key("tiering_policy.0.archived_restore_point.0.mode").HasValue("TierAfter"),
				check.That(data.ResourceName).Key("timezone").HasValue("UTC"),
			),
		},
	})
}

func (BackupProtectionPolicyVMDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_backup_policy_vm" "test" {
  name                = azurerm_backup_policy_vm.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, BackupProtectionPolicyVMResource{}.basicDaily(data, "V1"))
}

func (BackupProtectionPolicyVMDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_backup_policy_vm" "test" {
  name                           = "acctest-%[2]d"
  resource_group_name            = azurerm_resource_group.test.name
  recovery_vault_name            = azurerm_recovery_services_vault.test.name
  timezone                       = "UTC"
  policy_type                    = "V2"
  consistency_type               = "OnlyCrashConsistent"
  instant_restore_retention_days = 30

  instant_restore_resource_group {
    prefix = "acctest"
    suffix = "suffix"
  }

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

data "azurerm_backup_policy_vm" "test" {
  name                = azurerm_backup_policy_vm.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, BackupProtectionPolicyVMResource{}.template(data), data.RandomInteger)
}
