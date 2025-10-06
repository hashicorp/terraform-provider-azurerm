---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_routing_rule"
description: |-
  Manages a Network Manager Routing Rule.
---


# azurerm_network_manager_routing_rule

Manages a Network Manager Routing Rule.

!> **Note:** Terraform has enabled force deletion. This setting deletes the resource even if it's part of a deployed configuration. If the configuration is deployed, the service will perform a cleanup deployment in the background before the deletion.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "example" {
  name                = "example-network-manager"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Routing"]
}

resource "azurerm_network_manager_network_group" "example" {
  name               = "example-network-group"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_network_manager_routing_configuration" "example" {
  name               = "example-routing-configuration"
  network_manager_id = azurerm_network_manager.example.id
}

resource "azurerm_network_manager_routing_rule_collection" "example" {
  name                     = "example-routing-rule-collection"
  routing_configuration_id = azurerm_network_manager_routing_configuration.example.id
  network_group_ids        = [azurerm_network_manager_network_group.example.id]
  description              = "example routing rule collection"
}

resource "azurerm_network_manager_routing_rule" "example" {
  name               = "example-routing-rule"
  rule_collection_id = azurerm_network_manager_routing_rule_collection.example.id
  description        = "example routing rule"

  destination {
    type    = "AddressPrefix"
    address = "10.0.0.0/24"
  }

  next_hop {
    type = "VirtualNetworkGateway"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Network Manager Routing Rule. Changing this forces a new resource to be created.

* `rule_collection_id` - (Required) The ID of the Network Manager Routing Rule Collection to which this rule belongs. Changing this forces a new resource to be created.

* `destination` - (Required) A `destination` block as defined below.

* `next_hop` - (Required) A `next_hop` block as defined below.

* `description` - (Optional) A description for the routing rule.

---

A `destination` block supports the following:

* `address` - (Required) The destination address.

* `type` - (Required) The type of destination. Possible values are `AddressPrefix` and `ServiceTag`.

---

A `next_hop` block supports the following:

* `type` - (Required) The type of next hop. Possible values are `Internet`, `NoNextHop`, `VirtualAppliance`, `VirtualNetworkGateway` and `VnetLocal`.

* `address` - (Optional) The address of the next hop. This is required if the next hop type is `VirtualAppliance`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Routing Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Routing Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Routing Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Routing Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Routing Rule.

## Import

Network Manager Routing Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_routing_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/manager1/routingConfigurations/conf1/ruleCollections/collection1/rules/rule1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2024-05-01
