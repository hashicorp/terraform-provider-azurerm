// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package disks_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/iscsitargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DisksPoolIscsiTargetLunResource struct{}

func TestAccDiskPoolIscsiTargetLun_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_iscsi_target_lun", "test")
	l := DisksPoolIscsiTargetLunResource{}
	data.ResourceTest(t, l, []acceptance.TestStep{
		{
			Config: l.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(l),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskPoolIscsiTargetLun_multipleLuns(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_iscsi_target_lun", "test0")
	l := DisksPoolIscsiTargetLunResource{}
	data.ResourceTest(t, l, []acceptance.TestStep{
		{
			Config: l.multipleLuns(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_disk_pool_iscsi_target_lun.test0").ExistsInAzure(l),
				check.That("azurerm_disk_pool_iscsi_target_lun.test1").ExistsInAzure(l),
			),
		},
		data.ImportStep(),
		{
			Config: l.multipleLuns(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_disk_pool_iscsi_target_lun.test0").ExistsInAzure(l),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskPoolIscsiTargetLun_updateName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_iscsi_target_lun", "test")
	l := DisksPoolIscsiTargetLunResource{}
	data.ResourceTest(t, l, []acceptance.TestStep{
		{
			Config: l.basic(data),
		},
		data.ImportStep(),
		{
			Config: l.updateName(data),
		},
		data.ImportStep(),
	})
}

func TestAccDiskPoolIscsiTargetLun_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_iscsi_target_lun", "test")
	l := DisksPoolIscsiTargetLunResource{}
	data.ResourceTest(t, l, []acceptance.TestStep{
		{
			Config: l.basic(data),
		},
		{
			Config:      l.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_disk_pool_iscsi_target_lun"),
		},
	})
}

func TestAccDiskPoolIscsiTargetLun_destroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_pool_iscsi_target_lun", "test")
	l := DisksPoolIscsiTargetLunResource{}
	data.ResourceTest(t, l, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       l.basic,
			TestResource: l,
		}),
	})
}

func (r DisksPoolIscsiTargetLunResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IscsiTargetLunID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Disks.DisksPoolIscsiTargetClient
	resp, err := client.Get(ctx, id.IscsiTargetId)

	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", id.IscsiTargetId, err)
	}
	if resp.Model == nil {
		return nil, fmt.Errorf("malformed Iscsi Target response %q : %+v", id.IscsiTargetId, resp)
	}
	var luns []iscsitargets.IscsiLun
	if resp.Model.Properties.Luns != nil {
		luns = *resp.Model.Properties.Luns
	}
	for _, lun := range luns {
		if lun.ManagedDiskAzureResourceId == id.ManagedDiskId.ID() {
			return utils.Bool(true), nil
		}
	}
	return utils.Bool(false), nil
}

func (r DisksPoolIscsiTargetLunResource) Destroy(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Minute)
	defer cancel()
	id, err := parse.IscsiTargetLunID(state.ID)
	if err != nil {
		return nil, err
	}

	iscsiTargetId := id.IscsiTargetId

	locks.ByID(iscsiTargetId.ID())
	defer locks.UnlockByID(iscsiTargetId.ID())

	client := clients.Disks.DisksPoolIscsiTargetClient
	resp, err := client.Get(ctx, iscsiTargetId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(true), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", iscsiTargetId, err)
	}
	if resp.Model == nil {
		return nil, fmt.Errorf("malformed Iscsi Target response %q : %+v", iscsiTargetId.ID(), resp)
	}
	if resp.Model.Properties.Luns == nil {
		return utils.Bool(true), nil
	}
	luns := make([]iscsitargets.IscsiLun, 0)
	for _, lun := range *resp.Model.Properties.Luns {
		if lun.ManagedDiskAzureResourceId != id.ManagedDiskId.ID() {
			luns = append(luns, lun)
		}
	}
	sort.Slice(luns, func(i, j int) bool {
		return luns[i].ManagedDiskAzureResourceId < luns[j].ManagedDiskAzureResourceId
	})
	patch := iscsitargets.IscsiTargetUpdate{
		Properties: iscsitargets.IscsiTargetUpdateProperties{
			Luns: &luns,
		},
	}

	m := disks.DiskPoolIscsiTargetLunModel{}

	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("could not retrieve context deadline")
	}
	err = m.RetryError(time.Until(deadline), "waiting for delete DisksPool iscsi target", id.ID(), func() error {
		return client.UpdateThenPoll(ctx, iscsiTargetId, patch)
	})
	if err != nil {
		return nil, err
	}
	return utils.Bool(true), nil
}

func (r DisksPoolIscsiTargetLunResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_pool_iscsi_target_lun" "test" {
  iscsi_target_id                      = azurerm_disk_pool_iscsi_target.test.id
  disk_pool_managed_disk_attachment_id = azurerm_disk_pool_managed_disk_attachment.test[0].id
  name                                 = "test-0"
}
`, r.template(data, 1))
}

func (r DisksPoolIscsiTargetLunResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_pool_iscsi_target_lun" "import" {
  iscsi_target_id                      = azurerm_disk_pool_iscsi_target.test.id
  disk_pool_managed_disk_attachment_id = azurerm_disk_pool_managed_disk_attachment.test[0].id
  name                                 = "test-0"
}
`, r.basic(data))
}

func (r DisksPoolIscsiTargetLunResource) updateName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_disk_pool_iscsi_target_lun" "test" {
  iscsi_target_id                      = azurerm_disk_pool_iscsi_target.test.id
  disk_pool_managed_disk_attachment_id = azurerm_disk_pool_managed_disk_attachment.test[0].id
  name                                 = "updated-test-0"
}
`, r.template(data, 1))
}

func (r DisksPoolIscsiTargetLunResource) multipleLuns(data acceptance.TestData, diskCount int) string {
	tfCode := r.template(data, diskCount)
	for i := 0; i < diskCount; i++ {
		tfCode = fmt.Sprintf(`
%[1]s

resource "azurerm_disk_pool_iscsi_target_lun" "test%[2]d" {
  iscsi_target_id                      = azurerm_disk_pool_iscsi_target.test.id
  disk_pool_managed_disk_attachment_id = azurerm_disk_pool_managed_disk_attachment.test[%[2]d].id
  name                                 = "test-%[2]d"
}
`, tfCode, i)
	}
	return tfCode
}

func (r DisksPoolIscsiTargetLunResource) template(data acceptance.TestData, diskCount int) string {
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

locals {
  disk_count = %[4]d
}

resource "azurerm_managed_disk" "test" {
  count                = local.disk_count
  name                 = "acctest-diskspool-%[2]d${count.index}"
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
// DO NOT attempt to use mod operator to make two assignments into one because we must use "%%" to escape percent sign and that will break terrafmt
resource "azurerm_role_assignment" "disk_pool_operator" {
  count                = local.disk_count
  principal_id         = data.azuread_service_principal.test.id
  role_definition_name = "Disk Pool Operator"
  scope                = azurerm_managed_disk.test[count.index].id
}

resource "azurerm_role_assignment" "vm_contributor" {
  count                = local.disk_count
  principal_id         = data.azuread_service_principal.test.id
  role_definition_name = "Virtual Machine Contributor"
  scope                = azurerm_managed_disk.test[count.index].id
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
  count           = local.disk_count
  depends_on      = [azurerm_role_assignment.disk_pool_operator, azurerm_role_assignment.vm_contributor]
  disk_pool_id    = azurerm_disk_pool.test.id
  managed_disk_id = azurerm_managed_disk.test[count.index].id
}

resource "azurerm_disk_pool_iscsi_target" "test" {
  depends_on    = [azurerm_disk_pool_managed_disk_attachment.test]
  name          = "acctest-diskpool-%[3]s"
  acl_mode      = "Dynamic"
  disks_pool_id = azurerm_disk_pool.test.id
  target_iqn    = "iqn.2021-11.com.microsoft:test"
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString, diskCount)
}
