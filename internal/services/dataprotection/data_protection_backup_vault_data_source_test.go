// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataProtectionBackupVaultDataSource struct{}

func TestAccDataProtectionBackupVaultDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_protection_backup_vault", "test")
	r := DataProtectionBackupVaultDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("datastore_type").HasValue("VaultStore"),
				check.That(data.ResourceName).Key("redundancy").HasValue("LocallyRedundant"),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
	})
}

func (r DataProtectionBackupVaultDataSource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupvaults.ParseBackupVaultID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupVaultClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r DataProtectionBackupVaultDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-bv-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}
data "azurerm_data_protection_backup_vault" "test" {
  name                = azurerm_data_protection_backup_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
