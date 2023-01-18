---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_subscription_connection"
description: |-
  Manages a Network Manager Subscription Connection.
---

# azurerm_network_manager_subscription_connection

Manages a Network Manager Subscription Connection which may cross tenants.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "example" {
  name                = "example-networkmanager"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["SecurityAdmin"]
}

resource "azurerm_network_manager_subscription_connection" "example" {
  name               = "example-nsnmc"
  subscription_id    = data.azurerm_subscription.current.id
  network_manager_id = azurerm_network_manager.example.id
  description        = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Subscription Network Manager Connection. Changing this forces a new Network Subscription Network Manager Connection to be created.

* `subscription_id` - (Required) Specifies the ID of the target Subscription. Changing this forces a new resource to be created.

* `network_manager_id` - (Required) Specifies the ID of the Network Manager which the Subscription is connected to.

* `description` - (Optional) A description of the Network Manager Subscription Connection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Subscription Connection.

* `connection_state` - The Connection state of the Network Manager Subscription Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Subscription Network Manager Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Subscription Network Manager Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Network Subscription Network Manager Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Subscription Network Manager Connection.

## Import

Network Subscription Network Manager Connection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_subscription_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/networkManagerConnections/networkManagerConnection1
```
