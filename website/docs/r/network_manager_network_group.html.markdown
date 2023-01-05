---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_network_group"
description: |-
  Manages a Network Manager Network Group.
---

# azurerm_network_manager_network_group

Manages a Network Manager Network Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity", "SecurityAdmin"]
  description    = "example network manager"
}

resource "azurerm_network_manager_network_group" "example" {
  name               = "example-group"
  network_manager_id = azurerm_network_manager.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Manager Network Group. Changing this forces a new Network Manager Network Group to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Manager. Changing this forces a new Network Manager Network Group to be created.

* `description` - (Optional) A description of the Network Manager Network Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Network Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Network Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Network Group.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Network Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Network Group.

## Import

Network Manager Network Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_network_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/networkManagers/networkManager1/networkGroups/networkGroup1
```
