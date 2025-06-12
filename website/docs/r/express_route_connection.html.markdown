---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_connection"
description: |-
  Manages an Express Route Connection.
---

# azurerm_express_route_connection

Manages an Express Route Connection.

~> **Note:** The provider status of the Express Route Circuit must be set as provisioned while creating the Express Route Connection. See more details [here](https://docs.microsoft.com/azure/expressroute/expressroute-howto-circuit-portal-resource-manager#send-the-service-key-to-your-connectivity-provider-for-provisioning).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-vhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_express_route_gateway" "example" {
  name                = "example-expressroutegateway"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_hub_id      = azurerm_virtual_hub.example.id
  scale_units         = 1
}

resource "azurerm_express_route_port" "example" {
  name                = "example-erp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  peering_location    = "Equinix-Seattle-SE2"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "example" {
  name                  = "example-erc"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  express_route_port_id = azurerm_express_route_port.example.id
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
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
}

resource "azurerm_express_route_connection" "example" {
  name                             = "example-expressrouteconn"
  express_route_gateway_id         = azurerm_express_route_gateway.example.id
  express_route_circuit_peering_id = azurerm_express_route_circuit_peering.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Express Route Connection. Changing this forces a new resource to be created.

* `express_route_circuit_peering_id` - (Required) The ID of the Express Route Circuit Peering that this Express Route Connection connects with. Changing this forces a new resource to be created.

* `express_route_gateway_id` - (Required) The ID of the Express Route Gateway that this Express Route Connection connects with. Changing this forces a new resource to be created.

* `authorization_key` - (Optional) The authorization key to establish the Express Route Connection.

* `enable_internet_security` - (Optional) Is Internet security enabled for this Express Route Connection?

* `express_route_gateway_bypass_enabled` - (Optional) Specified whether Fast Path is enabled for Virtual Wan Firewall Hub. Defaults to `false`.

* `routing` - (Optional) A `routing` block as defined below.

* `routing_weight` - (Optional) The routing weight associated to the Express Route Connection. Possible value is between `0` and `32000`. Defaults to `0`.

---

A `routing` block supports the following:

* `associated_route_table_id` - (Optional) The ID of the Virtual Hub Route Table associated with this Express Route Connection.

* `inbound_route_map_id` - (Optional) The ID of the Route Map associated with this Express Route Connection for inbound routes.
 
* `outbound_route_map_id` - (Optional) The ID of the Route Map associated with this Express Route Connection for outbound routes.

* `propagated_route_table` - (Optional) A `propagated_route_table` block as defined below.

---

A `propagated_route_table` block supports the following:

* `labels` - (Optional) The list of labels to logically group route tables.

* `route_table_ids` - (Optional) A list of IDs of the Virtual Hub Route Table to propagate routes from Express Route Connection to the route table.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Express Route Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Express Route Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Express Route Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Express Route Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Express Route Connection.

## Import

Express Route Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/expressRouteGateways/expressRouteGateway1/expressRouteConnections/connection1
```
