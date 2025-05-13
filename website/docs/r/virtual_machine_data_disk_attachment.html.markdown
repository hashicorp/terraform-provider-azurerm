---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_data_disk_attachment"
description: |-
  Manages attaching a Disk to a Virtual Machine.
---

# azurerm_virtual_machine_data_disk_attachment

Manages attaching a Disk to a Virtual Machine.

~> **Note:** Data Disks can be attached either directly on the `azurerm_virtual_machine` resource, or using the `azurerm_virtual_machine_data_disk_attachment` resource - but the two cannot be used together. If both are used against the same Virtual Machine, spurious changes will occur.

-> **Note:** only Managed Disks are supported via this separate resource, Unmanaged Disks can be attached using the `storage_data_disk` block in the `azurerm_virtual_machine` resource.

## Example Usage

```hcl
variable "prefix" {
  default = "example"
}

locals {
  vm_name = "${var.prefix}-vm"
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "main" {
  name                = "${var.prefix}-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "main" {
  name                = "${var.prefix}-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.internal.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "example" {
  name                  = local.vm_name
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.main.id]
  vm_size               = "Standard_F2"

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
    computer_name  = local.vm_name
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_managed_disk" "example" {
  name                 = "${local.vm_name}-disk1"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 10
}

resource "azurerm_virtual_machine_data_disk_attachment" "example" {
  managed_disk_id    = azurerm_managed_disk.example.id
  virtual_machine_id = azurerm_virtual_machine.example.id
  lun                = "10"
  caching            = "ReadWrite"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_machine_id` - (Required) The ID of the Virtual Machine to which the Data Disk should be attached. Changing this forces a new resource to be created.

* `managed_disk_id` - (Required) The ID of an existing Managed Disk which should be attached. Changing this forces a new resource to be created.

* `lun` - (Required) The Logical Unit Number of the Data Disk, which needs to be unique within the Virtual Machine. Changing this forces a new resource to be created.

* `caching` - (Required) Specifies the caching requirements for this Data Disk. Possible values include `None`, `ReadOnly` and `ReadWrite`.

* `create_option` - (Optional) The Create Option of the Data Disk, such as `Empty` or `Attach`. Defaults to `Attach`. Changing this forces a new resource to be created.

* `write_accelerator_enabled` - (Optional) Specifies if Write Accelerator is enabled on the disk. This can only be enabled on `Premium_LRS` managed disks with no caching and [M-Series VMs](https://docs.microsoft.com/azure/virtual-machines/workloads/sap/how-to-enable-write-accelerator). Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Machine Data Disk attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Data Disk Attachment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Data Disk Attachment.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Machine Data Disk Attachment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Data Disk Attachment.

## Import

Virtual Machines Data Disk Attachments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_data_disk_attachment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1/dataDisks/disk1
```

-> **Note:** This is a Terraform Unique ID matching the format: `{virtualMachineID}/dataDisks/{diskName}`
