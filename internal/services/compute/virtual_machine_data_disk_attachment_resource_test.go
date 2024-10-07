// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualMachineDataDiskAttachmentResource struct{}

func TestAccVirtualMachineDataDiskAttachment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_machine_id").Exists(),
				check.That(data.ResourceName).Key("managed_disk_id").Exists(),
				check.That(data.ResourceName).Key("lun").HasValue("0"),
				check.That(data.ResourceName).Key("caching").HasValue("None"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineDataDiskAttachment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_machine_data_disk_attachment"),
		},
	})
}

func TestAccVirtualMachineDataDiskAttachment_destroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccVirtualMachineDataDiskAttachment_multipleDisks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "first")
	r := VirtualMachineDataDiskAttachmentResource{}

	secondResourceName := "azurerm_virtual_machine_data_disk_attachment.second"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleDisks(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_machine_id").Exists(),
				check.That(data.ResourceName).Key("managed_disk_id").Exists(),
				check.That(data.ResourceName).Key("lun").HasValue("10"),
				check.That(data.ResourceName).Key("caching").HasValue("None"),
				acceptance.TestCheckResourceAttrSet(secondResourceName, "virtual_machine_id"),
				acceptance.TestCheckResourceAttrSet(secondResourceName, "managed_disk_id"),
				acceptance.TestCheckResourceAttr(secondResourceName, "lun", "20"),
				acceptance.TestCheckResourceAttr(secondResourceName, "caching", "ReadOnly"),
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

func TestAccVirtualMachineDataDiskAttachment_updatingCaching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("caching").HasValue("None"),
			),
		},
		{
			Config: r.readOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("caching").HasValue("ReadOnly"),
			),
		},
		{
			Config: r.readWrite(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("caching").HasValue("ReadWrite"),
			),
		},
	})
}

func TestAccVirtualMachineDataDiskAttachment_updatingWriteAccelerator(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.writeAccelerator(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("write_accelerator_enabled").HasValue("false"),
			),
		},
		{
			Config: r.writeAccelerator(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("write_accelerator_enabled").HasValue("true"),
			),
		},
		{
			Config: r.writeAccelerator(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("write_accelerator_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccVirtualMachineDataDiskAttachment_managedServiceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedServiceIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_machine_id").Exists(),
				check.That(data.ResourceName).Key("managed_disk_id").Exists(),
				check.That(data.ResourceName).Key("lun").HasValue("0"),
				check.That(data.ResourceName).Key("caching").HasValue("None"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualMachineDataDiskAttachment_virtualMachineExtension(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualMachineExtensionPrep(data),
		},
		{
			Config: r.virtualMachineExtensionComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_machine_id").Exists(),
				check.That(data.ResourceName).Key("managed_disk_id").Exists(),
			),
		},
	})
}

func TestAccVirtualMachineDataDiskAttachment_virtualMachineApplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualMachineApplicationPrep(data),
		},
		{
			Config: r.virtualMachineApplicationComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_machine_id").Exists(),
				check.That(data.ResourceName).Key("managed_disk_id").Exists(),
				check.That("azurerm_linux_virtual_machine.test").Key("gallery_application.0.order").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualMachineDataDiskAttachment_importImplicitDataDisk(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_machine_data_disk_attachment", "test")
	r := VirtualMachineDataDiskAttachmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.importImplicitDataDisk(data),
		},
		{
			ResourceName:  data.ResourceName,
			ImportState:   true,
			ImportStateId: fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Compute/virtualMachines/acctvm-%d/dataDisks/acctestVMIDD-%d", os.Getenv("ARM_SUBSCRIPTION_ID"), data.RandomInteger, data.RandomInteger, data.RandomInteger),
			ExpectError:   regexp.MustCompile("the value of `create_option` for the imported `azurerm_virtual_machine_data_disk_attachment` instance must be `Attach` or `Empty`"),
		},
	})
}

func (t VirtualMachineDataDiskAttachmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
	resp, err := clients.Compute.VirtualMachinesClient.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving Data Disk Attachment %q for %s: %+v", *id, virtualMachineId, err)
	}

	var disk *virtualmachines.DataDisk
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if profile := props.StorageProfile; profile != nil {
				if dataDisks := profile.DataDisks; dataDisks != nil {
					for _, dataDisk := range *dataDisks {
						// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.EqualFold`
						if *dataDisk.Name == id.Name {
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

func (VirtualMachineDataDiskAttachmentResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DataDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
	resp, err := client.Compute.VirtualMachinesClient.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving Data Disk Attachment %q for %s: %+v", *id, virtualMachineId, err)
	}

	outputDisks := make([]virtualmachines.DataDisk, 0)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.StorageProfile != nil {
			for _, disk := range *props.StorageProfile.DataDisks {
				// deliberately not using strings.Equals as this is case sensitive
				if *disk.Name == id.Name {
					continue
				}

				outputDisks = append(outputDisks, disk)
			}
			model.Properties.StorageProfile.DataDisks = &outputDisks

			// fixes #2485
			model.Identity = nil
			// fixes #1600
			model.Resources = nil

			if err := client.Compute.VirtualMachinesClient.CreateOrUpdateThenPoll(ctx, virtualMachineId, *model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return nil, fmt.Errorf("updating Virtual Machine %q: %+v", id, err)
			}
		}
	}

	return utils.Bool(true), nil
}

func (r VirtualMachineDataDiskAttachmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, r.template(data))
}

func (r VirtualMachineDataDiskAttachmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "import" {
  managed_disk_id    = azurerm_virtual_machine_data_disk_attachment.test.managed_disk_id
  virtual_machine_id = azurerm_virtual_machine_data_disk_attachment.test.virtual_machine_id
  lun                = azurerm_virtual_machine_data_disk_attachment.test.lun
  caching            = azurerm_virtual_machine_data_disk_attachment.test.caching
}
`, r.basic(data))
}

func (VirtualMachineDataDiskAttachmentResource) managedServiceIdentity(data acceptance.TestData) string {
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

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineDataDiskAttachmentResource) multipleDisks(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "first" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "10"
  caching            = "None"
}

resource "azurerm_managed_disk" "second" {
  name                 = "%d-disk2"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "second" {
  managed_disk_id    = azurerm_managed_disk.second.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "20"
  caching            = "ReadOnly"
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualMachineDataDiskAttachmentResource) readOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "ReadOnly"
}
`, r.template(data))
}

func (r VirtualMachineDataDiskAttachmentResource) readWrite(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "ReadWrite"
}
`, r.template(data))
}

func (VirtualMachineDataDiskAttachmentResource) writeAccelerator(data acceptance.TestData, enabled bool) string {
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

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id           = azurerm_managed_disk.test.id
  virtual_machine_id        = azurerm_virtual_machine.test.id
  lun                       = "0"
  caching                   = "None"
  write_accelerator_enabled = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, enabled)
}

func (VirtualMachineDataDiskAttachmentResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualMachineDataDiskAttachmentResource) virtualMachineExtensionPrep(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualMachineDataDiskAttachmentResource) virtualMachineExtensionComplete(data acceptance.TestData) string {
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

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "11"
  caching            = "ReadWrite"
}
`, r.virtualMachineExtensionPrep(data), data.RandomInteger)
}

func (VirtualMachineDataDiskAttachmentResource) virtualMachineApplicationPrep(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r VirtualMachineDataDiskAttachmentResource) virtualMachineApplicationComplete(data acceptance.TestData) string {
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

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_linux_virtual_machine.test.id
  lun                = "11"
  caching            = "ReadWrite"
}
`, r.virtualMachineApplicationPrep(data), data.RandomInteger)
}

func (r VirtualMachineDataDiskAttachmentResource) importImplicitDataDisk(data acceptance.TestData) string {
	// currently implicit data disk is only supported in "eastus2" and "westus2".
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

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "test" {
  name               = "acctestVMIDD-%d"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.test.id
}

resource "azurerm_managed_disk" "test2" {
  name                 = "%d-disk2"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test2.id
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "1"
  caching            = "None"
}
`, data.RandomInteger, location, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
