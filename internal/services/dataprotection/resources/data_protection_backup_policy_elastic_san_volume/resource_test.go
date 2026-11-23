// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupPolicyElasticSanVolumeResource struct{}

func TestAccDataProtectionBackupPolicyElasticSanVolume_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_elastic_san_volume", "test")
	r := DataProtectionBackupPolicyElasticSanVolumeResource{}
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

func TestAccDataProtectionBackupPolicyElasticSanVolume_basicWeekly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_elastic_san_volume", "test")
	r := DataProtectionBackupPolicyElasticSanVolumeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWeekly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataProtectionBackupPolicyElasticSanVolume_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_elastic_san_volume", "test")
	r := DataProtectionBackupPolicyElasticSanVolumeResource{}
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

func TestAccDataProtectionBackupPolicyElasticSanVolume_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_elastic_san_volume", "test")
	r := DataProtectionBackupPolicyElasticSanVolumeResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DataProtectionBackupPolicyElasticSanVolumeResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupPolicyClient20260601.BackupPoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupPolicyElasticSanVolumeResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dp-esan-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dbv-esan-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DataProtectionBackupPolicyElasticSanVolumeResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_elastic_san_volume" "test" {
  name                            = "acctest-dbp-esan-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  daily_backup_schedule           = "R/2026-02-08T13:00:00+00:00/P1D"
  default_retention_duration      = "P7D"
  daily_retention_duration        = "P14D"
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupPolicyElasticSanVolumeResource) basicWeekly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_elastic_san_volume" "test" {
  name                            = "acctest-dbp-esan-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  weekly_backup_schedule          = "R/2026-02-08T13:00:00+00:00/P1W"
  default_retention_duration      = "P7D"
  monthly_retention_duration      = "P3M"
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupPolicyElasticSanVolumeResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_elastic_san_volume" "import" {
  name                            = azurerm_data_protection_backup_policy_elastic_san_volume.test.name
  data_protection_backup_vault_id = azurerm_data_protection_backup_policy_elastic_san_volume.test.data_protection_backup_vault_id
  daily_backup_schedule           = azurerm_data_protection_backup_policy_elastic_san_volume.test.daily_backup_schedule
  default_retention_duration      = azurerm_data_protection_backup_policy_elastic_san_volume.test.default_retention_duration
  daily_retention_duration        = azurerm_data_protection_backup_policy_elastic_san_volume.test.daily_retention_duration
}
`, r.basic(data))
}

func (r DataProtectionBackupPolicyElasticSanVolumeResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_elastic_san_volume" "test" {
  name                            = "acctest-dbp-esan-%d"
  data_protection_backup_vault_id = azurerm_data_protection_backup_vault.test.id
  daily_backup_schedule           = "R/2026-02-08T13:00:00+00:00/P1D"
  default_retention_duration      = "P7D"
  daily_retention_duration        = "P14D"
  weekly_retention_duration       = "P28D"
  monthly_retention_duration      = "P3M"
}
`, r.template(data), data.RandomInteger)
}
