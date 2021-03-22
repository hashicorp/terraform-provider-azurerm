---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit_connection"
description: |-
  Manages a network ExpressRouteCircuitConnection.
---

# azurerm_express_route_circuit_connection

Manages a network ExpressRouteCircuitConnection.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-network"
  location = "West Europe"
}

resource "azurerm_express_route_port" "example" {
  name = "example-expressrouteport"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
}

resource "azurerm_express_route_circuit" "example" {
  name = "example-expressroutecircuit"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
}

resource "azurerm_express_route_circuit_peering" "example" {
  name = "example-expressroutecircuitpeering"
  resource_group_name = azurerm_resource_group.example.name
  circuit_name = azurerm_express_route_circuit.example.name
}

resource "azurerm_express_route_circuit_connection" "example" {
  name = "example-expressroutecircuitconnection"
  resource_group_name = azurerm_resource_group.example.name
  circuit_name = azurerm_express_route_circuit.example.name
  peering_name = azurerm_express_route_circuit_peering.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this network ExpressRouteCircuitConnection. Changing this forces a new network ExpressRouteCircuitConnection to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the network ExpressRouteCircuitConnection should exist. Changing this forces a new network ExpressRouteCircuitConnection to be created.

* `circuit_name` - (Required) The name of the express route circuit. Changing this forces a new network ExpressRouteCircuitConnection to be created.

* `peering_name` - (Required) The name of the peering. Changing this forces a new network ExpressRouteCircuitConnection to be created.

---

* `address_prefix` - (Optional) /29 IP address space to carve out Customer addresses for tunnels.

* `authorization_key` - (Optional) The authorization key.

* `ipv6circuit_connection_config` - (Optional) A `ipv6circuit_connection_config` block as defined below.

---

An `ipv6circuit_connection_config` block exports the following:

* `address_prefix` - (Optional) /125 IP address space to carve out customer addresses for global reach.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the network ExpressRouteCircuitConnection.

* `circuit_connection_status` - Express Route Circuit connection state.

* `type` - Type of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the network ExpressRouteCircuitConnection.
* `read` - (Defaults to 5 minutes) Used when retrieving the network ExpressRouteCircuitConnection.
* `delete` - (Defaults to 30 minutes) Used when deleting the network ExpressRouteCircuitConnection.

## Import

network ExpressRouteCircuitConnections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_circuit_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/connection1
```