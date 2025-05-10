---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_watcher"
description: |-
  Manages a Network Watcher.

---

# azurerm_network_watcher

Manages a Network Watcher.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "production-nwwatcher"
  location = "West Europe"
}

resource "azurerm_network_watcher" "example" {
  name                = "production-nwwatcher"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Network Watcher. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Network Watcher. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Watcher.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Watcher.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Watcher.
* `update` - (Defaults to 30 minutes) Used when updating the Network Watcher.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Watcher.

## Import

Network Watchers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_watcher.watcher1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkWatchers/watcher1
```
