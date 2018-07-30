---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_state"
sidebar_current: "docs-azurerm-resource-compute-virtualmachine-state"
description: |-
    Manages a Virtual Machine state to allow for a machine to be Running,
    Powered Off and Deallocated without destroying/re-creating
---

# azurerm_virtual_machine_state

Manages a Virtual Machine Extension to provide post deployment configuration
and run automated tasks.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "West US"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctni"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_virtual_machine" "test" {
  name                  = "acctvm"
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
    name          = "myosdisk1"
    vhd_uri       = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}/myosdisk1.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags {
    environment = "staging"
  }
}

resource "azurerm_virtual_machine_extension" "test" {
  name                 = "hostname"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_machine_name = "${azurerm_virtual_machine.test.name}"
  publisher            = "Microsoft.Azure.Extensions"
  type                 = "CustomScript"
  type_handler_version = "2.0"

  settings = <<SETTINGS
	{
		"commandToExecute": "hostname && uptime"
	}
SETTINGS

  tags {
    environment = "Production"
  }
}

resource "azurerm_virtual_machine_state" "test" {
  virtual_machine_name = "${azurerm_virtual_machine.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  state = "Deallocated"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual machine extension peering. Changing
    this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the virtual network. Changing this forces a new resource to be
    created.

* `state` - (Required) The desired state of the Virtual Machine. Valid options are `running`,
    `stopped` (stopped but will still incur costs) and `deallocated` (stopped and won't incur costs).

## Attributes Reference

The following attributes are exported:

* `id` - The state of the virtual machine.

## Import

Virtual Machine States can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_state.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/myVM
```
