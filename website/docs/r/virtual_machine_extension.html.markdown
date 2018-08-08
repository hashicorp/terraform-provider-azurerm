---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_extension"
sidebar_current: "docs-azurerm-resource-compute-virtual-machine-extension"
description: |-
    Manages a Virtual Machine Extension to provide post deployment
    configuration and run automated tasks.
---

# azurerm_virtual_machine_extension

Manages a Virtual Machine Extension to provide post deployment configuration
and run automated tasks.

~> **NOTE:** Custom Script Extensions for Linux & Windows require that the `commandToExecute` returns a `0` exit code to be classified as successfully deployed. You can achieve this by appending `exit 0` to the end of your `commandToExecute`.

-> **NOTE:** Custom Script Extensions require that the Azure Virtual Machine Guest Agent is running on the Virtual Machine.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_virtual_network" "example" {
  # ...
}

resource "azurerm_subnet" "example" {
  # ...
}

resource "azurerm_network_interface" "example" {
  # ...
}

resource "azurerm_storage_account" "example" {
  # ...
}

resource "azurerm_storage_container" "example" {
  # ...
}

resource "azurerm_virtual_machine" "example" {
  # ...
}

resource "azurerm_virtual_machine_extension" "example" {
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual machine extension peering. Changing
    this forces a new resource to be created.

* `location` - (Required) The location where the extension is created. Changing
    this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the virtual network. Changing this forces a new resource to be
    created.

* `virtual_machine_name` - (Required) The name of the virtual machine. Changing
    this forces a new resource to be created.

* `publisher` - (Required) The publisher of the extension, available publishers
    can be found by using the Azure CLI.

* `type` - (Required) The type of extension, available types for a publisher can
    be found using the Azure CLI.

~> **Note:** The `Publisher` and `Type` of Virtual Machine Extensions can be found using the Azure CLI, via:
```shell
$ az vm extension image list --location westus -o table
```

* `type_handler_version` - (Required) Specifies the version of the extension to
    use, available versions can be found using the Azure CLI.

* `auto_upgrade_minor_version` - (Optional) Specifies if the platform deploys
    the latest minor version update to the `type_handler_version` specified.

* `settings` - (Required) The settings passed to the extension, these are
    specified as a JSON object in a string.

~> **Please Note:** Certain VM Extensions require that the keys in the `settings` block are case sensitive. If you're seeing unhelpful errors, please ensure the keys are consistent with how Azure is expecting them (for instance, for the `JsonADDomainExtension` extension, the keys are expected to be in `TitleCase`.)

* `protected_settings` - (Optional) The protected_settings passed to the
    extension, like settings, these are specified as a JSON object in a string.

~> **Please Note:** Certain VM Extensions require that the keys in the `protected_settings` block are case sensitive. If you're seeing unhelpful errors, please ensure the keys are consistent with how Azure is expecting them (for instance, for the `JsonADDomainExtension` extension, the keys are expected to be in `TitleCase`.)

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Machine Extension.

## Import

Virtual Machine Extensions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_extension.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/myVM/extensions/hostname
```
