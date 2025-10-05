---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_routing_rule_collection"
description: |-
  Manages a Network Manager Routing Rule Collection.
---

# azurerm_network_manager_routing_rule_collection

Manages a Network Manager Routing Rule Collection.

!> **Note:** Terraform has enabled force deletion. This setting deletes the resource even if it's part of a deployed configuration. If the configuration is deployed, the service will perform a cleanup deployment in the background before the deletion.

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
  network_group_ids        = ["azurerm_network_manager_network_group.example.id"]
  description              = "example routing rule collection"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Manager Routing Rule Collection. Changing this forces a new Network Manager Routing Rule Collection to be created.

* `routing_configuration_id` - (Required) The ID of the Network Manager Routing Configuration. Changing this forces a new Network Manager Routing Rule Collection to be created.

* `network_group_ids` - (Required) A list of Network Group IDs which this Network Manager Routing Rule Collection applies to.

---

* `bgp_route_propagation_enabled` - (Optional) Whether to enable the BGP route propagation. Defaults to `false`.

* `description` - (Optional) The description of the Network Manager Routing Rule Collection.

## Attribute Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Routing Rule Collection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Routing Rule Collection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Routing Rule Collection.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Routing Rule Collection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Routing Rule Collection.

## Import

Network Manager Routing Rule Collections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_routing_rule_collection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/manager1/routingConfigurations/conf1/ruleCollections/collection1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2024-05-01
