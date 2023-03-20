---
subcategory: "Hybrid Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hybrid_compute_machine_extension"
description: |-
  Manages a Hybrid Compute Machine Extension.
---

# azurerm_hybrid_compute_machine_extension

Manages a Hybrid Compute Machine Extension.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

data "azurerm_hybrid_compute_machine" "example" {
  name                = "existing-hcmachine"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_hybrid_compute_machine_extension" "example" {
  name                      = "example"
  location                  = "West Europe"
  hybrid_compute_machine_id = data.azurerm_hybrid_compute_machine.example.id
  publisher                 = "Microsoft.Azure.Monitor"
  type                      = "AzureMonitorLinuxAgent"
}
```

## Arguments Reference

The following arguments are supported:

* `hybrid_compute_machine_id` - (Required) The ID of the Hybrid Compute Machine Extension. Changing this forces a new Hybrid Compute Machine Extension to be created.

* `location` - (Required) The Azure Region where the Hybrid Compute Machine Extension should exist. Changing this forces a new Hybrid Compute Machine Extension to be created.

* `name` - (Required) The name which should be used for this Hybrid Compute Machine Extension. Changing this forces a new Hybrid Compute Machine Extension to be created.

* `publisher` - (Required) The name of the extension handler publisher, such as `Microsoft.Azure.Monitor`. Changing this forces a new Hybrid Compute Machine Extension to be created.

* `type` - (Required) Specifies the type of the extension. For example `CustomScriptExtension` or `AzureMonitorLinuxAgent`. Changing this forces a new Hybrid Compute Machine Extension to be created.

---

* `automatic_upgrade_enabled` - (Optional) Indicates whether the extension should be automatically upgraded by the platform if there is a newer version available. Supported values are `true` and `false`.

* `force_update_tag` - (Optional) How the extension handler should be forced to update even if the extension configuration has not changed.

* `protected_settings` - (Optional) Json formatted protected settings for the extension.

* `settings` - (Optional) Json formatted public settings for the extension.

* `tags` - (Optional) A mapping of tags which should be assigned to the Hybrid Compute Machine Extension.

* `type_handler_version` - (Optional) Specifies the version of the script handler.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Hybrid Compute Machine Extension.

* `instance_view` - A `instance_view` block as defined below.

---

A `instance_view` block exports the following:

* `name` - The name of this Hybrid Compute Machine Extension.

* `status` - A `status` block as defined below.

* `type` - Specifies the type of the extension.

* `type_handler_version` - Specifies the version of the script handler.

---

A `status` block exports the following:

* `code` - The status code of the instance.

* `display_status` - The short localizable label for the status.

* `level` - The level code. Possible values are `Info`, `Warning` or `Error`.

* `message` - The detailed status message, including for alerts and error messages.

* `time` - The time of the status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Hybrid Compute Machine Extension.
* `read` - (Defaults to 5 minutes) Used when retrieving the Hybrid Compute Machine Extension.
* `update` - (Defaults to 30 minutes) Used when updating the Hybrid Compute Machine Extension.
* `delete` - (Defaults to 30 minutes) Used when deleting the Hybrid Compute Machine Extension.

## Import

Hybrid Compute Machine Extensions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hybrid_compute_machine_extension.example C:/Program Files/Git/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HybridCompute/machines/hcmachine1/extensions/ext1
```