// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetvms"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

// TestAccVirtualMachineScaleSetManagedDisk_attachedToScaleSetInstanceResize proves the core behaviour of this
// resource: a Managed Disk that is attached to a Virtual Machine Scale Set instance can be updated (here changing
// the SKU, which requires the disk to be offline) by hot-detaching it from the instance, updating it, and
// re-attaching it - without deallocating the instance. This is the scenario that fails with `azurerm_managed_disk`.
func TestAccVirtualMachineScaleSetManagedDisk_attachedToScaleSetInstanceResize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_scale_set_managed_disk", "test")
	r := VirtualMachineScaleSetManagedDiskResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.attached(data, "Standard_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// attach the disk to the first scale set instance out-of-band so the disk's `managedBy` becomes the instance
				data.CheckWithClientForResource(func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
					diskId, err := commonids.ParseManagedDiskID(state.Attributes["id"])
					if err != nil {
						return err
					}

					instanceId := virtualmachinescalesetvms.NewVirtualMachineScaleSetVirtualMachineID(diskId.SubscriptionId, diskId.ResourceGroupName, fmt.Sprintf("acctestvmss-%d", data.RandomInteger), "0")
					input := virtualmachinescalesetvms.AttachDetachDataDisksRequest{
						DataDisksToAttach: &[]virtualmachinescalesetvms.DataDisksToAttach{
							{
								DiskId: diskId.ID(),
								Lun:    pointer.To(int64(0)),
							},
						},
					}

					ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
					defer cancel()
					if err := clients.Compute.VirtualMachineScaleSetVMsClient.AttachDetachDataDisksThenPoll(ctx2, instanceId, input); err != nil {
						return fmt.Errorf("attaching %s to %s: %+v", diskId, instanceId, err)
					}

					return nil
				}, "azurerm_virtual_machine_scale_set_managed_disk.test"),
			),
		},
		{
			Config: r.attached(data, "StandardSSD_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("StandardSSD_LRS"),
			),
		},
	})
}

func (r VirtualMachineScaleSetManagedDiskResource) attached(data acceptance.TestData, skuName string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_linux_virtual_machine_scale_set" "test" {
  name                            = "acctestvmss-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  sku                             = "Standard_F2"
  instances                       = 1
  admin_username                  = "adminuser"
  admin_password                  = "P@ssword1234!"
  disable_password_authentication = false
  overprovision                   = false

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "internal"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}

resource "azurerm_virtual_machine_scale_set_managed_disk" "test" {
  name                 = "acctestd-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  storage_account_type = "%[3]s"

  creation {
    option = "Empty"
  }

  disk_size_gb = 1
}
`, r.template(data), data.RandomInteger, skuName)
}
