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

resource "azurerm_express_route_circuit" "example1" {
  name                  = "expressRoute1"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags = {
    environment = "Production"
  }
}

resource "azurerm_express_route_circuit" "example2" {
  name                  = "expressRoute2"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags = {
    environment = "Production"
  }
}

resource "azurerm_express_route_circuit_peering" "example1" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.example.name
  resource_group_name           = azurerm_resource_group.example.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "123.0.0.0/30"
  secondary_peer_address_prefix = "123.0.0.4/30"
  vlan_id                       = 300
}
resource "azurerm_express_route_circuit_peering" "example2" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.example.name
  resource_group_name           = azurerm_resource_group.example.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "123.0.0.0/30"
  secondary_peer_address_prefix = "123.0.0.4/30"
  vlan_id                       = 300
}

resource "azurerm_express_route_circuit_connection" "example" {
  name = "example-expressroutecircuitconnection"
  resource_group_name = azurerm_resource_group.example.name
  circuit_name = azurerm_express_route_circuit.example.name
  peering_id = azurerm_express_route_circuit_peering.example1.id
  peer_peering_id = azurerm_express_route_circuit_peering.example2.id
  address_prefix = "192.169.8.0/29"
  authorization_key = "00000000-0000-0000-0000-000000000000"
  ipv6circuit_connection_config {
    address_prefix = "2002:db01::/125"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this network ExpressRouteCircuitConnection. Changing this forces a new network ExpressRouteCircuitConnection to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the network ExpressRouteCircuitConnection should exist. Changing this forces a new network ExpressRouteCircuitConnection to be created.

* `circuit_name` - (Required) The name of the express route circuit. Changing this forces a new network ExpressRouteCircuitConnection to be created.

* `address_prefix` - (Required) /29 IP address space to carve out Customer addresses for tunnels.

---

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