// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource struct{}

func TestAccDataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_protection_backup_policy_postgresql_flexible_server", "test")
	r := DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataProtectionBackupPolicyPostgresqlFlexibleServerResource{}),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("vault_id").Exists(),
				check.That(data.ResourceName).Key("backup_repeating_time_intervals.#").HasValue("1"),
				check.That(data.ResourceName).Key("default_retention_rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("default_retention_rule.0.life_cycle.0.duration").HasValue("P4M"),
				check.That(data.ResourceName).Key("default_retention_rule.0.life_cycle.0.data_store_type").HasValue("VaultStore"),
			),
		},
	})
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_protection_backup_policy_postgresql_flexible_server" "test" {
  name     = azurerm_data_protection_backup_policy_postgresql_flexible_server.test.name
  vault_id = azurerm_data_protection_backup_vault.test.id
}
`, DataProtectionBackupPolicyPostgresqlFlexibleServerResource{}.basic(data))
}
