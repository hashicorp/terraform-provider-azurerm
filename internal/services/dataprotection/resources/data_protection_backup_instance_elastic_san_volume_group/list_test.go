// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_elastic_san_volume_group_test

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

func TestAccDataProtectionBackupInstanceElasticSanVolumeGroup_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_elastic_san_volume_group", "test")
	r := DataProtectionBackupInstanceElasticSanVolumeGroupResource{}
	listResourceAddress := "azurerm_data_protection_backup_instance_elastic_san_volume_group.list"

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
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
				},
			},
		},
	})
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_san_volume_group" "test2" {
  name           = "acctestesvg2%d"
  elastic_san_id = azurerm_elastic_san.test.id
}

resource "azurerm_elastic_san_volume" "test2" {
  name            = "vol1"
  volume_group_id = azurerm_elastic_san_volume_group.test2.id
  size_in_gib     = 1
}

resource "azurerm_data_protection_backup_instance_elastic_san_volume_group" "test2" {
  name                         = "acctest-dbi-esan2-%d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  backup_policy_id             = azurerm_data_protection_backup_policy_elastic_san_volume.test.id
  elastic_san_volume_group_id  = azurerm_elastic_san_volume_group.test2.id
  snapshot_resource_group_name = azurerm_resource_group.test.name
  volume_name                  = azurerm_elastic_san_volume.test2.name

  depends_on = [
    azurerm_role_assignment.snapshot_exporter,
    azurerm_role_assignment.snapshot_contributor,
  ]
}
`, r.basic(data), data.RandomInteger, data.RandomInteger)
}

func (DataProtectionBackupInstanceElasticSanVolumeGroupResource) basicListQuery() string {
	return `
list "azurerm_data_protection_backup_instance_elastic_san_volume_group" "list" {
  provider = azurerm

  config {
    vault_id = azurerm_data_protection_backup_vault.test.id
  }
}
`
}
