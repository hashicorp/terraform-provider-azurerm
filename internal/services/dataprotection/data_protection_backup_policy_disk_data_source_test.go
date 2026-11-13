// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataProtectionBackupPolicyDiskDataSource struct{}

func TestAccDataProtectionBackupPolicyDiskDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_protection_backup_policy_disk", "test")
	r := DataProtectionBackupPolicyDiskDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataProtectionBackupPolicyDiskResource{}),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("vault_id").Exists(),
				check.That(data.ResourceName).Key("default_retention_duration").HasValue("P7D"),
				check.That(data.ResourceName).Key("backup_repeating_time_intervals.#").HasValue("1"),
			),
		},
	})
}

func (r DataProtectionBackupPolicyDiskDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_protection_backup_policy_disk" "test" {
  name     = azurerm_data_protection_backup_policy_disk.test.name
  vault_id = azurerm_data_protection_backup_vault.test.id
}
`, DataProtectionBackupPolicyDiskResource{}.basic(data))
}
