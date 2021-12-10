---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit_connection"
description: |-
  Manages an Express Route Circuit Connection.
---

# azurerm_express_route_circuit_connection

Manages an Express Route Circuit Connection.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_express_route_port" "example" {
  name                = "example-erport"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  peering_location    = "Equinix-Seattle-SE2"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "example" {
  name                  = "example-ercircuit"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  express_route_port_id = azurerm_express_route_port.example.id
  bandwidth_in_gbps     = 5

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_port" "example2" {
  name                = "example-erport2"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  peering_location    = "Allied-Toronto-King-West"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "example2" {
  name                  = "example-ercircuit2"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  express_route_port_id = azurerm_express_route_port.example2.id
  bandwidth_in_gbps     = 5

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_circuit_peering" "example" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.example.name
  resource_group_name           = azurerm_resource_group.example.name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.1.0/30"
  vlan_id                       = 100
}

resource "azurerm_express_route_circuit_peering" "example2" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.example2.name
  resource_group_name           = azurerm_resource_group.example.name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.1.0/30"
  vlan_id                       = 100
}

resource "azurerm_express_route_circuit_connection" "example" {
  name                = "example-ercircuitconnection"
  peering_id          = azurerm_express_route_circuit_peering.example.id
  peer_peering_id     = azurerm_express_route_circuit_peering.example2.id
  address_prefix_ipv4 = "192.169.9.0/29"
  authorization_key   = "846a1918-b7a2-4917-b43c-8c4cdaee006a"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Express Route Circuit Connection. Changing this forces a new Express Route Circuit Connection to be created.

* `peering_id` - (Required) The ID of the Express Route Circuit Private Peering that this Express Route Circuit Connection connects with. Changing this forces a new Express Route Circuit Connection to be created.
  
* `peer_peering_id` - (Required) The ID of the peered Express Route Circuit Private Peering. Changing this forces a new Express Route Circuit Connection to be created.

* `address_prefix_ipv4` - (Required) The IPv4 address space from which to allocate customer address for global reach. Changing this forces a new Express Route Circuit Connection to be created.

---

* `authorization_key` - (Optional) The authorization key which is associated with the Express Route Circuit Connection.

* `address_prefix_ipv6` - (Optional) The IPv6 address space from which to allocate customer addresses for global reach.

-> **NOTE**: `address_prefix_ipv6` cannot be set when ExpressRoute Circuit Connection with ExpressRoute Circuit based on ExpressRoute Port.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Express Route Circuit Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Express Route Circuit Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Express Route Circuit Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Express Route Circuit Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Express Route Circuit Connection.

## Import

Express Route Circuit Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_circuit_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/expressRouteCircuits/circuit1/peerings/peering1/connections/connection1
```
