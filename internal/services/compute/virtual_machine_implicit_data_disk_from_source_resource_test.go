// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualMachineImplicitDataDiskFromSourceResource struct{}

func TestAccVirtualMachineImplicitDataDiskFromSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}

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

func TestAccVirtualMachineImplicitDataDiskFromSource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_machine_implicit_data_disk_from_source"),
		},
	})
}

func TestAccVirtualMachineImplicitDataDiskFromSource_destroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccVirtualMachineImplicitDataDiskFromSource_multipleDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "first")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}

	secondResourceName := "azurerm_virtual_machine_implicit_data_disk_from_source.second"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			ResourceName:      secondResourceName,
			ImportState:       true,
			ImportStateVerify: true,
		},
	})
}

func TestAccVirtualMachineImplicitDataDiskFromSource_updatingCaching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.readOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.readWrite(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachineImplicitDataDiskFromSource_updatingWriteAccelerator(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.writeAccelerator(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.writeAccelerator(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.writeAccelerator(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachineImplicitDataDiskFromSource_managedServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedServiceIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineImplicitDataDiskFromSource_virtualMachineExtension(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualMachineExtensionPrep(data),
		},
		{
			Config: r.virtualMachineExtensionComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachineImplicitDataDiskFromSource_virtualMachineApplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualMachineApplicationPrep(data),
		},
		{
			Config: r.virtualMachineApplicationComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualMachineImplicitDataDiskFromSource_detachImplicitDataDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_implicit_data_disk_from_source", "test")
	r := VirtualMachineImplicitDataDiskFromSourceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
		},
		{
			Config: r.detachImplicitDataDisk(data),
			Check: acceptance.ComposeTestCheckFunc(
				func(state *terraform.State) error {
					client, err := testclient.Build()
					if err != nil {
						return fmt.Errorf("building client: %+v", err)
					}

					ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(5*time.Minute))
					defer cancel()

					diskClient := client.Compute.DisksClient

					id, err := commonids.ParseManagedDiskID(fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Compute/disks/acctestVMIDD-%d", os.Getenv("ARM_SUBSCRIPTION_ID"), data.RandomInteger, data.RandomInteger))
					if err != nil {
						return err
					}

					_, err = diskClient.Get(ctx, *id)
					if err != nil {
						return fmt.Errorf("retrieving %s: %+v", id, err)
					}

					return nil
				},
			),
		},
	})
}

func (r VirtualMachineImplicitDataDiskFromSourceResource) detachImplicitDataDisk(data acceptance.TestData) string {
	// currently only supported in "eastus2" and "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }

    virtual_machine {
      detach_implicit_data_disk_on_deletion = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F2"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "%d-disk1"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}
`, data.RandomInteger, location, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (t VirtualMachineImplicitDataDiskFromSourceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
	resp, err := clients.Compute.VirtualMachinesClient.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving Implicit Data Disk %s: %+v", *id, err)
	}

	var disk *virtualmachines.DataDisk
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if profile := props.StorageProfile; profile != nil {
				if dataDisks := profile.DataDisks; dataDisks != nil {
					for _, dataDisk := range *dataDisks {
						if pointer.From(dataDisk.Name) == id.Name {
							disk = &dataDisk
							break
						}
					}
				}
			}
		}
	}

	return pointer.To(disk != nil), nil
}

func (VirtualMachineImplicitDataDiskFromSourceResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
	resp, err := client.Compute.VirtualMachinesClient.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving Implicit Data Disk %s: %+v", *id, err)
	}

	outputDisks := make([]virtualmachines.DataDisk, 0)
	var toBeDeletedDisk *virtualmachines.DataDisk
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.StorageProfile != nil {
			for _, disk := range *props.StorageProfile.DataDisks {
				if pointer.From(disk.Name) != id.Name {
					outputDisks = append(outputDisks, disk)
				} else {
					toBeDeletedDisk = &disk
				}
			}

			props.StorageProfile.DataDisks = &outputDisks

			// fixes #2485
			model.Identity = nil
			// fixes #1600
			model.Resources = nil

			if err := client.Compute.VirtualMachinesClient.CreateOrUpdateThenPoll(ctx, virtualMachineId, *model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return nil, fmt.Errorf("deleting implicit data disk %s: %+v", *id, err)
			}

			// delete the data disk which was created by Azure Service when creating this resource
			detachDataDisk := client.Features.VirtualMachine.DetachImplicitDataDiskOnDeletion
			if !detachDataDisk && toBeDeletedDisk != nil && toBeDeletedDisk.ManagedDisk != nil && toBeDeletedDisk.ManagedDisk.Id != nil {
				diskClient := client.Compute.DisksClient
				diskId, err := commonids.ParseManagedDiskID(*toBeDeletedDisk.ManagedDisk.Id)
				if err != nil {
					return nil, err
				}

				err = diskClient.DeleteThenPoll(ctx, *diskId)
				if err != nil {
					return nil, fmt.Errorf("deleting Managed Disk %s: %+v", *diskId, err)
				}
			}
		}
	}

	return pointer.To(true), nil
}

func (r VirtualMachineImplicitDataDiskFromSourceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "test" {
  name               = "acctestVMIDD-%d"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineImplicitDataDiskFromSourceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "import" {
  name               = azurerm_virtual_machine_implicit_data_disk_from_source.test.name
  virtual_machine_id = azurerm_virtual_machine_implicit_data_disk_from_source.test.virtual_machine_id
  lun                = azurerm_virtual_machine_implicit_data_disk_from_source.test.lun
  create_option      = azurerm_virtual_machine_implicit_data_disk_from_source.test.create_option
  disk_size_gb       = azurerm_virtual_machine_implicit_data_disk_from_source.test.disk_size_gb
  source_resource_id = azurerm_virtual_machine_implicit_data_disk_from_source.test.source_resource_id
}
`, r.basic(data))
}

func (VirtualMachineImplicitDataDiskFromSourceResource) managedServiceIdentity(data acceptance.TestData) string {
	// currently only supported in "eastus2" and "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F2"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "%d-disk1"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "test" {
  name               = "acctestVMIDD-%d"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_managed_disk.test.id
}
`, data.RandomInteger, location, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineImplicitDataDiskFromSourceResource) multipleDisks(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "first" {
  name               = "acctestVMIDD-%d"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.test.id
}

resource "azurerm_managed_disk" "second" {
  name                 = "%d-disk2"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_snapshot" "second" {
  name                = "acctestss2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.second.id
}

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "second" {
  name               = "acctestVMIDD2-%d"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "20"
  caching            = "ReadOnly"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.second.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineImplicitDataDiskFromSourceResource) readOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "test" {
  name               = "acctestVMIDD-%d"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "ReadOnly"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineImplicitDataDiskFromSourceResource) readWrite(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "test" {
  name               = "acctestVMIDD-%d"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "ReadWrite"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.test.id
}
`, r.template(data), data.RandomInteger)
}

func (VirtualMachineImplicitDataDiskFromSourceResource) writeAccelerator(data acceptance.TestData, enabled bool) string {
	// currently only supported in "eastus2" and "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_M8ms"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "%d-disk1"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "test" {
  name                      = "acctestVMIDD-%d"
  virtual_machine_id        = azurerm_virtual_machine.test.id
  lun                       = "0"
  create_option             = "Copy"
  disk_size_gb              = 20
  source_resource_id        = azurerm_snapshot.test.id
  write_accelerator_enabled = %t
}
`, data.RandomInteger, location, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, enabled)
}

func (VirtualMachineImplicitDataDiskFromSourceResource) template(data acceptance.TestData) string {
	// currently only supported in "eastus2" and "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F2"

  delete_os_disk_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hn%d"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "%d-disk1"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}
`, data.RandomInteger, location, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineImplicitDataDiskFromSourceResource) virtualMachineExtensionPrep(data acceptance.TestData) string {
	// currently only supported in "eastus2" and "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctestvm%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  vm_size               = "Standard_F4"

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_profile {
    computer_name  = "testvm"
    admin_username = "tfuser123"
    admin_password = "Password1234!"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_virtual_machine_extension" "test" {
  name                 = "random-script"
  virtual_machine_id   = azurerm_virtual_machine.test.id
  publisher            = "Microsoft.Azure.Extensions"
  type                 = "CustomScript"
  type_handler_version = "2.0"

  settings = <<SETTINGS
	{
		"commandToExecute": "hostname"
	}
SETTINGS

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, location, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineImplicitDataDiskFromSourceResource) virtualMachineExtensionComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctest%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "test" {
  name               = "acctestVMIDD-%d"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "11"
  caching            = "ReadWrite"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.test.id
}
`, r.virtualMachineExtensionPrep(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineImplicitDataDiskFromSourceResource) virtualMachineApplicationPrep(data acceptance.TestData) string {
	// currently only supported in "eastus2" and "westus2".
	location := "westus2"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "ACCTESTGAL%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_gallery_application" "test" {
  name              = "acc-app%[1]d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"
}

resource "azurerm_storage_account" "test" {
  name                     = "stacc%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "container"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "scripts"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_content         = "exit 0"
}

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  manage_action {
    install = "exit 0"
    remove  = "exit 0"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestni%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                  = "acctestvm%[1]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  network_interface_ids = [azurerm_network_interface.test.id]
  size                  = "Standard_F4"

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  admin_username                  = "tfuser123"
  admin_password                  = "Password1234!"
  disable_password_authentication = false

  os_disk {
    name                 = "myosdisk1"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  gallery_application {
    version_id = azurerm_gallery_application_version.test.id
    order      = 0
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, location)
}

func (r VirtualMachineImplicitDataDiskFromSourceResource) virtualMachineApplicationComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_disk" "test" {
  name                 = "acctest%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "test" {
  name               = "acctestVMIDD-%d"
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "11"
  caching            = "ReadWrite"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.test.id
}
`, r.virtualMachineApplicationPrep(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
