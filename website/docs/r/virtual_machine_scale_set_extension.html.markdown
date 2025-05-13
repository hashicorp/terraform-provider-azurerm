---
layout: "azurerm"
subcategory: "Compute"
page_title: "Azure Resource Manager: azurerm_virtual_machine_scale_set_extension"
description: |-
  Manages an Extension for a Virtual Machine Scale Set.
---

# azurerm_virtual_machine_scale_set_extension

Manages an Extension for a Virtual Machine Scale Set.

~> **Note:** This resource is not intended to be used with the `azurerm_virtual_machine_scale_set` resource - instead it's intended for this to be used with the `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` resources.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_linux_virtual_machine_scale_set" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Standard_F2"
  admin_username      = "adminuser"
  instances           = 1

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  network_interface {
    name = "example"

    ip_configuration {
      name = "internal"
    }
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }
}

resource "azurerm_virtual_machine_scale_set_extension" "example" {
  name                         = "example"
  virtual_machine_scale_set_id = azurerm_linux_virtual_machine_scale_set.example.id
  publisher                    = "Microsoft.Azure.Extensions"
  type                         = "CustomScript"
  type_handler_version         = "2.0"
  settings = jsonencode({
    "commandToExecute" = "echo $HOSTNAME"
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the Virtual Machine Scale Set Extension. Changing this forces a new resource to be created.

* `virtual_machine_scale_set_id` - (Required) The ID of the Virtual Machine Scale Set. Changing this forces a new resource to be created.

-> **Note:** This should be the ID from the `azurerm_linux_virtual_machine_scale_set` or `azurerm_windows_virtual_machine_scale_set` resource - when using the older `azurerm_virtual_machine_scale_set` resource extensions should instead be defined inline.

* `publisher` - (Required) Specifies the Publisher of the Extension. Changing this forces a new resource to be created.

* `type` - (Required) Specifies the Type of the Extension. Changing this forces a new resource to be created.

* `type_handler_version` - (Required) Specifies the version of the extension to use, available versions can be found using the Azure CLI.

~> **Note:** The `Publisher` and `Type` of Virtual Machine Scale Set Extensions can be found using the Azure CLI, via:

```shell
az vmss extension image list --location westus -o table
```

---

* `auto_upgrade_minor_version` - (Optional) Should the latest version of the Extension be used at Deployment Time, if one is available? This won't auto-update the extension on existing installation. Defaults to `true`.

* `automatic_upgrade_enabled` - (Optional) Should the Extension be automatically updated whenever the Publisher releases a new version of this VM Extension? 

* `failure_suppression_enabled` - (Optional) Should failures from the extension be suppressed? Possible values are `true` or `false`. Defaults to `false`.

-> **Note:** Operational failures such as not connecting to the VM will not be suppressed regardless of the `failure_suppression_enabled` value.

* `force_update_tag` - (Optional) A value which, when different to the previous value can be used to force-run the Extension even if the Extension Configuration hasn't changed.

* `protected_settings` - (Optional) A JSON String which specifies Sensitive Settings (such as Passwords) for the Extension.

~> **Note:** Keys within the `protected_settings` block are notoriously case-sensitive, where the casing required (e.g. TitleCase vs snakeCase) depends on the Extension being used. Please refer to the documentation for the specific Virtual Machine Extension you're looking to use for more information.

* `protected_settings_from_key_vault` - (Optional) A `protected_settings_from_key_vault` block as defined below.

~> **Note:** `protected_settings_from_key_vault` cannot be used with `protected_settings`

* `provision_after_extensions` - (Optional) An ordered list of Extension names which this should be provisioned after.

* `settings` - (Optional) A JSON String which specifies Settings for the Extension.

~> **Note:** Keys within the `settings` block are notoriously case-sensitive, where the casing required (e.g. TitleCase vs snakeCase) depends on the Extension being used. Please refer to the documentation for the specific Virtual Machine Extension you're looking to use for more information.

---

A `protected_settings_from_key_vault` block supports the following:

* `secret_url` - (Required) The URL to the Key Vault Secret which stores the protected settings.

* `source_vault_id` - (Required) The ID of the source Key Vault.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Machine Scale Set Extension.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Scale Set Extension.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Scale Set Extension.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Machine Scale Set Extension.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Scale Set Extension.

## Import

Virtual Machine Scale Set Extensions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_scale_set_extension.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/extensions/extension1
```
