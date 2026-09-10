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

func TestAccDataProtectionBackupInstanceElasticSanVolumeGroup_listByVaultID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_elastic_san_volume_group", "test")
	r := DataProtectionBackupInstanceElasticSanVolumeGroupResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.listConfig(data),
			},
			{
				Query:  true,
				Config: r.listQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_data_protection_backup_instance_elastic_san_volume_group.list", 2),
				},
			},
		},
	})
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) listConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_elastic_san_volume_group" "list" {
  count = 2

  name                         = "acctest-dbi-esan-${count.index}-%d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  backup_policy_id             = azurerm_data_protection_backup_policy_elastic_san_volume_group.test.id
  elastic_san_volume_group_id  = azurerm_elastic_san_volume_group.test.id
  snapshot_resource_group_name = azurerm_resource_group.test.name
  volume_name                  = count.index == 0 ? azurerm_elastic_san_volume.first.name : azurerm_elastic_san_volume.second.name

  depends_on = [
    azurerm_role_assignment.snapshot_exporter,
    azurerm_role_assignment.snapshot_contributor,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (DataProtectionBackupInstanceElasticSanVolumeGroupResource) listQuery() string {
	return `
list "azurerm_data_protection_backup_instance_elastic_san_volume_group" "list" {
  provider = azurerm
  config {
    vault_id = azurerm_data_protection_backup_vault.test.id
  }
}
`
}
