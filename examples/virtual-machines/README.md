## Examples for the Virtual Machine resources

In 1.x versions of the Provider, Terraform has a single resource for Virtual Machines: `azurerm_virtual_machine`.

Version 2.0 of the Azure Provider introduces several new resources which supersede the existing `azurerm_virtual_machine` resource:

* `azurerm_linux_virtual_machine`
* `azurerm_windows_virtual_machine`

[More details can be found in this issue](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2807) - however these resources will replace the existing `azurerm_virtual_machine` resource in the long-term.

This directory contains 4 sub-directories:

* `./virtual_machine` - which are examples of how to use the `azurerm_virtual_machine` resource.
* `./linux` - which are examples of how to use the `azurerm_linux_virtual_machine` resource.
* `./windows` - which are examples of how to use the `azurerm_windows_virtual_machine` resource.
