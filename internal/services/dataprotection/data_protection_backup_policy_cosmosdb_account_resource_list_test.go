// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccDataProtectionBackupPolicyCosmosdbAccount_listByDataProtectionBackupVaultID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_cosmosdb_account", "test")
	r := DataProtectionBackupPolicyCosmosdbAccountResource{}
	listResourceAddress := "azurerm_data_protection_backup_policy_cosmosdb_account.list"

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQueryByDataProtectionBackupVaultID(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_cosmosdb_account" "test" {
  count = 3

  name                            = "acctest-dbp-cosmos-${count.index}-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  default_retention_duration      = "P10Y"
  full_backup_schedule            = "R/2026-02-08T10:00:00+00:00/P1W"
}

resource "azurerm_data_protection_backup_policy_disk" "other" {
  name                            = "acctest-dbp-disk-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-19T06:33:16+00:00/PT4H"]
  default_retention_duration      = "P7D"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) basicQueryByDataProtectionBackupVaultID() string {
	return `
list "azurerm_data_protection_backup_policy_cosmosdb_account" "list" {
  provider = azurerm

  config {
    data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  }
}
`
}
