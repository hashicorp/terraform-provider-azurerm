// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppSnapshotResource struct{}

func TestAccNetAppSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")
	r := NetAppSnapshotResource{}

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

func TestAccNetAppSnapshot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")
	r := NetAppSnapshotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_netapp_snapshot"),
		},
	})
}

func TestAccNetAppSnapshot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_snapshot", "test")
	r := NetAppSnapshotResource{}

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

func (t NetAppSnapshotResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := snapshots.ParseSnapshotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.SnapshotClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r NetAppSnapshotResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_name         = azurerm_netapp_volume.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppSnapshotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "import" {
  name                = azurerm_netapp_snapshot.test.name
  location            = azurerm_netapp_snapshot.test.location
  resource_group_name = azurerm_netapp_snapshot.test.resource_group_name
  account_name        = azurerm_netapp_snapshot.test.account_name
  pool_name           = azurerm_netapp_snapshot.test.pool_name
  volume_name         = azurerm_netapp_snapshot.test.volume_name
}
`, r.basic(data))
}

func (r NetAppSnapshotResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-NetAppSnapshot-%d"
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_name         = azurerm_netapp_volume.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (NetAppSnapshotResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    netapp {
      prevent_volume_destruction             = false
      delete_backups_on_backup_vault_destroy = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "netapp"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Premium"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Premium"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
