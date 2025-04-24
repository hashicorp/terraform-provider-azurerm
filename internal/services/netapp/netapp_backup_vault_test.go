// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backupvaults"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppBackupVaultResource struct{}

func (t NetAppBackupVaultResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := backupvaults.ParseBackupVaultID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.BackupVaultsClient.Get(ctx, *id)
	if err != nil {
		if resp.HttpResponse.StatusCode == http.StatusNotFound {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func TestAccNetAppBackupVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_backup_vault", "test")
	r := NetAppBackupVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppBackupVault_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_backup_vault", "test")
	r := NetAppBackupVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (r NetAppBackupVaultResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_backup_vault" "test" {
  name                = "acctest-NetAppBackupVault-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  account_name        = azurerm_netapp_account.test.name

  tags = {
    "testTag" = "testTagValue"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppBackupVaultResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_backup_vault" "test" {
  name                = "acctest-NetAppBackupVault-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  account_name        = azurerm_netapp_account.test.name

  tags = {
    "testTag" = "testTagValue",
    "FoO"     = "BaR"
  }
}
`, r.template(data), data.RandomInteger)
}

func (NetAppBackupVaultResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    netapp {
      delete_backups_on_backup_vault_destroy = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2023-08-17T08:01:00Z",
    "SkipASMAzSecPack" = "true"
  }
}


`, data.RandomInteger, data.Locations.Primary)
}
