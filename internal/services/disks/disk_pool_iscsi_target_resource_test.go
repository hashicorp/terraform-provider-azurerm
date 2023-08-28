// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package disks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DisksPoolIscsiTargetResource struct{}

func TestDiskPoolIscsiTarget_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_iscsi_target", "test")
	i := DisksPoolIscsiTargetResource{}

	data.ResourceTest(t, i, []acceptance.TestStep{
		{
			Config: i.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(i),
			),
		},
		data.ImportStep(),
	})
}

func TestDiskPoolIscsiTarget_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_iscsi_target", "test")
	i := DisksPoolIscsiTargetResource{}
	data.ResourceTest(t, i, []acceptance.TestStep{
		{
			Config: i.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(i),
			),
		},
		{
			Config:      i.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_disk_pool_iscsi_target"),
		},
	})
}

func (i DisksPoolIscsiTargetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := iscsitargets.ParseIscsiTargetID(state.ID)
	if err != nil {
		return nil, err
	}
	client := clients.Disks.DisksPoolIscsiTargetClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (i DisksPoolIscsiTargetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-diskspool-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/24"]
  delegation {
    name = "diskspool"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/read"]
      name    = "Microsoft.StoragePool/diskPools"
    }
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-diskspool-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  create_option        = "Empty"
  storage_account_type = "Premium_LRS"
  disk_size_gb         = 4
  max_shares           = 2
  zone                 = "1"
}

data "azuread_service_principal" "test" {
  display_name = "StoragePool Resource Provider"
}

locals {
  roles = ["Disk Pool Operator", "Virtual Machine Contributor"]
}

resource "azurerm_role_assignment" "test" {
  count                = length(local.roles)
  principal_id         = data.azuread_service_principal.test.id
  role_definition_name = local.roles[count.index]
  scope                = azurerm_managed_disk.test.id
}

resource "azurerm_disk_pool" "test" {
  name                = "acctest-diskspool-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zones               = ["1"]
  sku_name            = "Basic_B1"
  subnet_id           = azurerm_subnet.test.id
  tags = {
    "env" = "qa"
  }
}

resource "azurerm_disk_pool_managed_disk_attachment" "test" {
  depends_on      = [azurerm_role_assignment.test]
  disk_pool_id    = azurerm_disk_pool.test.id
  managed_disk_id = azurerm_managed_disk.test.id
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (i DisksPoolIscsiTargetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_pool_iscsi_target" "test" {
  depends_on    = [azurerm_disk_pool_managed_disk_attachment.test]
  name          = "acctest-diskpool-%s"
  acl_mode      = "Dynamic"
  disks_pool_id = azurerm_disk_pool.test.id
  target_iqn    = "iqn.2021-11.com.microsoft:test"
}
`, i.template(data), data.RandomString)
}

func (i DisksPoolIscsiTargetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_pool_iscsi_target" "import" {
  name          = "acctest-diskpool-%s"
  acl_mode      = "Dynamic"
  disks_pool_id = azurerm_disk_pool.test.id
  target_iqn    = "iqn.2021-11.com.microsoft:test"
}
`, i.basic(data), data.RandomString)
}
