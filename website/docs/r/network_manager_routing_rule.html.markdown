---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_routing_rule"
description: |-
  Manages a Network Manager Routing Rule.
---

# azurerm_network_manager_routing_rule

Manages a Network Manager Routing Rule.

## Example Usage

```hcl
resource "azurerm_network_manager_routing_rule" "example" {
  name = "example"
  rule_collection_id = "TODO"

  destination {
    type = "TODO"
    address = "TODO"    
  }

  next_hop {
    address = "TODO"
    type = "TODO"    
  }
}
```

## Arguments Reference

The following arguments are supported:

* `destination` - (Required) A `destination` block as defined below.

* `name` - (Required) The name which should be used for this Network Manager Routing Rule. Changing this forces a new Network Manager Routing Rule to be created.

* `next_hop` - (Required) A `next_hop` block as defined below.

* `rule_collection_id` - (Required) The ID of the TODO. Changing this forces a new Network Manager Routing Rule to be created.

---

* `description` - (Optional) TODO.

---

A `destination` block supports the following:

* `address` - (Required) TODO. Changing this forces a new Network Manager Routing Rule to be created.

* `type` - (Required) TODO. Changing this forces a new Network Manager Routing Rule to be created.

---

A `next_hop` block supports the following:

* `address` - (Required) TODO. Changing this forces a new Network Manager Routing Rule to be created.

* `type` - (Required) TODO. Changing this forces a new Network Manager Routing Rule to be created.

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