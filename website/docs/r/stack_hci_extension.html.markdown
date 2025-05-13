---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_extension"
description: |-
  Manages an Azure Stack HCI Extension.
---

# azurerm_stack_hci_extension

Manages an Azure Stack HCI Extension.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-hci-ext"
  location = "West Europe"
}

resource "azurerm_stack_hci_extension" "example" {
  name                               = "AzureMonitorWindowsAgent"
  arc_setting_id                     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-hci/providers/Microsoft.AzureStackHCI/clusters/hci-cl/arcSettings/default"
  publisher                          = "Microsoft.Azure.Monitor"
  type                               = "MicrosoftMonitoringAgent"
  auto_upgrade_minor_version_enabled = true
  automatic_upgrade_enabled          = true
  type_handler_version               = "1.22.0"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Extension. Changing this forces a new resource to be created.

* `arc_setting_id` - (Required) The ID of the Azure Stack HCI Cluster Arc Setting. Changing this forces a new resource to be created.

* `publisher` - (Required) The name of the extension handler publisher, such as `Microsoft.Azure.Monitor`. Changing this forces a new resource to be created.

* `type` - (Required) Specifies the type of the extension. For example `CustomScriptExtension` or `AzureMonitorLinuxAgent`. Changing this forces a new resource to be created.

* `auto_upgrade_minor_version_enabled` - (Optional) Indicates whether the extension should use a newer minor version if one is available at deployment time. Once deployed, however, the extension will not upgrade minor versions unless redeployed, even with this property set to true. Changing this forces a new resource to be created. Possible values are `true` and `false`. Defaults to `true`.

* `automatic_upgrade_enabled` - (Optional) Indicates whether the extension should be automatically upgraded by the platform if there is a newer version available. Possible values are `true` and `false`. Defaults to `true`.

* `protected_settings` - (Optional) The json formatted protected settings for the extension.

* `settings` - (Optional) The json formatted public settings for the extension.

* `type_handler_version` - (Optional) Specifies the version of the script handler.

-> **Note:** `type_handler_version` cannot be set when `automatic_upgrade_enabled` is set to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Stack HCI Extension.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Extension.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Extension.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Extension.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Extension.

## Import

Azure Stack HCI Extension can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_extension.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/clusters/cluster1/arcSettings/default/extensions/extension1
```
