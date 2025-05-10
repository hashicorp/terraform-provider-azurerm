---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_implicit_data_disk_from_source"
description: |-
  Manages an implicit Data Disk of a Virtual Machine.
---

# azurerm_virtual_machine_implicit_data_disk_from_source

Manages an implicit Data Disk of a Virtual Machine.

~> **Note:** The Implicit Data Disk will be deleted instantly after this resource is destroyed. If you want to detach this disk only, you may set `detach_implicit_data_disk_on_deletion` field to `true` within the `virtual_machine` block in the provider `features` block.

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

resource "azurerm_snapshot" "example" {
  name                = "${local.vm_name}-snapshot1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.example.id
}

resource "azurerm_virtual_machine_implicit_data_disk_from_source" "example" {
  name               = "${local.vm_name}-implicitdisk1"
  virtual_machine_id = azurerm_virtual_machine.test.id
  lun                = "0"
  caching            = "None"
  create_option      = "Copy"
  disk_size_gb       = 20
  source_resource_id = azurerm_snapshot.test.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Data Disk. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The ID of the Virtual Machine to which the Data Disk should be attached. Changing this forces a new resource to be created.

* `lun` - (Required) The Logical Unit Number of the Data Disk, which needs to be unique within the Virtual Machine. Changing this forces a new resource to be created.

* `create_option` - (Required) Specifies the Create Option of the Data Disk. The only possible value is `Copy`. Changing this forces a new resource to be created.

* `disk_size_gb` - (Required) Specifies the size of the Data Disk in gigabytes. Changing this forces a new resource to be created.

* `source_resource_id` - (Required) The ID of the source resource which this Data Disk was created from. Changing this forces a new resource to be created.

* `caching` - (Optional) Specifies the caching requirements for this Data Disk. Possible values are `ReadOnly` and `ReadWrite`.

* `write_accelerator_enabled` - (Optional) Specifies if Write Accelerator is enabled on the disk. This can only be enabled on `Premium_LRS` managed disks with no caching and [M-Series VMs](https://docs.microsoft.com/azure/virtual-machines/workloads/sap/how-to-enable-write-accelerator). Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the implicit Data Disk of the Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the implicit Data Disk of the Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the implicit Data Disk of the Virtual Machine.
* `update` - (Defaults to 30 minutes) Used when updating the implicit Data Disk of the Virtual Machine.
* `delete` - (Defaults to 30 minutes) Used when deleting the implicit Data Disk of the Virtual Machine.

## Import

The implicit Data Disk of the Virtual Machine can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_implicit_data_disk_from_source.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1/dataDisks/disk1
```

-> **Note:** This is a Terraform Unique ID matching the format: `{virtualMachineID}/dataDisks/{diskName}`
