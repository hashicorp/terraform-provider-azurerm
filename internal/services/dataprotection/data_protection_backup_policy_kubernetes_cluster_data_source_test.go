// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataProtectionBackupPolicyKubernetesClusterDataSource struct{}

func TestAccDataProtectionBackupPolicyKubernetesClusterDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_protection_backup_policy_kubernetes_cluster", "test")
	r := DataProtectionBackupPolicyKubernetesClusterDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataProtectionBackupPolicyKubernetesClusterResource{}),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("vault_id").Exists(),
				check.That(data.ResourceName).Key("backup_repeating_time_intervals.#").HasValue("1"),
				check.That(data.ResourceName).Key("default_retention_rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("default_retention_rule.0.life_cycle.0.duration").HasValue("P7D"),
				check.That(data.ResourceName).Key("default_retention_rule.0.life_cycle.0.data_store_type").HasValue("OperationalStore"),
			),
		},
	})
}

func (r DataProtectionBackupPolicyKubernetesClusterDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_protection_backup_policy_kubernetes_cluster" "test" {
  name     = azurerm_data_protection_backup_policy_kubernetes_cluster.test.name
  vault_id = azurerm_data_protection_backup_vault.test.id
}
`, DataProtectionBackupPolicyKubernetesClusterResource{}.basic(data))
}
