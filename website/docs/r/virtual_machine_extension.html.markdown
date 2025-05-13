---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_extension"
description: |-
    Manages a Virtual Machine Extension to provide post deployment
    configuration and run automated tasks.
---

# azurerm_virtual_machine_extension

Manages a Virtual Machine Extension to provide post deployment configuration
and run automated tasks.

~> **Note:** Custom Script Extensions for Linux & Windows require that the `commandToExecute` returns a `0` exit code to be classified as successfully deployed. You can achieve this by appending `exit 0` to the end of your `commandToExecute`.

-> **Note:** Custom Script Extensions require that the Azure Virtual Machine Guest Agent is running on the Virtual Machine.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "acctvn"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "acctsub"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "acctni"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                = "example-machine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = file("~/.ssh/id_rsa.pub")
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_virtual_machine_extension" "example" {
  name                 = "hostname"
  virtual_machine_id   = azurerm_linux_virtual_machine.example.id
  publisher            = "Microsoft.Azure.Extensions"
  type                 = "CustomScript"
  type_handler_version = "2.0"

  settings = <<SETTINGS
 {
  "commandToExecute": "hostname && uptime"
 }
SETTINGS


  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual machine extension peering. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The ID of the Virtual Machine. Changing this forces a new resource to be created

* `publisher` - (Required) The publisher of the extension, available publishers can be found by using the Azure CLI. Changing this forces a new resource to be created.

* `type` - (Required) The type of extension, available types for a publisher can be found using the Azure CLI.

~> **Note:** The `Publisher` and `Type` of Virtual Machine Extensions can be found using the Azure CLI, via:

```shell
az vm extension image list --location westus -o table
```

* `type_handler_version` - (Required) Specifies the version of the extension to use, available versions can be found using the Azure CLI.

* `auto_upgrade_minor_version` - (Optional) Specifies if the platform deploys the latest minor version update to the `type_handler_version` specified.

* `automatic_upgrade_enabled` - (Optional) Should the Extension be automatically updated whenever the Publisher releases a new version of this VM Extension? 
* `settings` - (Optional) The settings passed to the extension, these are specified as a JSON object in a string.

~> **Note:** Certain VM Extensions require that the keys in the `settings` block are case sensitive. If you're seeing unhelpful errors, please ensure the keys are consistent with how Azure is expecting them (for instance, for the `JsonADDomainExtension` extension, the keys are expected to be in `TitleCase`.)

* `failure_suppression_enabled` - (Optional) Should failures from the extension be suppressed? Possible values are `true` or `false`. Defaults to `false`.

-> **Note:** Operational failures such as not connecting to the VM will not be suppressed regardless of the `failure_suppression_enabled` value.

* `protected_settings` - (Optional) The protected_settings passed to the extension, like settings, these are specified as a JSON object in a string.

~> **Note:** Certain VM Extensions require that the keys in the `protected_settings` block are case sensitive. If you're seeing unhelpful errors, please ensure the keys are consistent with how Azure is expecting them (for instance, for the `JsonADDomainExtension` extension, the keys are expected to be in `TitleCase`.)

* `protected_settings_from_key_vault` - (Optional) A `protected_settings_from_key_vault` block as defined below.

~> **Note:** `protected_settings_from_key_vault` cannot be used with `protected_settings`

* `provision_after_extensions` - (Optional) Specifies the collection of extension names after which this extension needs to be provisioned.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `protected_settings_from_key_vault` block supports the following:

* `secret_url` - (Required) The URL to the Key Vault Secret which stores the protected settings.

* `source_vault_id` - (Required) The ID of the source Key Vault.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Machine Extension.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Extension.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Extension.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Machine Extension.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Extension.

## Import

Virtual Machine Extensions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_extension.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/myVM/extensions/extensionName
```
