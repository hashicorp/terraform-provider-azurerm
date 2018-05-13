---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_data_disk_attachment"
sidebar_current: "docs-azurerm-resource-compute-virtual-machine-data-disk-attachment"
description: |-
  Manages attaching a Disk to a Virtual Machine.
---

# azurerm_virtual_machine_data_disk_attachment

Manages attaching a Disk to a Virtual Machine.

~> **NOTE:** Data Disks can be attached either directly on the `azurerm_virtual_machine` resource, or using the `azurerm_virtual_machine_data_disk_attachment` resource - but the two cannot be used together. If both are used against the same Virtual Machine, spurious changes will occur.

## Example Usage

```hcl
variable "prefix" {
  default = "example"
}

locals {
  vm_name = "${var.prefix}-vm"
}

resource "azurerm_resource_group" "main" {
  name = "${var.prefix}-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "main" {
  name                = "${var.prefix}-network"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.main.name}"
  virtual_network_name = "${azurerm_virtual_network.main.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "main" {
  name                = "${var.prefix}-nic"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"

  ip_configuration {
    name                          = "internal"
    subnet_id                     = "${azurerm_subnet.internal.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_virtual_machine" "test" {
  name                  = "${local.vm_name}"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "${local.vm_name}"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  name               = "${local.vm_name}-data1"
  virtual_machine_id = "${azurerm_virtual_machine.main.id}"
  create_option      = "Empty"
  managed_disk_type  = "Standard_LRS"
  disk_size_gb       = 10
  lun                = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this Disk Attachment, which needs to be unique within the Virtual Machine. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The ID of the Virtual Machine to which the Data Disk should be attached. Changing this forces a new resource to be created.

* `create_option` - (Required) The Create Option of the Data Disk, such as `Empty` or `Attach`. Changing this forces a new resource to be created.

* `lun` - (Required) The Logical Unit Number of the Data Disk, which needs to be unique within the Virtual Machine.

* `vhd_uri` - (Optional) The URI of a Blob in a Storage Account where the VHD for this Disk should be placed. Cannot be specified when `managed_disk_id` or `managed_disk_type` is specified.

* `managed_disk_id` - (Optional) The ID of an existing Managed Disk which should be attached. When set, `create_option` should be set to `Attach`.

* `managed_disk_type` - (Optional) Specifies the type of managed disk to create. Value you must be either `Standard_LRS` or `Premium_LRS`. Cannot be used when `vhd_uri` is specified.

* `caching` - (Optional) Specifies the caching requirements for this Data Disk, such as `ReadWrite`.

* `disk_size_gb` - (Optional) Specifies the size of the Data Disk in GB.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Machine Data Disk attachment.

## Import

Virtual Machines Data Disk Attachments can be imported using the `resource id`, e.g.

```hcl
terraform import azurerm_virtual_machine_data_disk_attachment.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/virtualMachines/machine1/dataDisks/disk1
```
