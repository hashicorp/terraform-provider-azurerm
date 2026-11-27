// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume_test

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

func TestAccDataProtectionBackupPolicyElasticSanVolume_listByDataProtectionBackupVaultID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_elastic_san_volume", "test")
	r := DataProtectionBackupPolicyElasticSanVolumeResource{}
	listResourceAddress := "azurerm_data_protection_backup_policy_elastic_san_volume.list"

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

func (r DataProtectionBackupPolicyElasticSanVolumeResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_elastic_san_volume" "test" {
  count = 3

  name                            = "acctest-dbp-esan-${count.index}-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  daily_backup_schedule           = "R/2026-02-08T13:00:00+00:00/P1D"
  default_retention_duration      = "P7D"
}

resource "azurerm_data_protection_backup_policy_disk" "other" {
  name                            = "acctest-dbp-disk-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-19T06:33:16+00:00/PT4H"]
  default_retention_duration      = "P7D"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r DataProtectionBackupPolicyElasticSanVolumeResource) basicQueryByDataProtectionBackupVaultID() string {
	return `
list "azurerm_data_protection_backup_policy_elastic_san_volume" "list" {
  provider = azurerm

  config {
    data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  }
}
`
}
