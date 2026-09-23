// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BackupProtectionPolicyFileShareDataSource struct{}

func TestAccDataSourceBackupPolicyFileShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("recovery_vault_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.time").HasValue("23:00"),
				check.That(data.ResourceName).Key("backup_tier").HasValue("snapshot"),
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("snapshot_retention_in_days").HasValue("0"),
				check.That(data.ResourceName).Key("timezone").HasValue("UTC"),
			),
		},
	})
}

func TestAccDataSourceBackupPolicyFileShare_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("retention_daily.0.count").HasValue("10"),
				check.That(data.ResourceName).Key("retention_weekly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_monthly.0.weeks.#").HasValue("2"),
				check.That(data.ResourceName).Key("retention_yearly.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("retention_yearly.0.months.#").HasValue("2"),
			),
		},
	})
}

func TestAccDataSourceBackupPolicyFileShare_vaultStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_backup_policy_file_share", "test")
	r := BackupProtectionPolicyFileShareDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.vaultStandard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("backup.0.frequency").HasValue("Hourly"),
				check.That(data.ResourceName).Key("backup.0.hourly.0.interval").HasValue("4"),
				check.That(data.ResourceName).Key("backup.0.hourly.0.start_time").HasValue("10:00"),
				check.That(data.ResourceName).Key("backup.0.hourly.0.window_duration").HasValue("12"),
				check.That(data.ResourceName).Key("backup_tier").HasValue("vault-standard"),
				check.That(data.ResourceName).Key("snapshot_retention_in_days").HasValue("9"),
			),
		},
	})
}

func (BackupProtectionPolicyFileShareDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_backup_policy_file_share" "test" {
  name                = azurerm_backup_policy_file_share.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, BackupProtectionPolicyFileShareResource{}.basicDaily(data))
}

func (BackupProtectionPolicyFileShareDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_backup_policy_file_share" "test" {
  name                = azurerm_backup_policy_file_share.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, BackupProtectionPolicyFileShareResource{}.completeDaily(data))
}

func (BackupProtectionPolicyFileShareDataSource) vaultStandard(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_backup_policy_file_share" "test" {
  name                = azurerm_backup_policy_file_share.test.name
  recovery_vault_name = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, BackupProtectionPolicyFileShareResource{}.vaultStandard(data))
}
