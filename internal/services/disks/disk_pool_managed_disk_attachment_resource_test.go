// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package disks_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DisksPoolManagedDiskAttachmentResource struct{}

func TestAccDiskPoolDiskAttachment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_managed_disk_attachment", "test")
	a := DisksPoolManagedDiskAttachmentResource{}
	data.ResourceTest(t, a, []acceptance.TestStep{
		{
			Config: a.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(a),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskPoolDiskAttachment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_managed_disk_attachment", "test")
	a := DisksPoolManagedDiskAttachmentResource{}
	data.ResourceTest(t, a, []acceptance.TestStep{
		{
			Config: a.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(a),
		},
		{
			Config:      a.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_disk_pool_managed_disk_attachment"),
		},
	})
}

func TestAccDiskPoolDiskAttachment_multipleDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_managed_disk_attachment", "test")
	a := DisksPoolManagedDiskAttachmentResource{}
	secondResourceName := "azurerm_disk_pool_managed_disk_attachment.second"
	data.ResourceTest(t, a, []acceptance.TestStep{
		{
			Config: a.multipleDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(a),
				check.That(secondResourceName).ExistsInAzure(a),
			),
		},
		data.ImportStep(),
		{
			Config: a.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(a),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskPoolDiskAttachment_destroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_managed_disk_attachment", "test")
	a := DisksPoolManagedDiskAttachmentResource{}
	data.ResourceTest(t, a, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       a.basic,
			TestResource: a,
		}),
	})
}

func (a DisksPoolManagedDiskAttachmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DiskPoolManagedDiskAttachmentID(state.ID)
	if err != nil {
		return nil, err
	}
	poolId := id.DiskPoolId
	diskId := id.ManagedDiskId
	client := clients.Disks.DiskPoolsClient
	resp, err := client.Get(ctx, poolId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	targetDiskId := diskId.ID()
	if resp.Model == nil || resp.Model.Properties.Disks == nil {
		return utils.Bool(false), nil
	}
	for _, disk := range *resp.Model.Properties.Disks {
		if disk.Id == targetDiskId {
			return utils.Bool(true), nil
		}
	}
	return utils.Bool(false), nil
}

func (a DisksPoolManagedDiskAttachmentResource) Destroy(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Minute)
	defer cancel()
	id, err := parse.DiskPoolManagedDiskAttachmentID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Disks.DiskPoolsClient
	pool, err := client.Get(ctx, id.DiskPoolId)
	if err != nil {
		return nil, err
	}
	if pool.Model == nil || pool.Model.Properties.Disks == nil {
		return nil, err
	}
	attachedDisks := pool.Model.Properties.Disks
	remainingDisks := make([]diskpools.Disk, 0)
	for _, attachedDisk := range *attachedDisks {
		if attachedDisk.Id != id.ManagedDiskId.ID() {
			remainingDisks = append(remainingDisks, attachedDisk)
		}
	}

	err = client.UpdateThenPoll(ctx, id.DiskPoolId, diskpools.DiskPoolUpdate{
		Properties: diskpools.DiskPoolUpdateProperties{
			Disks: &remainingDisks,
		},
	})
	if err != nil {
		return nil, err
	}
	return utils.Bool(true), nil
}

func (a DisksPoolManagedDiskAttachmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_disk_pool_managed_disk_attachment" "test" {
  depends_on      = [azurerm_role_assignment.test]
  disk_pool_id    = azurerm_disk_pool.test.id
  managed_disk_id = azurerm_managed_disk.test.id
}
`, a.template(data))
}

func (a DisksPoolManagedDiskAttachmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-diskpool-%[2]d"
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
  name                 = "acctest-diskpool-%[2]d"
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
  name                = "acctest-diskpool-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zones               = ["1"]
  sku_name            = "Basic_B1"
  subnet_id           = azurerm_subnet.test.id
  tags = {
    "env" = "qa"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (a DisksPoolManagedDiskAttachmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_pool_managed_disk_attachment" "import" {
  disk_pool_id    = azurerm_disk_pool.test.id
  managed_disk_id = azurerm_managed_disk.test.id
}
`, a.basic(data))
}

func (a DisksPoolManagedDiskAttachmentResource) multipleDisks(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "second" {
  name                 = "acctest-diskspool-%d-2"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  create_option        = "Empty"
  storage_account_type = "Premium_LRS"
  disk_size_gb         = 4
  max_shares           = 2
  zone                 = "1"
}

resource "azurerm_role_assignment" "second" {
  count                = length(local.roles)
  principal_id         = data.azuread_service_principal.test.id
  role_definition_name = local.roles[count.index]
  scope                = azurerm_managed_disk.second.id
}

resource "azurerm_disk_pool_managed_disk_attachment" "second" {
  depends_on      = [azurerm_role_assignment.second]
  disk_pool_id    = azurerm_disk_pool.test.id
  managed_disk_id = azurerm_managed_disk.second.id
}
`, a.basic(data), data.RandomInteger)
}
