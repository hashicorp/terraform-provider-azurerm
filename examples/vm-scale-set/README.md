## Examples for the Virtual Machine Scale Set resources

In 1.x versions of the Provider, Terraform has a single resource for Virtual Machine Scale Sets: `azurerm_virtual_machine_scale_set`.

Version 2.0 of the Azure Provider will introduce several new resources which supersede the existing `azurerm_virtual_machine_scale_set` resource:

* `azurerm_linux_virtual_machine_scale_set`
* `azurerm_virtual_machine_scale_set_extension`
* `azurerm_windows_virtual_machine_scale_set`

[More details can be found in this issue](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2807) - however these resources will replace the existing `azurerm_virtual_machine_scale_set` resource in the long-term.

Version 1.43 of the Azure Provider will allow you to opt-into the Beta of these resources to provide feedback - [more information on how to do that can be found in this guide](https://terraform.io/docs/providers/azurerm/guides/2.0-beta.html).

This directory contains 4 sub-directories:

* `./virtual_machine_scale_set` - which are examples of how to use the `azurerm_virtual_machine_scale_set` resource.
* `./linux` - which are examples of how to use the `azurerm_linux_virtual_machine_scale_set` resource.
* `./extensions` - which are examples of how to use the `azurerm_virtual_machine_scale_set_extension` resource.
* `./windows` - which are examples of how to use the `azurerm_windows_virtual_machine_scale_set` resource.
