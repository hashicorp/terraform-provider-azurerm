// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupInstanceElasticSanVolumeGroupResource struct{}

func TestAccDataProtectionBackupInstanceElasticSanVolumeGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_elastic_san_volume_group", "test")
	r := DataProtectionBackupInstanceElasticSanVolumeGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstanceElasticSanVolumeGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_elastic_san_volume_group", "test")
	r := DataProtectionBackupInstanceElasticSanVolumeGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDataProtectionBackupInstanceElasticSanVolumeGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_elastic_san_volume_group", "test")
	r := DataProtectionBackupInstanceElasticSanVolumeGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume_name").HasValue("vol2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupInstanceElasticSanVolumeGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_elastic_san_volume_group", "test")
	r := DataProtectionBackupInstanceElasticSanVolumeGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume_name").HasValue("vol1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume_name").HasValue("vol2"),
			),
		},
		data.ImportStep(),
	})
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstanceresources.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient20260601.BackupInstancesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_san" "test" {
  name                = "acctestes%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  base_size_in_tib    = 1

  sku {
    name = "Premium_LRS"
  }
}

resource "azurerm_elastic_san_volume_group" "test" {
  name           = "acctestesvg%d"
  elastic_san_id = azurerm_elastic_san.test.id
}

resource "azurerm_elastic_san_volume" "first" {
  name            = "vol1"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  size_in_gib     = 1
}

resource "azurerm_elastic_san_volume" "second" {
  name            = "vol2"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  size_in_gib     = 1
}

resource "azurerm_role_assignment" "snapshot_exporter" {
  scope                = azurerm_elastic_san.test.id
  role_definition_name = "Elastic SAN Snapshot Exporter"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}

resource "azurerm_role_assignment" "snapshot_contributor" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Disk Snapshot Contributor"
  principal_id         = azurerm_data_protection_backup_vault.test.identity[0].principal_id
}
`, DataProtectionBackupPolicyElasticSanVolumeGroupResource{}.basic(data), data.RandomIntOfLength(8), data.RandomInteger)
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_elastic_san_volume_group" "test" {
  name                         = "acctest-dbi-esan-%d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  backup_policy_id             = azurerm_data_protection_backup_policy_elastic_san_volume_group.test.id
  elastic_san_volume_group_id  = azurerm_elastic_san_volume_group.test.id
  snapshot_resource_group_name = azurerm_resource_group.test.name
  volume_name                  = azurerm_elastic_san_volume.first.name

  depends_on = [
    azurerm_role_assignment.snapshot_exporter,
    azurerm_role_assignment.snapshot_contributor,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_elastic_san_volume_group" "import" {
  name                         = azurerm_data_protection_backup_instance_elastic_san_volume_group.test.name
  location                     = azurerm_data_protection_backup_instance_elastic_san_volume_group.test.location
  vault_id                     = azurerm_data_protection_backup_instance_elastic_san_volume_group.test.vault_id
  backup_policy_id             = azurerm_data_protection_backup_instance_elastic_san_volume_group.test.backup_policy_id
  elastic_san_volume_group_id  = azurerm_data_protection_backup_instance_elastic_san_volume_group.test.elastic_san_volume_group_id
  snapshot_resource_group_name = azurerm_data_protection_backup_instance_elastic_san_volume_group.test.snapshot_resource_group_name
  volume_name                  = azurerm_data_protection_backup_instance_elastic_san_volume_group.test.volume_name
}
`, r.basic(data))
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_elastic_san_volume_group" "test" {
  name                         = "acctest-dbi-esan-%d"
  location                     = azurerm_resource_group.test.location
  vault_id                     = azurerm_data_protection_backup_vault.test.id
  backup_policy_id             = azurerm_data_protection_backup_policy_elastic_san_volume_group.test.id
  elastic_san_volume_group_id  = azurerm_elastic_san_volume_group.test.id
  snapshot_resource_group_name = azurerm_resource_group.test.name
  volume_name                  = azurerm_elastic_san_volume.second.name

  depends_on = [
    azurerm_role_assignment.snapshot_exporter,
    azurerm_role_assignment.snapshot_contributor,
  ]
}
`, r.template(data), data.RandomInteger)
}
